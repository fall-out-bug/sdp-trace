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

var outputSafetyMarkers = []string{
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
}

var secretSafetyMarkers = []string{
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
	// Profile files are publication artifacts. Build and safety finalization
	// happen before the write so callers never persist an unscanned trust claim.
	record, err := buildProfileForWrite(kind, runsRoot, reportDir, opts)
	if err != nil {
		return Record{}, err
	}
	if err := writeProfileRecord(outPath, record); err != nil {
		return Record{}, err
	}
	return record, nil
}

func buildProfileForWrite(kind, runsRoot, reportDir string, opts ProfileOptions) (Record, error) {
	// BuildProfile may produce pass, fail, cannot_verify, or not_assessed
	// records. finalizeRecordForWrite is the final output-safety gate that can
	// still replace an otherwise useful record with a redacted failure record.
	record, err := BuildProfile(kind, runsRoot, reportDir, opts)
	if err != nil {
		return Record{}, err
	}
	return finalizeRecordForWrite(record), nil
}

func writeProfileRecord(outPath string, record Record) error {
	// The filesystem write is deliberately dumb: all trust decisions have
	// already been made, and this layer only materializes the exact record.
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(outPath, append(data, '\n'), 0o644); err != nil {
		return err
	}
	return nil
}

func BuildProfile(kind, runsRoot, reportDir string, opts ProfileOptions) (Record, error) {
	if ciEnvelopeProfile(kind) {
		// GitLab and Buildkite profiles are envelope-backed because this package
		// does not obtain provider-native attestations for those systems.
		return BuildCIEnvelopeProfile(kind, runsRoot, reportDir, opts.EnvelopePath)
	}
	switch kind {
	case KindCustomerPKI:
		// Customer PKI is external evidence: the profile can pass only after
		// policy, signer, freshness, signature, and run binding all agree.
		return BuildCustomerPKI(runsRoot, reportDir, opts)
	default:
		return Record{}, fmt.Errorf("unsupported witness kind %q", kind)
	}
}

func ciEnvelopeProfile(kind string) bool {
	return kind == KindGitLabCI || kind == KindBuildkite
}

func BuildCIEnvelopeProfile(kind, runsRoot, reportDir, envelopePath string) (Record, error) {
	record := baseRecord(kind)
	// Current artifact hashes are computed before reading the envelope so the
	// envelope can be compared with live local evidence, not trusted by itself.
	if err := populateCIEnvelopeArtifacts(&record, runsRoot, reportDir); err != nil {
		return Record{}, err
	}
	if !applyCIEnvelopeInputState(&record, kind, runsRoot, envelopePath) {
		return record, nil
	}
	return record, nil
}

func applyCIEnvelopeInputState(record *Record, kind, runsRoot, envelopePath string) bool {
	if strings.TrimSpace(envelopePath) == "" {
		// Ambient CI variables without a portable envelope are observation only;
		// they cannot establish a replayable provider witness.
		applyCIMissingEnvelopeState(record, kind)
		return false
	}
	envelope, ok := loadSafeCIEnvelopeRecord(record, envelopePath)
	if !ok {
		return false
	}
	applyCIEnvelopeRecordValues(record, kind, envelope)
	return applyCIEnvelopeTrustDecision(record, kind, runsRoot, envelope)
}

func populateCIEnvelopeArtifacts(record *Record, runsRoot, reportDir string) error {
	// Artifact population defines the comparison set for later envelope
	// validation. A missing or unreadable run artifact blocks the profile
	// instead of silently shrinking the evidence surface.
	runArtifacts, err := hashRunArtifacts(runsRoot)
	if err != nil {
		return err
	}
	record.RunArtifacts = runArtifacts
	if reportDir == "" {
		return nil
	}
	reportArtifacts, err := hashReportArtifacts(reportDir)
	if err != nil {
		return err
	}
	record.ReportArtifacts = reportArtifacts
	return nil
}

func applyCIMissingEnvelopeState(record *Record, kind string) {
	// Existing provider environment is weaker than a signed or exported
	// envelope. Preserve that distinction in the reason so downstream gates do
	// not treat "CI was present" as proof of CI-witnessed trust.
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
}

func loadSafeCIEnvelopeRecord(record *Record, envelopePath string) (EnvelopeInput, bool) {
	// Envelope loading is the first point where external CI profile facts can
	// affect the record, so malformed or unsafe inputs lower the record directly.
	// The boolean return forces callers to stop before copying untrusted envelope
	// values into the witness record.
	// The output record records the failure class, not the unsafe envelope body.
	var envelope EnvelopeInput
	if err := readSafeJSON(envelopePath, &envelope); err != nil {
		// Malformed or unsafe envelope input is cannot_verify because no
		// trustworthy envelope facts can be extracted.
		applyMalformedEnvelopeState(record)
		return envelope, false
	}
	if unsafeEnvelopeFields(envelope) {
		// Unsafe fields in a parsed envelope are an output-safety failure, not
		// just missing evidence, because persisting them could leak credentials.
		applyUnsafeEnvelopeState(record)
		return envelope, false
	}
	return envelope, true
}

func applyMalformedEnvelopeState(record *Record) {
	applyProfileState(record, StatusCannotVerify, stateCannotVerify, ReasonMalformedInput)
	record.TrustScope = TrustScopeLocalObserved
}

