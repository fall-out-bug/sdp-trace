package witness

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

const (
	statePass         = "pass"
	stateFail         = "fail"
	stateCannotVerify = "cannot_verify"
	stateNotAssessed  = "not_assessed"

	independenceExternal = "external_independent"
	independenceCIJob    = "ci_isolated_job"
	independenceSameJob  = "ci_same_job"
)

var safetyClasses = []string{
	"ci_token",
	"oidc_token",
	"jwt_body",
	"private_key_material",
	"provider_token",
	"authenticated_provider_url",
	"raw_job_log",
	"private_filesystem_path",
	"unsafe_personal_identifier",
	"free_text_parser_error_with_input",
	"customer_directory_payload",
	"ldap_payload",
	"saml_payload",
	"cloud_payload",
	"vault_payload",
	"hsm_payload",
	"kms_payload",
	"pki_payload",
}

type EnvelopeInput struct {
	SchemaVersion       string           `json:"schema_version"`
	ProfileID           string           `json:"profile_id"`
	ProfileVersion      string           `json:"profile_version"`
	ProviderKind        string           `json:"provider_kind"`
	RequestedTrustScope string           `json:"requested_trust_scope"`
	GeneratedAt         string           `json:"generated_at"`
	Source              SourceIdentity   `json:"source"`
	CI                  CIIdentity       `json:"ci"`
	RunArtifacts        []ArtifactDigest `json:"run_artifacts"`
	ReportArtifacts     []ArtifactDigest `json:"report_artifacts"`
	ProfileStates       ProfileStates    `json:"profile_states"`
}

type ProfileOptions struct {
	EnvelopePath               string
	CustomerPKIAuthorityPolicy string
	CustomerPKIPublicCert      string
	CustomerPKIPublicKey       string
	CustomerPKIPayloadDigest   string
	CustomerPKIFreshness       string
}

type CustomerPKIAuthorityPolicy struct {
	SchemaVersion      string `json:"schema_version"`
	ProfileID          string `json:"profile_id"`
	AllowedSignerID    string `json:"allowed_signer_id"`
	PublicKeySHA256    string `json:"public_key_sha256"`
	PolicyDigest       string `json:"policy_digest"`
	KeyCustodyState    string `json:"key_custody_state"`
	RevocationRequired bool   `json:"revocation_required"`
	RevocationState    string `json:"revocation_state"`
}

type CustomerPKIFreshnessEvidence struct {
	SchemaVersion string `json:"schema_version"`
	SignerID      string `json:"signer_id"`
	PayloadDigest string `json:"payload_digest"`
	RunID         string `json:"run_id"`
	PolicyDigest  string `json:"policy_digest"`
	IssuedAt      string `json:"issued_at"`
	ValidUntil    string `json:"valid_until"`
	Nonce         string `json:"nonce"`
	Signature     string `json:"signature"`
}

func WriteProfile(kind, outPath, runsRoot, reportDir string, opts ProfileOptions) (Record, error) {
	if strings.TrimSpace(outPath) == "" {
		return Record{}, errors.New("witness requires --out <file>")
	}
	record, err := BuildProfile(kind, runsRoot, reportDir, opts)
	if err != nil {
		return Record{}, err
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return Record{}, err
	}
	record = finalizeRecordForWrite(record)
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return Record{}, err
	}
	if err := os.WriteFile(outPath, append(data, '\n'), 0o644); err != nil {
		return Record{}, err
	}
	return record, nil
}

func BuildProfile(kind, runsRoot, reportDir string, opts ProfileOptions) (Record, error) {
	switch kind {
	case KindGitLabCI, KindBuildkite:
		return BuildCIEnvelopeProfile(kind, runsRoot, reportDir, opts.EnvelopePath)
	case KindCustomerPKI:
		return BuildCustomerPKI(runsRoot, reportDir, opts)
	default:
		return Record{}, fmt.Errorf("unsupported witness kind %q", kind)
	}
}

