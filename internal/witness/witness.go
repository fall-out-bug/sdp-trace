package witness

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

const (
	KindGitHubActions       = "github-actions"
	KindGitLabCI            = "gitlab-ci"
	KindBuildkite           = "buildkite"
	KindCustomerPKI         = "customer-pki"
	StatusPass              = "pass"
	StatusFail              = "fail"
	StatusCannotVerify      = "cannot_verify"
	StatusNotAssessed       = "not_assessed"
	TrustScopeCIWitnessed   = "ci_witnessed"
	TrustScopeLocalObserved = "local_observed"
	TrustScopeExternal      = "external_witnessed"
	ReasonCIIdentityPresent = "ci_identity_present"
	ReasonProfileVerified   = "witness_profile_verified"
	ReasonMissingCIIdentity = "missing_ci_identity"
	ReasonMissingCIOIDC     = "missing_ci_oidc"
	ReasonInvalidCIOIDC     = "invalid_ci_oidc"
	ReasonEnvOnly           = "witness_environment_only_insufficient"
	ReasonMissingIdentity   = "witness_identity_missing"
	ReasonMissingSigner     = "witness_signer_authority_missing"
	ReasonMissingFreshness  = "witness_freshness_missing"
	ReasonStaleFreshness    = "witness_freshness_stale"
	ReasonMissingArtifact   = "witness_artifact_digest_missing"
	ReasonArtifactMismatch  = "witness_artifact_digest_mismatch"
	ReasonIdentityMismatch  = "witness_identity_mismatch"
	ReasonSourceMissing     = "witness_source_binding_missing"
	ReasonSourceMismatch    = "witness_source_mismatch"
	ReasonRunMissing        = "witness_run_binding_missing"
	ReasonRunMismatch       = "witness_run_mismatch"
	ReasonPolicyMissing     = "witness_policy_binding_missing"
	ReasonPolicyMismatch    = "witness_policy_mismatch"
	ReasonSignerMismatch    = "witness_signer_mismatch"
	ReasonUnsupported       = "witness_unsupported_profile"
	ReasonUnsafeOutput      = "witness_unsafe_output_candidate"
	ReasonPrivateKeyInput   = "witness_private_key_input_rejected"
	ReasonRevocationNA      = "witness_revocation_not_assessed"
	ReasonCertRevoked       = "witness_certificate_revoked"
	ReasonKeyCustodyNA      = "witness_key_custody_not_assessed"
	ReasonMalformedInput    = "witness_malformed_input"
	githubOIDCIssuer        = "https://token.actions.githubusercontent.com"
)

type Record struct {
	SchemaVersion         string           `json:"schema_version,omitempty"`
	Kind                  string           `json:"kind"`
	ProfileID             string           `json:"profile_id,omitempty"`
	ProfileVersion        string           `json:"profile_version,omitempty"`
	ProviderKind          string           `json:"provider_kind,omitempty"`
	Status                string           `json:"status"`
	TrustScope            string           `json:"trust_scope"`
	RequestedTrustScope   string           `json:"requested_trust_scope,omitempty"`
	EstablishedTrustScope string           `json:"established_trust_scope,omitempty"`
	Reason                string           `json:"reason"`
	ReasonCodes           []string         `json:"reason_codes,omitempty"`
	GeneratedAt           string           `json:"generated_at"`
	MissingIdentityFields []string         `json:"missing_identity_fields,omitempty"`
	Source                SourceIdentity   `json:"source"`
	CI                    CIIdentity       `json:"ci"`
	OIDC                  *OIDCClaims      `json:"oidc,omitempty"`
	RunArtifacts          []ArtifactDigest `json:"run_artifacts"`
	ReportArtifacts       []ArtifactDigest `json:"report_artifacts"`
	ProfileStates         *ProfileStates   `json:"profile_states"`
	OutputSafety          *OutputSafety    `json:"output_safety"`
}

type SourceIdentity struct {
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	CommitSHA  string `json:"commit_sha"`
}

type CIIdentity struct {
	Provider   string `json:"provider"`
	ServerURL  string `json:"server_url"`
	Workflow   string `json:"workflow"`
	Job        string `json:"job"`
	RunID      string `json:"run_id"`
	RunAttempt string `json:"run_attempt"`
	Actor      string `json:"actor"`
}

type ArtifactDigest struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type ProfileStates struct {
	IdentityState        string `json:"identity_state"`
	SignerAuthorityState string `json:"signer_authority_state"`
	FreshnessState       string `json:"freshness_state"`
	ArtifactBindingState string `json:"artifact_binding_state"`
	SourceBindingState   string `json:"source_binding_state"`
	RunBindingState      string `json:"run_binding_state"`
	PolicyBindingState   string `json:"policy_binding_state"`
	IndependenceState    string `json:"independence_state"`
	KeyCustodyState      string `json:"key_custody_state,omitempty"`
}