func applyUnsafeEnvelopeState(record *Record) {
	applyProfileState(record, StatusFail, stateFail, ReasonUnsafeOutput)
	record.ProfileStates = defaultProfileStates(stateFail, "cannot_verify")
}

func applyCIEnvelopeRecordValues(record *Record, kind string, envelope EnvelopeInput) {
	// Envelope metadata is copied only after safe parsing. Empty metadata keeps
	// deterministic defaults so missing fields do not become implicit claims.
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
}

func applyCIEnvelopeTrustDecision(record *Record, kind, runsRoot string, envelope EnvelopeInput) bool {
	// Validation first checks envelope self-consistency and artifact binding,
	// then separately binds the claimed run ID to discovered local runs.
	state := validateCIEnvelope(kind, envelope, record.RunArtifacts)
	if state.reason != "" {
		applyProfileState(record, state.status, state.scope, state.reason)
		return false
	}
	return setCIEnvelopeRunBindingState(record, runsRoot, envelope.CI.RunID)
}

func setCIEnvelopeRunBindingState(record *Record, runsRoot, witnessRunID string) bool {
	if runIDMatches(runsRoot, witnessRunID) {
		// A passing envelope profile requires both valid envelope states and a
		// run ID that resolves to a current discovered run artifact.
		record.Status = StatusPass
		record.TrustScope = TrustScopeCIWitnessed
		record.EstablishedTrustScope = TrustScopeCIWitnessed
		record.Reason = ReasonProfileVerified
		return true
	}
	record.ProfileStates.RunBindingState = stateFail
	applyProfileState(record, StatusFail, stateFail, ReasonRunMismatch)
	return false
}

func BuildCustomerPKI(runsRoot, reportDir string, opts ProfileOptions) (Record, error) {
	// Customer PKI can only establish external witness trust after local
	// artifacts, authority policy, public key, and freshness evidence all bind.
	// Early cannot_verify/fail records are returned with artifact digests already
	// captured where possible.
	// The happy path sets CI identity from freshness evidence only after the
	// external input files load safely.
	record, err := newCustomerPKIRecord(runsRoot, reportDir)
	if err != nil {
		return Record{}, err
	}
	// Customer PKI starts with local artifact hashes, then imports only the
	// minimal external inputs needed to authorize those artifacts.
	inputs, ok := prepareCustomerPKIValidation(&record, opts)
	if !ok {
		return record, nil
	}
	if !validateCustomerPKIRecord(&record, inputs.states, inputs.publicKey, inputs.policy, inputs.freshness, runsRoot, opts.CustomerPKIPayloadDigest) {
		return record, nil
	}
	applyCustomerPKIPass(&record)
	return record, nil
}

type customerPKIValidationInputs struct {
	publicKey ed25519.PublicKey
	policy    CustomerPKIAuthorityPolicy
	freshness CustomerPKIFreshnessEvidence
	states    *ProfileStates
}

func prepareCustomerPKIValidation(record *Record, opts ProfileOptions) (customerPKIValidationInputs, bool) {
	// External inputs become validation context only after safety checks have
	// accepted the public trust anchor, authority policy, and freshness evidence.
	publicKey, policy, freshness, ok := loadCustomerPKIInputs(record, opts)
	if !ok {
		return customerPKIValidationInputs{}, false
	}
	record.CI = CIIdentity{Provider: KindCustomerPKI, RunID: freshness.RunID}
	states := customerPKIPassStates(policy)
	record.ProfileStates = states
	return customerPKIValidationInputs{publicKey: publicKey, policy: policy, freshness: freshness, states: states}, true
}

func applyCustomerPKIPass(record *Record) {
	// Passing Customer PKI records establish external trust only after every
	// policy, freshness, signature, and run-binding gate has returned pass.
	record.Status = StatusPass
	record.TrustScope = TrustScopeExternal
	record.EstablishedTrustScope = TrustScopeExternal
	record.Reason = ReasonProfileVerified
}

func validateCustomerPKIRecord(record *Record, states *ProfileStates, publicKey ed25519.PublicKey, policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence, runsRoot, payloadDigest string) bool {
	// Authority, freshness, and signature are checked as separate gates so the
	// failing profile state can name the exact evidence boundary that broke.
	if !validateCustomerPKIAuthority(record, states, publicKey, policy, freshness) {
		return false
	}
	if !validateCustomerPKIFreshness(record, states, runsRoot, payloadDigest, freshness) {
		return false
	}
	if !verifyFreshnessSignature(publicKey, freshness) {
		customerPKIFail(record, states, "freshness", ReasonSignerMismatch)
		return false
	}
	return true
}