func BuildCIEnvelopeProfile(kind, runsRoot, reportDir, envelopePath string) (Record, error) {
	record := baseRecord(kind)
	runArtifacts, err := hashRunArtifacts(runsRoot)
	if err != nil {
		return Record{}, err
	}
	record.RunArtifacts = runArtifacts
	if reportDir != "" {
		reportArtifacts, err := hashReportArtifacts(reportDir)
		if err != nil {
			return Record{}, err
		}
		record.ReportArtifacts = reportArtifacts
	}
	if strings.TrimSpace(envelopePath) == "" {
		reason := ReasonMissingIdentity
		if ambientCIEnvPresent(kind) {
			reason = ReasonEnvOnly
		}
		record.Status = StatusCannotVerify
		record.TrustScope = TrustScopeLocalObserved
		record.EstablishedTrustScope = stateCannotVerify
		record.Reason = reason
		record.ReasonCodes = []string{reason}
		record.ProfileStates = defaultProfileStates(stateCannotVerify, independenceSameJob)
		return record, nil
	}
	var envelope EnvelopeInput
	if err := readSafeJSON(envelopePath, &envelope); err != nil {
		record.Status = StatusCannotVerify
		record.TrustScope = TrustScopeLocalObserved
		record.EstablishedTrustScope = stateCannotVerify
		record.Reason = ReasonMalformedInput
		record.ReasonCodes = []string{ReasonMalformedInput}
		return record, nil
	}
	if unsafeEnvelopeFields(envelope) {
		record.Status = StatusFail
		record.TrustScope = stateFail
		record.EstablishedTrustScope = stateFail
		record.Reason = ReasonUnsafeOutput
		record.ReasonCodes = []string{ReasonUnsafeOutput}
		record.ProfileStates = defaultProfileStates(stateFail, "cannot_verify")
		return record, nil
	}
	record.SchemaVersion = defaultString(envelope.SchemaVersion, "sdp-trace-witness-profile-result/v1")
	record.ProfileID = defaultString(envelope.ProfileID, kind+"-v1")
	record.ProfileVersion = defaultString(envelope.ProfileVersion, "1.0")
	record.ProviderKind = defaultString(envelope.ProviderKind, kind)
	record.RequestedTrustScope = defaultString(envelope.RequestedTrustScope, TrustScopeCIWitnessed)
	record.Source = envelope.Source
	record.CI = envelope.CI
	record.ProfileStates = &envelope.ProfileStates
	if strings.TrimSpace(envelope.GeneratedAt) != "" {
		record.GeneratedAt = envelope.GeneratedAt
	}
	if state := validateCIEnvelope(kind, envelope, runArtifacts); state.reason != "" {
		applyProfileState(&record, state.status, state.scope, state.reason)
		return record, nil
	}
	if !runIDMatches(runsRoot, envelope.CI.RunID) {
		record.ProfileStates.RunBindingState = stateFail
		applyProfileState(&record, StatusFail, stateFail, ReasonRunMismatch)
		return record, nil
	}
	record.Status = StatusPass
	record.TrustScope = TrustScopeCIWitnessed
	record.EstablishedTrustScope = TrustScopeCIWitnessed
	record.Reason = ReasonProfileVerified
	return record, nil
}

func BuildCustomerPKI(runsRoot, reportDir string, opts ProfileOptions) (Record, error) {
	record, err := newCustomerPKIRecord(runsRoot, reportDir)
	if err != nil {
		return Record{}, err
	}
	publicKey, policy, freshness, ok := loadCustomerPKIInputs(&record, opts)
	if !ok {
		return record, nil
	}
	record.CI = CIIdentity{Provider: KindCustomerPKI, RunID: freshness.RunID}
	states := customerPKIPassStates(policy)
	record.ProfileStates = states

	if !validateCustomerPKIAuthority(&record, states, publicKey, policy, freshness) {
		return record, nil
	}
	if !validateCustomerPKIFreshness(&record, states, runsRoot, opts.CustomerPKIPayloadDigest, freshness) {
		return record, nil
	}
	if !verifyFreshnessSignature(publicKey, freshness) {
		customerPKIFail(&record, states, "freshness", ReasonSignerMismatch)
		return record, nil
	}
	record.Status = StatusPass
	record.TrustScope = TrustScopeExternal
	record.EstablishedTrustScope = TrustScopeExternal
	record.Reason = ReasonProfileVerified
	return record, nil
}

func newCustomerPKIRecord(runsRoot, reportDir string) (Record, error) {
	record := baseRecord(KindCustomerPKI)
	runArtifacts, err := hashRunArtifacts(runsRoot)
	if err != nil {
		return Record{}, err
	}
	record.RunArtifacts = runArtifacts
	if reportDir != "" {
		reportArtifacts, err := hashReportArtifacts(reportDir)
		if err != nil {
			return Record{}, err
		}
		record.ReportArtifacts = reportArtifacts
	}
	record.ProfileID = "customer-pki-v1"
	record.ProfileVersion = "1.0"
	record.ProviderKind = KindCustomerPKI
	record.RequestedTrustScope = TrustScopeExternal
	return record, nil
}