type OutputSafety struct {
	State                 string   `json:"state"`
	VerifiedAbsentClasses []string `json:"verified_absent_classes"`
}

type OIDCClaims struct {
	Issuer     string `json:"issuer"`
	Subject    string `json:"subject"`
	Audience   string `json:"audience"`
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	SHA        string `json:"sha"`
}

type TokenFetcher func(env map[string]string) (string, error)

func BuildGitHubActions(runsRoot, reportDir string, env map[string]string) (Record, error) {
	return BuildGitHubActionsWithFetcher(runsRoot, reportDir, env, FetchGitHubOIDCToken)
}

func BuildGitHubActionsWithFetcher(runsRoot, reportDir string, env map[string]string, fetcher TokenFetcher) (Record, error) {
	record := baseRecord(KindGitHubActions)
	record.Source = SourceIdentity{
		Repository: env["GITHUB_REPOSITORY"],
		Ref:        env["GITHUB_REF"],
		CommitSHA:  env["GITHUB_SHA"],
	}
	record.CI = CIIdentity{
		Provider:   KindGitHubActions,
		ServerURL:  env["GITHUB_SERVER_URL"],
		Workflow:   env["GITHUB_WORKFLOW"],
		Job:        env["GITHUB_JOB"],
		RunID:      env["GITHUB_RUN_ID"],
		RunAttempt: env["GITHUB_RUN_ATTEMPT"],
		Actor:      env["GITHUB_ACTOR"],
	}

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

	missing := missingGitHubIdentity(env)
	if len(missing) > 0 {
		applyProfileState(&record, StatusCannotVerify, stateCannotVerify, ReasonMissingCIIdentity)
		record.TrustScope = TrustScopeLocalObserved
		record.ProfileStates = defaultProfileStates(stateCannotVerify, independenceSameJob)
		record.MissingIdentityFields = missing
		return record, nil
	}
	oidcMissing := missingGitHubOIDC(env)
	if len(oidcMissing) > 0 {
		applyProfileState(&record, StatusCannotVerify, stateCannotVerify, ReasonMissingCIOIDC)
		record.TrustScope = TrustScopeLocalObserved
		record.ProfileStates = defaultProfileStates(stateCannotVerify, independenceCIJob)
		record.MissingIdentityFields = oidcMissing
		return record, nil
	}
	token, err := fetcher(env)
	if err != nil {
		applyProfileState(&record, StatusCannotVerify, stateCannotVerify, ReasonInvalidCIOIDC)
		record.TrustScope = TrustScopeLocalObserved
		record.ProfileStates = defaultProfileStates(stateCannotVerify, independenceCIJob)
		return record, nil
	}
	claims, err := parseOIDCClaims(token)
	if err != nil || !claimsMatchEnvironment(claims, env) {
		applyProfileState(&record, StatusCannotVerify, stateCannotVerify, ReasonInvalidCIOIDC)
		record.TrustScope = TrustScopeLocalObserved
		record.ProfileStates = defaultProfileStates(stateCannotVerify, independenceCIJob)
		return record, nil
	}
	record.OIDC = &claims
	record.Status = StatusPass
	record.TrustScope = TrustScopeCIWitnessed
	record.EstablishedTrustScope = TrustScopeCIWitnessed
	record.Reason = ReasonCIIdentityPresent
	record.ReasonCodes = []string{ReasonCIIdentityPresent}
	record.ProfileStates = defaultProfileStates(statePass, independenceCIJob)
	return record, nil
}

func WriteGitHubActions(outPath, runsRoot, reportDir string, env map[string]string) (Record, error) {
	if strings.TrimSpace(outPath) == "" {
		return Record{}, errors.New("witness requires --out <file>")
	}
	record, err := BuildGitHubActions(runsRoot, reportDir, env)
	if err != nil {
		return Record{}, err
	}
	record = finalizeRecordForWrite(record)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return Record{}, err
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return Record{}, err
	}
	if err := os.WriteFile(outPath, append(data, '\n'), 0o644); err != nil {
		return Record{}, err
	}
	return record, nil
}

func Load(path string) (Record, error) {
	var record Record
	data, err := os.ReadFile(path)
	if err != nil {
		return Record{}, err
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return Record{}, err
	}
	return record, nil
}

func IsPassingCI(record Record) bool {
	return record.Kind == KindGitHubActions &&
		record.Status == StatusPass &&
		record.TrustScope == TrustScopeCIWitnessed
}

func EnvironmentFromOS() map[string]string {
	env := map[string]string{}
	for _, entry := range os.Environ() {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			env[key] = value
		}
	}
	return env
}