func newCustomerPKIRecord(runsRoot, reportDir string) (Record, error) {
	record := baseRecord(KindCustomerPKI)
	// Artifact digests are captured before external policy validation. The
	// external signer authorizes this concrete payload, not an abstract run.
	// Report artifacts are optional supporting material and do not replace run
	// artifact binding.
	// RequestedTrustScope is set to external before validation so failures still
	// show which trust upgrade was attempted.
	// ProviderKind stays customer-pki because this profile is not tied to a CI
	// vendor.
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
	if missing := missingCustomerPKIInputs(opts); len(missing) > 0 {
		// Missing PKI inputs leave the external authority unresolved; this is
		// cannot_verify rather than fail because no contradictory evidence exists.
		record.MissingIdentityFields = missing
		customerPKICannotVerify(record, ReasonMissingIdentity)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	if privateKeyInput(opts.CustomerPKIPublicCert) || privateKeyInput(opts.CustomerPKIPublicKey) {
		// Private key material in an input slot is a hard safety failure because
		// witness generation must never require or preserve signing secrets.
		customerPKIInputFail(record, ReasonPrivateKeyInput)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	return loadCustomerPKISafeInputs(record, opts)
}

func loadCustomerPKISafeInputs(record *Record, opts ProfileOptions) (ed25519.PublicKey, CustomerPKIAuthorityPolicy, CustomerPKIFreshnessEvidence, bool) {
	// Load the trust anchor before policy JSON so malformed key material cannot
	// be hidden behind otherwise well-formed policy metadata.
	publicKey, err := loadCustomerPublicKey(opts)
	if err != nil {
		customerPKICannotVerify(record, ReasonMalformedInput)
		return nil, CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	policy, freshness, ok := loadCustomerPKIJSONInputs(record, opts)
	return publicKey, policy, freshness, ok
}

func loadCustomerPKIJSONInputs(record *Record, opts ProfileOptions) (CustomerPKIAuthorityPolicy, CustomerPKIFreshnessEvidence, bool) {
	var policy CustomerPKIAuthorityPolicy
	if err := readSafeJSON(opts.CustomerPKIAuthorityPolicy, &policy); err != nil {
		// Policy JSON is authority evidence. If it cannot be read safely, no
		// signer can be accepted for the profile.
		customerPKICannotVerify(record, ReasonPolicyMissing)
		return CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	var freshness CustomerPKIFreshnessEvidence
	if err := readSafeJSON(opts.CustomerPKIFreshness, &freshness); err != nil {
		// Freshness JSON links the signed payload to a time window and run ID; a
		// missing value keeps the profile open instead of failing unrelated gates.
		customerPKICannotVerify(record, ReasonMissingFreshness)
		return CustomerPKIAuthorityPolicy{}, CustomerPKIFreshnessEvidence{}, false
	}
	return policy, freshness, true
}

func customerPKICannotVerify(record *Record, reason string) {
	// cannot_verify keeps the requested external profile visible while making
	// clear that no external trust scope was established.
	applyProfileState(record, StatusCannotVerify, stateCannotVerify, reason)
	record.ProfileStates = defaultProfileStates(stateCannotVerify, independenceExternal)
}

func customerPKIInputFail(record *Record, reason string) {
	// Unsafe PKI inputs fail the profile immediately because the evidence source
	// violates the witness safety contract.
	applyProfileState(record, StatusFail, stateFail, reason)
	record.ProfileStates = defaultProfileStates(stateFail, independenceExternal)
}

func customerPKIPassStates(policy CustomerPKIAuthorityPolicy) *ProfileStates {
	// These optimistic states are provisional: they become authoritative only
	// after validateCustomerPKIRecord completes every external-evidence check.
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
	// Authority validation rejects signer, key, policy, and revocation issues
	// before any freshness signature is allowed to establish trust.
	issue, ok := nextCustomerPKIAuthorityIssue(publicKey, policy, freshness)
	if !ok {
		return true
	}
	if issue.notAssessed {
		states.SignerAuthorityState = stateNotAssessed
		applyProfileState(record, StatusNotAssessed, stateNotAssessed, issue.reason)
		return false
	}
	customerPKIFail(record, states, issue.field, issue.reason)
	return false
}

type customerPKIAuthorityIssue struct {
	field       string
	reason      string
	notAssessed bool
	matches     bool
}

type customerPKIAuthorityCheck struct {
	field       string
	reason      string
	notAssessed bool
	matches     func() bool
}

func nextCustomerPKIAuthorityIssue(publicKey ed25519.PublicKey, policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence) (customerPKIAuthorityIssue, bool) {
	// The issue order is intentional: identity/key mismatches are stronger than
	// revocation gaps, and a required-but-absent revocation check stays
	// not_assessed instead of being reported as a false failure.
	// Each row names the profile state that should be lowered if it matches.
	// Matching stops at the first issue so reason-code precedence is stable.
	for _, check := range orderedCustomerPKIAuthorityChecks(publicKey, policy, freshness) {
		if check.matches() {
			return customerPKIAuthorityIssue{
				field:       check.field,
				reason:      check.reason,
				notAssessed: check.notAssessed,
				matches:     true,
			}, true
		}
	}
	return customerPKIAuthorityIssue{}, false
}

func orderedCustomerPKIAuthorityChecks(publicKey ed25519.PublicKey, policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence) []customerPKIAuthorityCheck {
	// Signer identity fails before key digest so callers see the authority
	// allow-list problem first.
	return []customerPKIAuthorityCheck{
		{field: "signer", reason: ReasonSignerMismatch, matches: func() bool { return customerPKISignerMismatch(policy, freshness) }},
		{field: "signer", reason: ReasonSignerMismatch, matches: func() bool { return customerPKIPublicKeyMismatch(publicKey, policy) }},
		{field: "policy", reason: ReasonPolicyMismatch, matches: func() bool { return customerPKIPolicyDigestMismatch(policy, freshness) }},
		{field: "signer", reason: ReasonRevocationNA, notAssessed: true, matches: func() bool { return customerPKIRevocationAssessmentRequired(policy) }},
		{field: "signer", reason: ReasonCertRevoked, matches: func() bool { return customerPKIRevoked(policy) }},
	}
}

func customerPKISignerMismatch(policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence) bool {
	if policy.ProfileID != "customer-pki-v1" {
		// Profile ID binds the policy to this verifier contract; a policy for a
		// different profile cannot authorize customer-pki-v1 evidence.
		return true
	}
	if policy.AllowedSignerID == "" {
		// A blank signer allow-list does not grant universal signer authority.
		return true
	}
	return policy.AllowedSignerID != freshness.SignerID
}

func customerPKIPublicKeyMismatch(publicKey ed25519.PublicKey, policy CustomerPKIAuthorityPolicy) bool {
	// Empty key digests are allowed for early integrations, but a provided
	// digest must bind exactly to the parsed public key.
	return policy.PublicKeySHA256 != "" && policy.PublicKeySHA256 != digestBytes(publicKey)
}

func customerPKIPolicyDigestMismatch(policy CustomerPKIAuthorityPolicy, freshness CustomerPKIFreshnessEvidence) bool {
	// Policy digests are optional compatibility evidence; when present, they
	// prevent freshness evidence from silently switching authority policy.
	return policy.PolicyDigest != "" && policy.PolicyDigest != freshness.PolicyDigest
}

func customerPKIRevocationAssessmentRequired(policy CustomerPKIAuthorityPolicy) bool {
	// Required revocation evidence without a state is an assessment gap. It must
	// block pass without inventing a revoked verdict.
	return policy.RevocationRequired && policy.RevocationState == ""
}

func customerPKIRevoked(policy CustomerPKIAuthorityPolicy) bool {
	// An explicit revoked state is contradictory external authority evidence and
	// therefore fails signer authority.
	return policy.RevocationState == "revoked"
}

func validateCustomerPKIFreshness(record *Record, states *ProfileStates, runsRoot, payloadDigest string, freshness CustomerPKIFreshnessEvidence) bool {
	if invalidFreshnessPayloadDigest(payloadDigest, freshness.PayloadDigest) {
		// Payload digest mismatch means the signer did not authorize the artifact
		// payload currently being assessed.
		customerPKIFail(record, states, "artifact", ReasonArtifactMismatch)
		return false
	}
	if !freshnessCurrent(freshness, time.Now().UTC()) {
		// Expired or future-dated freshness evidence cannot establish a live
		// external witness even if its signature is valid.
		customerPKIFail(record, states, "freshness", ReasonStaleFreshness)
		return false
	}
	if !runIDMatches(runsRoot, freshness.RunID) {
		// The signed run ID must resolve to a discovered local run so the
		// external freshness evidence binds to this evidence set.
		customerPKIFail(record, states, "run", ReasonRunMismatch)
		return false
	}
	return true
}

func invalidFreshnessPayloadDigest(expected, actual string) bool {
	// A digest must both match the expected payload and have strong hex shape;
	// malformed values are treated as binding failures.
	return expected != actual || !strongDigest(actual)
}

func customerPKIFail(record *Record, states *ProfileStates, field, reason string) {
	// The profile-level fail and the specific state fail are updated together so
	// downstream gates can explain which evidence class contradicted the claim.
	failCustomerPKIState(states, field)
	applyProfileState(record, StatusFail, stateFail, reason)
}

func failCustomerPKIState(states *ProfileStates, field string) {
	setter := customerPKIStateSetters[field]
	if setter == nil {
		// Unknown failure fields fall back to identity failure so the profile
		// cannot accidentally pass an unmapped Customer PKI condition.
		setter = func(states *ProfileStates) { states.IdentityState = stateFail }
	}
	setter(states)
}

var customerPKIStateSetters = map[string]func(*ProfileStates){
	"artifact":  func(states *ProfileStates) { states.ArtifactBindingState = stateFail },
	"freshness": func(states *ProfileStates) { states.FreshnessState = stateFail },
	"policy":    func(states *ProfileStates) { states.PolicyBindingState = stateFail },
	"run":       func(states *ProfileStates) { states.RunBindingState = stateFail },
	"signer":    func(states *ProfileStates) { states.SignerAuthorityState = stateFail },
}

type profileDecision struct {
	status string
	scope  string
	reason string
}

func validateCIEnvelope(kind string, envelope EnvelopeInput, current []ArtifactDigest) profileDecision {
	if envelope.ProfileID != kind+"-v1" {
		// The envelope profile ID must match the requested provider contract; a
		// different profile may use incompatible state semantics.
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonUnsupported}
	}
	if state := validateCIEnvelopeIdentity(envelope); state.reason != "" {
		return state
	}
	if state := validateCIEnvelopeStates(envelope.ProfileStates); state.reason != "" {
		return state
	}
	return validateCIEnvelopeArtifacts(envelope.RunArtifacts, current)
}

func validateCIEnvelopeIdentity(envelope EnvelopeInput) profileDecision {
	if missingEnvelopeIdentity(envelope) {
		// Missing CI identity leaves provenance unbound to the run; it cannot be
		// downgraded to a profile mismatch or treated as passing evidence.
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingIdentity}
	}
	return profileDecision{}
}

func validateCIEnvelopeArtifacts(runArtifacts, current []ArtifactDigest) profileDecision {
	if len(runArtifacts) == 0 {
		// A CI envelope without run artifact digests is not replayable evidence.
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingArtifact}
	}
	if !artifactSetsMatch(runArtifacts, current) {
		// Digest mismatch is contradictory evidence: the envelope no longer
		// describes the artifacts currently present on disk.
		return profileDecision{StatusFail, stateFail, ReasonArtifactMismatch}
	}
	return profileDecision{}
}

func missingEnvelopeIdentity(envelope EnvelopeInput) bool {
	// Commit SHA and run ID are the minimum portable identity tuple for binding
	// an envelope to source and CI execution.
	return strings.TrimSpace(envelope.Source.CommitSHA) == "" || strings.TrimSpace(envelope.CI.RunID) == ""
}

func validateCIEnvelopeStates(states ProfileStates) profileDecision {
	// State validation preserves the verifier's evidence taxonomy: missing
	// evidence remains cannot_verify while contradictory bindings fail.
	// The first matching state wins so the output reason stays deterministic and
	// mirrors the profile repair order.
	// Independence is evaluated last because it cannot rescue missing identity,
	// signer, freshness, binding, or artifact evidence.
	if decision := validateCIEnvelopeEvidenceStates(states); decision.reason != "" {
		return decision
	}
	if decision := validateCIEnvelopeBindingStates(states); decision.reason != "" {
		return decision
	}
	if !ciEnvelopeIndependent(states) {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonEnvOnly}
	}
	return profileDecision{}
}