func loadCustomerPKIInputs(record *Record, opts ProfileOptions) (ed25519.PublicKey, CustomerPKIAuthorityPolicy, CustomerPKIFreshnessEvidence, bool) {
	requiredMissing := missingCustomerPKIInputs(opts)
	if len(requiredMissing) > 0 {
		record.MissingIdentityFields = requiredMissing
		customerPKICannotVerify(record, ReasonMissingIdentity)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	if privateKeyInput(opts.CustomerPKIPublicCert) || privateKeyInput(opts.CustomerPKIPublicKey) {
		customerPKIInputFail(record, ReasonPrivateKeyInput)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	publicKey, err := loadCustomerPublicKey(opts)
	if err != nil {
		customerPKICannotVerify(record, ReasonMalformedInput)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	var policy CustomerPKIAuthorityPolicy
	if err := readSafeJSON(opts.CustomerPKIAuthorityPolicy, &policy); err != nil {
		customerPKICannotVerify(record, ReasonPolicyMissing)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	var freshness CustomerPKIFreshnessEvidence
	if err := readSafeJSON(opts.CustomerPKIFreshness, &freshness); err != nil {
		customerPKICannotVerify(record, ReasonMissingFreshness)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	return publicKey, policy, freshness, true
}

func customerPKICannotVerify(record *Record, reason string) {
	applyProfileState(record, StatusCannotVerify, stateCannotVerify, reason)
	record.ProfileStates = defaultProfileStates(stateCannotVerify, independenceExternal)
}

func customerPKIInputFail(record *Record, reason string) {
	applyProfileState(record, StatusFail, stateFail, reason)
	record.ProfileStates = defaultProfileStates(stateFail, independenceExternal)
}

func customerPKIPassStates(policy CustomerPKIAuthorityPolicy) *ProfileStates {
	return &ProfileStates{
		IdentityState:        statePass,
		SignerAuthorityState: statePass,
		FreshnessState:       statePass,
		ArtifactBindingState: statePass,
		SourceBindingState:   stateNotAssessed,
		RunBindingState:      statePass,
		PolicyBindingState:   statePass,
		IndependenceState:    independenceExternal,
		KeyCustodyState:      defaultString(policy.KeyCustodyState, "not_assessed"),
	}
}

func validateCustomerPKIAuthority(record *Record, states *ProfileStates, publicKey ed25519.PublicKey, policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence) bool {
	if policy.ProfileID != "customer-pki-v1" || policy.AllowedSignerID == "" || policy.AllowedSignerID != freshness.SignerID {
		customerPKIFail(record, states, "signer", ReasonSignerMismatch)
		return false
	}
	if policy.PublicKeySHA256 != "" && policy.PublicKeySHA256 != digestBytes(publicKey) {
		customerPKIFail(record, states, "signer", ReasonSignerMismatch)
		return false
	}
	if policy.PolicyDigest != "" && policy.PolicyDigest != freshness.PolicyDigest {
		customerPKIFail(record, states, "policy", ReasonPolicyMismatch)
		return false
	}
	if policy.RevocationRequired && policy.RevocationState == "" {
		states.SignerAuthorityState = stateNotAssessed
		applyProfileState(record, StatusNotAssessed, stateNotAssessed, ReasonRevocationNA)
		return false
	}
	if policy.RevocationState == "revoked" {
		customerPKIFail(record, states, "signer", ReasonCertRevoked)
		return false
	}
	return true
}

func validateCustomerPKIFreshness(record *Record, states *ProfileStates, runsRoot, payloadDigest string, freshness CustomerPKIFreshnessEvidence) bool {
	if payloadDigest != freshness.PayloadDigest || !strongDigest(freshness.PayloadDigest) {
		customerPKIFail(record, states, "artifact", ReasonArtifactMismatch)
		return false
	}
	if !freshnessCurrent(freshness, time.Now().UTC()) {
		customerPKIFail(record, states, "freshness", ReasonStaleFreshness)
		return false
	}
	if !runIDMatches(runsRoot, freshness.RunID) {
		customerPKIFail(record, states, "run", ReasonRunMismatch)
		return false
	}
	return true
}

func customerPKIFail(record *Record, states *ProfileStates, field, reason string) {
	switch field {
	case "artifact":
		states.ArtifactBindingState = stateFail
	case "freshness":
		states.FreshnessState = stateFail
	case "policy":
		states.PolicyBindingState = stateFail
	case "run":
		states.RunBindingState = stateFail
	case "signer":
		states.SignerAuthorityState = stateFail
	default:
		states.IdentityState = stateFail
	}
	applyProfileState(record, StatusFail, stateFail, reason)
}

type profileDecision struct {
	status string
	scope  string
	reason string
}

func validateCIEnvelope(kind string, envelope EnvelopeInput, current []ArtifactDigest) profileDecision {
	if envelope.ProfileID != kind+"-v1" {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonUnsupported}
	}
	if strings.TrimSpace(envelope.Source.CommitSHA) == "" || strings.TrimSpace(envelope.CI.RunID) == "" {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingIdentity}
	}
	if state := validateCIEnvelopeStates(envelope.ProfileStates); state.reason != "" {
		return state
	}
	if len(envelope.RunArtifacts) == 0 {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingArtifact}
	}
	if !artifactSetsMatch(envelope.RunArtifacts, current) {
		return profileDecision{StatusFail, stateFail, ReasonArtifactMismatch}
	}
	return profileDecision{}
}

func validateCIEnvelopeStates(states ProfileStates) profileDecision {
	if states.IdentityState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingIdentity}
	}
	if states.SignerAuthorityState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingSigner}
	}
	if states.FreshnessState != statePass {
		if states.FreshnessState == stateFail {
			return profileDecision{StatusFail, stateFail, ReasonStaleFreshness}
		}
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingFreshness}
	}
	if states.SourceBindingState != statePass {
		return profileDecision{StatusFail, stateFail, ReasonSourceMismatch}
	}
	if states.RunBindingState != statePass {
		return profileDecision{StatusFail, stateFail, ReasonRunMismatch}
	}
	if states.PolicyBindingState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonPolicyMissing}
	}
	if states.ArtifactBindingState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingArtifact}
	}
	if states.IndependenceState != independenceCIJob && states.IndependenceState != independenceExternal {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonEnvOnly}
	}
	return profileDecision{}
}