func missingGitHubIdentity(env map[string]string) []string {
	required := []string{
		"GITHUB_ACTIONS",
		"GITHUB_SHA",
		"GITHUB_RUN_ID",
		"GITHUB_RUN_ATTEMPT",
		"GITHUB_WORKFLOW",
		"GITHUB_JOB",
		"GITHUB_ACTOR",
		"GITHUB_REPOSITORY",
		"GITHUB_REF",
		"GITHUB_SERVER_URL",
	}
	missing := make([]string, 0)
	for _, key := range required {
		if strings.TrimSpace(env[key]) == "" {
			missing = append(missing, key)
		}
	}
	if env["GITHUB_ACTIONS"] != "" && env["GITHUB_ACTIONS"] != "true" {
		missing = append(missing, "GITHUB_ACTIONS=true")
	}
	sort.Strings(missing)
	return missing
}

func missingGitHubOIDC(env map[string]string) []string {
	required := []string{"ACTIONS_ID_TOKEN_REQUEST_URL", "ACTIONS_ID_TOKEN_REQUEST_TOKEN"}
	missing := make([]string, 0)
	for _, key := range required {
		if strings.TrimSpace(env[key]) == "" {
			missing = append(missing, key)
		}
	}
	sort.Strings(missing)
	return missing
}

func FetchGitHubOIDCToken(env map[string]string) (string, error) {
	requestURL, err := url.Parse(env["ACTIONS_ID_TOKEN_REQUEST_URL"])
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(requestURL.Host, "actions.githubusercontent.com") {
		return "", fmt.Errorf("unexpected oidc request host: %s", requestURL.Host)
	}
	query := requestURL.Query()
	query.Set("audience", "sdp-trace")
	requestURL.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "bearer "+env["ACTIONS_ID_TOKEN_REQUEST_TOKEN"])
	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("oidc token request returned %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var payload struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	if payload.Value == "" {
		return "", errors.New("oidc token response missing value")
	}
	return payload.Value, nil
}

func parseOIDCClaims(token string) (OIDCClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return OIDCClaims{}, errors.New("invalid jwt shape")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return OIDCClaims{}, err
	}
	var raw struct {
		Issuer     string `json:"iss"`
		Subject    string `json:"sub"`
		Audience   any    `json:"aud"`
		Repository string `json:"repository"`
		Ref        string `json:"ref"`
		SHA        string `json:"sha"`
	}
	if err := json.Unmarshal(payload, &raw); err != nil {
		return OIDCClaims{}, err
	}
	return OIDCClaims{
		Issuer:     raw.Issuer,
		Subject:    raw.Subject,
		Audience:   audienceString(raw.Audience),
		Repository: raw.Repository,
		Ref:        raw.Ref,
		SHA:        raw.SHA,
	}, nil
}

func audienceString(value any) string {
	switch typed := value.(type) {
	case string:
		return typed
	case []any:
		parts := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok {
				parts = append(parts, text)
			}
		}
		return strings.Join(parts, ",")
	default:
		return ""
	}
}

func claimsMatchEnvironment(claims OIDCClaims, env map[string]string) bool {
	return claims.Issuer == githubOIDCIssuer &&
		claims.Audience == "sdp-trace" &&
		claims.Repository == env["GITHUB_REPOSITORY"] &&
		claims.Ref == env["GITHUB_REF"] &&
		claims.SHA == env["GITHUB_SHA"]
}

func hashRunArtifacts(runsRoot string) ([]ArtifactDigest, error) {
	runDirs, err := demo.DiscoverRunDirs(runsRoot)
	if err != nil {
		return nil, err
	}
	artifacts := make([]ArtifactDigest, 0, len(runDirs))
	for _, runDir := range runDirs {
		digest, err := hashFile(filepath.Join(runDir, "run.json"))
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, ArtifactDigest{
			Path:   filepath.ToSlash(filepath.Join(filepath.Base(runDir), "run.json")),
			SHA256: digest,
		})
	}
	return artifacts, nil
}

func hashReportArtifacts(reportDir string) ([]ArtifactDigest, error) {
	entries, err := os.ReadDir(reportDir)
	if err != nil {
		return nil, err
	}
	artifacts := make([]ArtifactDigest, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == "ci-witness.json" {
			continue
		}
		path := filepath.Join(reportDir, entry.Name())
		digest, err := hashFile(path)
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, ArtifactDigest{
			Path:   filepath.ToSlash(entry.Name()),
			SHA256: digest,
		})
	}
	sort.Slice(artifacts, func(i, j int) bool {
		return artifacts[i].Path < artifacts[j].Path
	})
	return artifacts, nil
}

func hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("hash %s: %w", path, err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