func validateCIEnvelopeEvidenceStates(states ProfileStates) profileDecision {
	// Evidence states are checked before binding states because missing identity,
	// signer authority, or freshness means the envelope cannot yet be trusted
	// enough to interpret source, run, policy, or artifact claims.
	if decision := validateCIEnvelopeIdentityAuthorityStates(states); decision.reason != "" {
		return decision
	}
	return validateCIEnvelopeFreshnessState(states)
}

func validateCIEnvelopeIdentityAuthorityStates(states ProfileStates) profileDecision {
	// Missing identity or signer authority means the envelope never reached a
	// replayable evidence boundary, so the verdict stays cannot_verify.
	if states.IdentityState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingIdentity}
	}
	if states.SignerAuthorityState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingSigner}
	}
	return profileDecision{}
}

func validateCIEnvelopeFreshnessState(states ProfileStates) profileDecision {
	// Explicit stale freshness is contradictory evidence and fails; any other
	// non-pass freshness state is missing evidence.
	if states.FreshnessState == stateFail {
		return profileDecision{StatusFail, stateFail, ReasonStaleFreshness}
	}
	if states.FreshnessState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingFreshness}
	}
	return profileDecision{}
}

func validateCIEnvelopeBindingStates(states ProfileStates) profileDecision {
	// Binding checks run only after core evidence is present; at that point
	// source/run mismatches are failures, while policy/artifact gaps remain
	// cannot_verify because the required comparison evidence is incomplete.
	if decision := validateCIEnvelopeRunSourceStates(states); decision.reason != "" {
		return decision
	}
	return validateCIEnvelopePolicyArtifactStates(states)
}