func baseRecord(kind string) Record {
	return Record{
		SchemaVersion:         "sdp-trace-witness-profile-result/v1",
		Kind:                  kind,
		ProfileID:             kind + "-v1",
		ProfileVersion:        "1.0",
		ProviderKind:          kind,
		Status:                StatusCannotVerify,
		TrustScope:            TrustScopeLocalObserved,
		RequestedTrustScope:   TrustScopeCIWitnessed,
		EstablishedTrustScope: stateCannotVerify,
		GeneratedAt:           time.Now().UTC().Format(time.RFC3339Nano),
		RunArtifacts:          []ArtifactDigest{},
		ReportArtifacts:       []ArtifactDigest{},
		OutputSafety: &OutputSafety{
			State:                 statePass,
			VerifiedAbsentClasses: safetyClasses,
		},
	}
}

func applyProfileState(record *Record, status, scope, reason string) {
	record.Status = status
	record.TrustScope = scope
	record.EstablishedTrustScope = scope
	record.Reason = reason
	record.ReasonCodes = []string{reason}
}

func finalizeRecordForWrite(record Record) Record {
	scanRecord := record
	scanRecord.OutputSafety = nil
	raw, err := json.Marshal(scanRecord)
	if err != nil {
		applyProfileState(&record, StatusFail, stateFail, ReasonUnsafeOutput)
		return record
	}
	if !forbiddenOutputPresent(raw) {
		if record.OutputSafety == nil {
			record.OutputSafety = &OutputSafety{}
		}
		record.OutputSafety.State = statePass
		record.OutputSafety.VerifiedAbsentClasses = safetyClasses
		return record
	}
	safe := baseRecord(record.Kind)
	safe.ProfileStates = defaultProfileStates(stateFail, "cannot_verify")
	safe.OutputSafety = &OutputSafety{State: stateFail, VerifiedAbsentClasses: safetyClasses}
	applyProfileState(&safe, StatusFail, stateFail, ReasonUnsafeOutput)
	return safe
}