func validateCIEnvelopeRunSourceStates(states ProfileStates) profileDecision {
	// Source and run mismatches contradict the claimed execution context, so
	// they fail rather than remaining open as missing evidence.
	if states.SourceBindingState != statePass {
		return profileDecision{StatusFail, stateFail, ReasonSourceMismatch}
	}
	if states.RunBindingState != statePass {
		return profileDecision{StatusFail, stateFail, ReasonRunMismatch}
	}
	return profileDecision{}
}

func validateCIEnvelopePolicyArtifactStates(states ProfileStates) profileDecision {
	// Policy and artifact binding gaps leave the envelope unverifiable because
	// the local verifier cannot prove the claimed evidence set.
	if states.PolicyBindingState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonPolicyMissing}
	}
	if states.ArtifactBindingState != statePass {
		return profileDecision{StatusCannotVerify, stateCannotVerify, ReasonMissingArtifact}
	}
	return profileDecision{}
}

func ciEnvelopeIndependent(states ProfileStates) bool {
	return states.IndependenceState == independenceCIJob || states.IndependenceState == independenceExternal
}

func baseRecord(kind string) Record {
	// Base records start as cannot_verify/local_observed until profile-specific
	// evidence raises or fails the established trust scope.
	// Empty artifact slices are intentional; nil would make generated records
	// shape-shift across profiles.
	// OutputSafety starts optimistic but is recalculated before write.
	// GeneratedAt is created locally and does not claim external witness time.
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
		OutputSafety:          passingOutputSafety(),
	}
}

func passingOutputSafety() *OutputSafety {
	// The pass state means the known unsafe classes were checked for absence; it
	// does not make any claim about the profile's underlying trust verdict.
	return &OutputSafety{
		State:                 statePass,
		VerifiedAbsentClasses: safetyClasses,
	}
}

func applyProfileState(record *Record, status, scope, reason string) {
	// Keep human status, trust scope, and machine-readable reason code aligned
	// whenever a profile decision is applied.
	record.Status = status
	record.TrustScope = scope
	record.EstablishedTrustScope = scope
	record.Reason = reason
	record.ReasonCodes = []string{reason}
}

func finalizeRecordForWrite(record Record) Record {
	// Scan a copy without OutputSafety so the safety attestation cannot satisfy
	// itself or hide unsafe material in the rest of the record.
	// The original record is only returned when the serialized form is free of
	// known unsafe output classes.
	// Serialization failures are treated as unsafe output because no reviewable
	// witness artifact can be produced from that record.
	// A replacement failure record avoids partial redaction of untrusted payload
	// content.
	raw, ok := outputSafetyScanBytes(record)
	if !ok {
		applyProfileState(&record, StatusFail, stateFail, ReasonUnsafeOutput)
		return record
	}
	if !forbiddenOutputPresent(raw) {
		// Passing output safety states only that known unsafe classes were absent
		// from the serialized record; it does not upgrade the profile verdict.
		applyOutputSafetyPass(&record)
		return record
	}
	return unsafeOutputRecord(record.Kind)
}

func outputSafetyScanBytes(record Record) ([]byte, bool) {
	// Scan bytes exclude OutputSafety so the attestation cannot recursively
	// satisfy the safety check that is about to be recorded.
	scanRecord := record
	scanRecord.OutputSafety = nil
	raw, err := json.Marshal(scanRecord)
	return raw, err == nil
}

func applyOutputSafetyPass(record *Record) {
	// Preserve the caller's record while refreshing only the output-safety
	// attestation after the serialized trust payload has been scanned.
	if record.OutputSafety == nil {
		record.OutputSafety = &OutputSafety{}
	}
	record.OutputSafety.State = statePass
	record.OutputSafety.VerifiedAbsentClasses = safetyClasses
}

func unsafeOutputRecord(kind string) Record {
	safe := baseRecord(kind)
	// Unsafe output candidates are replaced by a minimal failure record so the
	// published artifact carries the verdict without carrying the unsafe data.
	safe.ProfileStates = defaultProfileStates(stateFail, "cannot_verify")
	safe.OutputSafety = &OutputSafety{State: stateFail, VerifiedAbsentClasses: safetyClasses}
	applyProfileState(&safe, StatusFail, stateFail, ReasonUnsafeOutput)
	return safe
}

func forbiddenOutputPresent(raw []byte) bool {
	// Output safety checks marker families before publication; unsafe content is
	// never echoed in the resulting failure record.
	// Secret-like structured markers are checked before general lowercase marker
	// scanning.
	// Marker matching is deliberately string-based so it works on serialized JSON
	// without schema-specific traversal.
	if containsSecretLike(raw) {
		return true
	}
	// Marker checks catch unsafe payload classes even when they are not shaped
	// like standard tokens or private keys.
	text := strings.ToLower(string(raw))
	for _, marker := range outputSafetyMarkers {
		if strings.Contains(text, marker) {
			return true
		}
	}
	return jwtLike(text)
}

func defaultProfileStates(state, independence string) *ProfileStates {
	// Defaults make unassessed profile state explicit; callers must override
	// individual states when they have evidence for a narrower verdict.
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
		// Count differences are binding failures even if all shared paths match.
		return false
	}
	byPath := artifactDigestsByPath(current)
	for _, artifact := range expected {
		if byPath[artifact.Path] != artifact.SHA256 {
			return false
		}
	}
	return true
}

func artifactDigestsByPath(artifacts []ArtifactDigest) map[string]string {
	byPath := map[string]string{}
	for _, artifact := range artifacts {
		// Path is the binding key because witness artifacts are compared against
		// the selected run/report file set, not just a bag of digests.
		byPath[artifact.Path] = artifact.SHA256
	}
	return byPath
}

func runIDMatches(runsRoot, witnessRunID string) bool {
	if strings.TrimSpace(witnessRunID) == "" {
		// Empty run IDs cannot bind external or envelope evidence to a run.
		return false
	}
	runIDs, err := runIDsFromRoot(runsRoot)
	if err != nil || len(runIDs) == 0 {
		return false
	}
	return containsRunID(runIDs, witnessRunID)
}

func containsRunID(runIDs []string, witnessRunID string) bool {
	// Exact run ID matching prevents envelope or freshness evidence from binding
	// to a sibling run.
	for _, runID := range runIDs {
		if runID == witnessRunID {
			return true
		}
	}
	return false
}

func runIDsFromRoot(runsRoot string) ([]string, error) {
	// Demo discovery defines the replayable run set for witness binding.
	runDirs, err := demo.DiscoverRunDirs(runsRoot)
	if err != nil {
		return nil, err
	}
	return runIDsFromDirs(runDirs)
}