func forbiddenOutputPresent(raw []byte) bool {
	if containsSecretLike(raw) {
		return true
	}
	text := strings.ToLower(string(raw))
	for _, marker := range []string{
		"https://user:",
		"https://token:",
		"bearer ",
		"ghp_",
		"glpat-",
		"xoxb-",
		"/private/",
		"raw_job_log_sentinel",
		"oidc.jwt.",
		"customer_directory_payload",
		"ldap_payload",
		"saml_payload",
		"cloud_payload",
		"vault_payload",
		"hsm_payload",
		"kms_payload",
		"pki_payload",
		"free_text_parser_error_with_input",
	} {
		if strings.Contains(text, marker) {
			return true
		}
	}
	return jwtLike(text)
}

func defaultProfileStates(state, independence string) *ProfileStates {
	return &ProfileStates{
		IdentityState:        state,
		SignerAuthorityState: state,
		FreshnessState:       state,
		ArtifactBindingState: state,
		SourceBindingState:   state,
		RunBindingState:      state,
		PolicyBindingState:   state,
		IndependenceState:    independence,
	}
}

func artifactSetsMatch(expected, current []ArtifactDigest) bool {
	if len(expected) != len(current) {
		return false
	}
	byPath := map[string]string{}
	for _, artifact := range current {
		byPath[artifact.Path] = artifact.SHA256
	}
	for _, artifact := range expected {
		if byPath[artifact.Path] != artifact.SHA256 {
			return false
		}
	}
	return true
}

func runIDMatches(runsRoot, witnessRunID string) bool {
	if strings.TrimSpace(witnessRunID) == "" {
		return false
	}
	runIDs, err := runIDsFromRoot(runsRoot)
	if err != nil || len(runIDs) == 0 {
		return false
	}
	for _, runID := range runIDs {
		if runID == witnessRunID {
			return true
		}
	}
	return false
}

func runIDsFromRoot(runsRoot string) ([]string, error) {
	runDirs, err := demo.DiscoverRunDirs(runsRoot)
	if err != nil {
		return nil, err
	}
	runIDs := make([]string, 0, len(runDirs))
	for _, runDir := range runDirs {
		raw, err := os.ReadFile(filepath.Join(runDir, "run.json"))
		if err != nil {
			return nil, err
		}
		var payload struct {
			RunID string `json:"run_id"`
			ID    string `json:"id"`
		}
		if err := json.Unmarshal(raw, &payload); err != nil {
			return nil, err
		}
		switch {
		case strings.TrimSpace(payload.RunID) != "":
			runIDs = append(runIDs, payload.RunID)
		case strings.TrimSpace(payload.ID) != "":
			runIDs = append(runIDs, payload.ID)
		}
	}
	return runIDs, nil
}

func ambientCIEnvPresent(kind string) bool {
	prefixes := map[string][]string{
		KindGitLabCI:  {"GITLAB_CI", "CI_PIPELINE_ID", "CI_JOB_ID", "CI_COMMIT_SHA"},
		KindBuildkite: {"BUILDKITE", "BUILDKITE_BUILD_ID", "BUILDKITE_JOB_ID", "BUILDKITE_COMMIT"},
	}
	for _, key := range prefixes[kind] {
		if strings.TrimSpace(os.Getenv(key)) != "" {
			return true
		}
	}
	return false
}

func missingCustomerPKIInputs(opts ProfileOptions) []string {
	missing := []string{}
	required := map[string]string{
		"--customer-pki-authority-policy":   opts.CustomerPKIAuthorityPolicy,
		"--customer-pki-payload-digest":     opts.CustomerPKIPayloadDigest,
		"--customer-pki-freshness-evidence": opts.CustomerPKIFreshness,
	}
	for name, value := range required {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}
	if strings.TrimSpace(opts.CustomerPKIPublicCert) == "" && strings.TrimSpace(opts.CustomerPKIPublicKey) == "" {
		missing = append(missing, "--customer-pki-public-cert|--customer-pki-public-key")
	}
	sort.Strings(missing)
	return missing
}

func readSafeJSON(path string, target any) error {
	if unsafeInputPath(path) {
		return fmt.Errorf("unsafe input path")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if containsSecretLike(raw) {
		return fmt.Errorf("unsafe input content")
	}
	return json.Unmarshal(raw, target)
}

func unsafeInputPath(path string) bool {
	if strings.TrimSpace(path) == "" || strings.Contains(path, "\x00") {
		return true
	}
	lower := strings.ToLower(filepath.ToSlash(path))
	if strings.Contains(lower, "://") || strings.Contains(lower, "..") || strings.Contains(lower, "private.key") {
		return true
	}
	info, err := os.Lstat(path)
	return err == nil && info.Mode()&os.ModeSymlink != 0
}

func containsSecretLike(raw []byte) bool {
	lower := strings.ToLower(string(raw))
	for _, marker := range []string{
		"-----begin private key-----",
		"-----begin rsa private key-----",
		"-----begin ec private key-----",
		"token_secret_",
		"jwt_secret_",
		"oidc.jwt.",
		"bearer ",
		"ghp_",
		"glpat-",
		"xoxb-",
		"https://user:",
		"https://token:",
		"raw_job_log_sentinel",
		"customer_directory_payload",
		"ldap_payload",
		"saml_payload",
		"cloud_payload",
		"vault_payload",
		"hsm_payload",
		"kms_payload",
		"pki_payload",
		"free_text_parser_error_with_input",
	} {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return jwtLike(lower)
}

func jwtLike(text string) bool {
	for _, field := range strings.FieldsFunc(text, func(r rune) bool {
		switch r {
		case '"', '\'', ' ', '\n', '\t', '\r', ',', ':', '{', '}', '[', ']', '(', ')':
			return true
		default:
			return false
		}
	}) {
		parts := strings.Split(field, ".")
		if len(parts) == 3 && strings.HasPrefix(parts[0], "eyj") && len(parts[1]) >= 8 && len(parts[2]) >= 8 {
			return true
		}
	}
	return false
}

func unsafeEnvelopeFields(envelope EnvelopeInput) bool {
	values := []string{
		envelope.Source.Repository,
		envelope.Source.Ref,
		envelope.Source.CommitSHA,
		envelope.CI.Provider,
		envelope.CI.ServerURL,
		envelope.CI.Workflow,
		envelope.CI.Job,
		envelope.CI.RunID,
		envelope.CI.RunAttempt,
		envelope.CI.Actor,
	}
	for _, value := range values {
		if unsafeOutputString(value) {
			return true
		}
	}
	for _, artifact := range append(envelope.RunArtifacts, envelope.ReportArtifacts...) {
		if unsafeOutputString(artifact.Path) || unsafeOutputString(artifact.SHA256) {
			return true
		}
	}
	return false
}

func unsafeOutputString(value string) bool {
	if value == "" {
		return false
	}
	if containsSecretLike([]byte(value)) {
		return true
	}
	lower := strings.ToLower(value)
	return strings.Contains(lower, "/private/") || strings.Contains(lower, "@")
}

func privateKeyInput(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return containsSecretLike(raw)
}

func loadCustomerPublicKey(opts ProfileOptions) (ed25519.PublicKey, error) {
	path := opts.CustomerPKIPublicKey
	if path == "" {
		path = opts.CustomerPKIPublicCert
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if containsSecretLike(raw) {
		return nil, errors.New("private key input rejected")
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("public key or certificate PEM is required")
	}
	if block.Type == "CERTIFICATE" {
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}
		if key, ok := cert.PublicKey.(ed25519.PublicKey); ok {
			return key, nil
		}
		return nil, errors.New("certificate must contain ed25519 public key")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	edKey, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("public key must be ed25519")
	}
	return edKey, nil
}

func verifyFreshnessSignature(publicKey ed25519.PublicKey, evidence CustomerPKIFreshnessEvidence) bool {
	signature, err := base64.StdEncoding.DecodeString(evidence.Signature)
	if err != nil {
		return false
	}
	return ed25519.Verify(publicKey, []byte(freshnessPayload(evidence)), signature)
}

func freshnessPayload(evidence CustomerPKIFreshnessEvidence) string {
	return strings.Join([]string{
		evidence.PayloadDigest,
		evidence.RunID,
		evidence.PolicyDigest,
		evidence.SignerID,
		evidence.IssuedAt,
		evidence.ValidUntil,
		evidence.Nonce,
	}, "\n")
}

func freshnessCurrent(evidence CustomerPKIFreshnessEvidence, now time.Time) bool {
	issued, err := time.Parse(time.RFC3339, evidence.IssuedAt)
	if err != nil || issued.After(now.Add(time.Minute)) {
		return false
	}
	if evidence.ValidUntil == "" {
		return true
	}
	validUntil, err := time.Parse(time.RFC3339, evidence.ValidUntil)
	return err == nil && !validUntil.Before(now)
}

func digestBytes(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func strongDigest(value string) bool {
	if len(value) < 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