func runIDsFromDirs(runDirs []string) ([]string, error) {
	runIDs := make([]string, 0, len(runDirs))
	for _, runDir := range runDirs {
		// Skip run directories that lack an ID but fail on unreadable or
		// malformed run.json, preserving absent versus bad evidence.
		runID, ok, err := nonEmptyRunIDFromDir(runDir)
		if err != nil {
			return nil, err
		}
		if ok {
			runIDs = append(runIDs, runID)
		}
	}
	return runIDs, nil
}

func nonEmptyRunIDFromDir(runDir string) (string, bool, error) {
	runID, err := runIDFromDir(runDir)
	if err != nil {
		return "", false, err
	}
	// Empty run IDs are skipped instead of being treated as a wildcard match.
	return runID, runID != "", nil
}

func runIDFromDir(runDir string) (string, error) {
	// Accept both run_id and legacy id to keep old evidence replayable without
	// weakening the requirement that some run identity be present.
	// The returned ID is trimmed so whitespace-only legacy IDs do not bind.
	// Malformed run.json is an error rather than a non-match because it blocks
	// reliable run binding.
	raw, err := os.ReadFile(filepath.Join(runDir, "run.json"))
	if err != nil {
		return "", err
	}
	var payload struct {
		RunID string `json:"run_id"`
		ID    string `json:"id"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", err
	}
	if strings.TrimSpace(payload.RunID) != "" {
		return payload.RunID, nil
	}
	return strings.TrimSpace(payload.ID), nil
}

func ambientCIEnvPresent(kind string) bool {
	// Ambient variables explain why an envelope is required, but this helper
	// never upgrades trust by itself.
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
	// Required input reporting is stable and sorted so the missing-evidence
	// surface can be compared across runs.
	missing := []string{}
	for name, value := range requiredCustomerPKIInputs(opts) {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}
	return appendMissingCustomerPKIKeyInput(missing, opts)
}

func appendMissingCustomerPKIKeyInput(missing []string, opts ProfileOptions) []string {
	// Public key and public certificate are alternatives, so the missing-input
	// message reports them as one choice.
	if strings.TrimSpace(opts.CustomerPKIPublicCert) == "" && strings.TrimSpace(opts.CustomerPKIPublicKey) == "" {
		missing = append(missing, "--customer-pki-public-cert|--customer-pki-public-key")
	}
	sort.Strings(missing)
	return missing
}

func requiredCustomerPKIInputs(opts ProfileOptions) map[string]string {
	// These three external evidence files are mandatory; the public key/cert
	// alternative is checked separately because either form can carry the key.
	return map[string]string{
		"--customer-pki-authority-policy":   opts.CustomerPKIAuthorityPolicy,
		"--customer-pki-payload-digest":     opts.CustomerPKIPayloadDigest,
		"--customer-pki-freshness-evidence": opts.CustomerPKIFreshness,
	}
}

func readSafeJSON(path string, target any) error {
	if unsafeInputPath(path) {
		// Unsafe paths are rejected before reads so callers cannot smuggle
		// remote URLs, traversal, symlinks, or private-key filenames as evidence.
		return fmt.Errorf("unsafe input path")
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if containsSecretLike(raw) {
		// JSON evidence files must not carry secrets; a secret-like input cannot
		// be safely copied into witness outputs.
		return fmt.Errorf("unsafe input content")
	}
	return json.Unmarshal(raw, target)
}

func unsafeInputPath(path string) bool {
	// Reject unsafe text first, then reject symlinks so Customer PKI inputs stay
	// source-bound to explicit files.
	lower := strings.ToLower(filepath.ToSlash(path))
	if unsafeInputPathText(path, lower) {
		return true
	}
	return inputPathIsSymlink(path)
}

func unsafeInputPathText(path, lower string) bool {
	return emptyOrNULPath(path) ||
		unsafeLowerInputPathText(lower)
}

func emptyOrNULPath(path string) bool {
	return strings.TrimSpace(path) == "" || strings.Contains(path, "\x00")
}

func unsafeLowerInputPathText(lower string) bool {
	return strings.Contains(lower, "://") || strings.Contains(lower, "..") || strings.Contains(lower, "private.key")
}

func inputPathIsSymlink(path string) bool {
	info, err := os.Lstat(path)
	return err == nil && info.Mode()&os.ModeSymlink != 0
}

func containsSecretLike(raw []byte) bool {
	// This detector is intentionally conservative. It protects witness outputs
	// from known secret shapes without trying to classify or print the secret.
	// Unknown content is not declared safe by this helper; it only blocks known
	// high-risk markers.
	// The final JWT-shape check covers tokens that lack a provider prefix.
	// All matching is done on lowercase text so provider marker casing cannot
	// bypass the deny-list.
	lower := strings.ToLower(string(raw))
	for _, marker := range secretSafetyMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return jwtLike(lower)
}

func jwtLike(text string) bool {
	// JWT detection runs on split fields only; callers get a boolean decision,
	// never decoded token claims or token material.
	for _, field := range jwtCandidateFields(text) {
		if jwtCandidate(field) {
			return true
		}
	}
	return false
}

func jwtCandidateFields(text string) []string {
	// Split on common JSON and prose separators so token-shaped substrings are
	// checked without decoding them.
	return strings.FieldsFunc(text, func(r rune) bool {
		switch r {
		case '"', '\'', ' ', '\n', '\t', '\r', ',', ':', '{', '}', '[', ']', '(', ')':
			return true
		default:
			return false
		}
	})
}

func jwtCandidate(field string) bool {
	parts := strings.Split(field, ".")
	// JWT-like values are refused by shape only; this catches common bearer
	// token leaks without logging or decoding token material.
	return len(parts) == 3 && strings.HasPrefix(parts[0], "eyj") && len(parts[1]) >= 8 && len(parts[2]) >= 8
}

func unsafeEnvelopeFields(envelope EnvelopeInput) bool {
	// Envelope safety covers both identity metadata and artifact references
	// because either can leak credentials or private paths.
	return unsafeEnvelopeScalarFields(envelope) || unsafeEnvelopeArtifactFields(envelope)
}

func unsafeEnvelopeScalarFields(envelope EnvelopeInput) bool {
	// Scalar provider fields are copied into witness records, so they must be
	// safe before the envelope can influence the output record.
	// Structured artifact references are checked separately because their path
	// and digest rules differ.
	// The list mirrors fields copied by applyCIEnvelopeRecordValues.
	// Any unsafe scalar blocks the entire envelope before profile states are
	// trusted.
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
	return false
}

func unsafeEnvelopeArtifactFields(envelope EnvelopeInput) bool {
	for _, artifact := range append(envelope.RunArtifacts, envelope.ReportArtifacts...) {
		// Artifact paths and digests are persisted verbatim in the witness output.
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
		// Secret-shaped scalar values fail before broader personal/path checks.
		return true
	}
	lower := strings.ToLower(value)
	return strings.Contains(lower, "/private/") || strings.Contains(lower, "@")
}

func privateKeyInput(path string) bool {
	if strings.TrimSpace(path) == "" {
		return false
	}
	// This best-effort preflight catches accidental private-key files before
	// parsing public-key or certificate material.
	raw, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return containsSecretLike(raw)
}

func loadCustomerPublicKey(opts ProfileOptions) (ed25519.PublicKey, error) {
	// Load the chosen public trust anchor only after input-path and private-key
	// checks have already rejected unsafe material.
	raw, err := os.ReadFile(customerPKIPublicKeyPath(opts))
	if err != nil {
		return nil, err
	}
	return parseCustomerPublicKeyPEM(raw)
}

func customerPKIPublicKeyPath(opts ProfileOptions) string {
	if opts.CustomerPKIPublicKey == "" {
		// Certificate input is the fallback trust anchor only when a direct
		// public key path was not supplied.
		return opts.CustomerPKIPublicCert
	}
	return opts.CustomerPKIPublicKey
}

func parseCustomerPublicKeyPEM(raw []byte) (ed25519.PublicKey, error) {
	if containsSecretLike(raw) {
		// Public trust anchors must not include signing secrets.
		return nil, errors.New("private key input rejected")
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("public key or certificate PEM is required")
	}
	if block.Type == "CERTIFICATE" {
		// Certificates are accepted only as carriers for an Ed25519 public key;
		// revocation and custody are evaluated through explicit policy fields.
		return parseCertificatePublicKey(block.Bytes)
	}
	return parsePKIXPublicKey(block.Bytes)
}

func parseCertificatePublicKey(raw []byte) (ed25519.PublicKey, error) {
	cert, err := x509.ParseCertificate(raw)
	if err != nil {
		return nil, err
	}
	// Only Ed25519 keys are supported so signature verification has one narrow
	// algorithm contract.
	key, ok := cert.PublicKey.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("certificate must contain ed25519 public key")
	}
	return key, nil
}

func parsePKIXPublicKey(raw []byte) (ed25519.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(raw)
	if err != nil {
		return nil, err
	}
	// Reject other public-key algorithms rather than silently changing the
	// signature contract.
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
	// Signature verification covers the canonical freshness payload assembled
	// below, binding payload digest, run ID, policy digest, signer, time, and nonce.
	return ed25519.Verify(publicKey, []byte(freshnessPayload(evidence)), signature)
}

func freshnessPayload(evidence CustomerPKIFreshnessEvidence) string {
	// Newline joining gives the signed payload a stable field order without
	// depending on JSON map ordering or formatting.
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
		// A small skew allowance avoids rejecting near-current evidence while
		// still blocking far-future timestamps.
		return false
	}
	return freshnessValidUntilCurrent(evidence.ValidUntil, now)
}

func freshnessValidUntilCurrent(validUntilText string, now time.Time) bool {
	if validUntilText == "" {
		// Missing expiry is allowed for compatibility, but issued_at is still
		// mandatory and checked by freshnessCurrent.
		return true
	}
	validUntil, err := time.Parse(time.RFC3339, validUntilText)
	return err == nil && !validUntil.Before(now)
}

func digestBytes(data []byte) string {
	// Policy key binding uses the digest of the parsed public key bytes, not the
	// original PEM text, so formatting differences do not affect authority.
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func strongDigest(value string) bool {
	if len(value) < 64 {
		// Customer PKI freshness must bind to at least a SHA-256-sized hex digest.
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		// Blank profile fields inherit the caller's explicit fallback instead
		// of becoming empty trust-context claims.
		return fallback
	}
	return value
}
