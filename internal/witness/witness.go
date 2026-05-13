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

type rawOIDCClaims struct {
	Issuer     string `json:"iss"`
	Subject    string `json:"sub"`
	Audience   any    `json:"aud"`
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	SHA        string `json:"sha"`
}

type TokenFetcher func(env map[string]string) (string, error)

func BuildGitHubActions(runsRoot, reportDir string, env map[string]string) (Record, error) {
	return BuildGitHubActionsWithFetcher(runsRoot, reportDir, env, FetchGitHubOIDCToken)
}

func BuildGitHubActionsWithFetcher(runsRoot, reportDir string, env map[string]string, fetcher TokenFetcher) (Record, error) {
	// GitHub Actions witness generation starts as a local record and upgrades
	// only after artifact hashing, required environment fields, and OIDC claims
	// all bind to the same execution.
	// The injectable fetcher keeps network trust decisions testable without
	// weakening production request validation.
	// Source and CI fields are copied before validation so failure records remain
	// useful for explaining which identity material was present.
	record := baseRecord(KindGitHubActions)
	record.Source = githubSourceIdentity(env)
	record.CI = githubCIIdentity(env)

	if err := hydrateGitHubArtifacts(&record, runsRoot, reportDir); err != nil {
		return Record{}, err
	}
	// Identity and OIDC checks are separate so missing provider variables remain
	// distinguishable from invalid attestation evidence.
	if record, blocked := handleGitHubIdentityChecks(record, env); blocked {
		return record, nil
	}
	if record, blocked := handleGitHubOIDCChecks(record, env, fetcher); blocked {
		return record, nil
	}
	// Only after local artifacts, environment identity, and OIDC claims align
	// can the record be promoted to ci_witnessed.
	return passingGitHubRecord(record), nil
}

func githubSourceIdentity(env map[string]string) SourceIdentity {
	// GitHub source identity is an environment snapshot until OIDC claim binding
	// below promotes it across the CI trust boundary.
	return SourceIdentity{
		Repository: env["GITHUB_REPOSITORY"],
		Ref:        env["GITHUB_REF"],
		CommitSHA:  env["GITHUB_SHA"],
	}
}

func githubCIIdentity(env map[string]string) CIIdentity {
	// CI identity fields remain explanatory local evidence until the matching
	// OIDC token has been fetched and parsed.
	return CIIdentity{
		Provider:   KindGitHubActions,
		ServerURL:  env["GITHUB_SERVER_URL"],
		Workflow:   env["GITHUB_WORKFLOW"],
		Job:        env["GITHUB_JOB"],
		RunID:      env["GITHUB_RUN_ID"],
		RunAttempt: env["GITHUB_RUN_ATTEMPT"],
		Actor:      env["GITHUB_ACTOR"],
	}
}

func passingGitHubRecord(record Record) Record {
	// This is the sole promotion point from local_observed to ci_witnessed for
	// GitHub Actions records.
	record.Status = StatusPass
	record.TrustScope = TrustScopeCIWitnessed
	record.EstablishedTrustScope = TrustScopeCIWitnessed
	record.Reason = ReasonCIIdentityPresent
	record.ReasonCodes = []string{ReasonCIIdentityPresent}
	record.ProfileStates = defaultProfileStates(statePass, independenceCIJob)
	return record
}

func hydrateGitHubArtifacts(record *Record, runsRoot, reportDir string) error {
	// Artifact hashes are local evidence that the OIDC-backed identity must bind
	// to; they do not by themselves establish CI witness trust.
	// Report artifacts are optional because some callers only need run-output
	// binding.
	runArtifacts, err := hashRunArtifacts(runsRoot)
	if err != nil {
		return err
	}
	record.RunArtifacts = runArtifacts
	if reportDir != "" {
		reportArtifacts, err := hashReportArtifacts(reportDir)
		if err != nil {
			return err
		}
		record.ReportArtifacts = reportArtifacts
	}
	return nil
}

func handleGitHubIdentityChecks(record Record, env map[string]string) (Record, bool) {
	missing := missingGitHubIdentity(env)
	if len(missing) > 0 {
		// Missing GitHub identity keeps the record local_observed because the run
		// cannot be tied to a provider execution.
		return applyGitHubFailure(record, ReasonMissingCIIdentity, independenceSameJob, missing), true
	}
	oidcMissing := missingGitHubOIDC(env)
	if len(oidcMissing) > 0 {
		// A GitHub job without OIDC token provisioning cannot produce the
		// attested CI identity required for ci_witnessed scope.
		return applyGitHubFailure(record, ReasonMissingCIOIDC, independenceCIJob, oidcMissing), true
	}
	return record, false
}

func handleGitHubOIDCChecks(record Record, env map[string]string, fetcher TokenFetcher) (Record, bool) {
	token, err := fetcher(env)
	if err != nil {
		// Fetch failures leave the provider identity unverified; the token value
		// is intentionally not recorded in the witness output.
		return applyGitHubFailure(record, ReasonInvalidCIOIDC, independenceCIJob, nil), true
	}
	claims, err := parseOIDCClaims(token)
	if err != nil || !claimsMatchEnvironment(claims, env) {
		// Parsed claims must match both issuer/audience and source identity
		// before they can raise the trust scope.
		return applyGitHubFailure(record, ReasonInvalidCIOIDC, independenceCIJob, nil), true
	}
	record.OIDC = &claims
	return record, false
}

// applyGitHubFailure transitions the record to cannot_verify state and
// downgrades the trust scope from ci_witnessed to local_observed. This
// boundary change reflects that we are now making claims about locally
// captured environment state rather than cryptographically attested CI identity.
func applyGitHubFailure(record Record, reason, independence string, missing []string) Record {
	// Failure records preserve whatever local artifact evidence was already
	// collected, but lower the trust scope before returning.
	applyProfileState(&record, StatusCannotVerify, stateCannotVerify, reason)
	record.TrustScope = TrustScopeLocalObserved
	record.ProfileStates = defaultProfileStates(stateCannotVerify, independence)
	record.MissingIdentityFields = missing
	return record
}

func WriteGitHubActions(outPath, runsRoot, reportDir string, env map[string]string) (Record, error) {
	// Writing a witness is explicit: no output path means no durable claim is
	// produced.
	if strings.TrimSpace(outPath) == "" {
		return Record{}, errors.New("witness requires --out <file>")
	}
	record, err := BuildGitHubActions(runsRoot, reportDir, env)
	if err != nil {
		return Record{}, err
	}
	// The GitHub-specific writer uses the same output-safety finalizer as
	// portable profiles before creating a durable witness artifact.
	record = finalizeRecordForWrite(record)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return Record{}, err
	}
	return record, writeRecord(outPath, record)
}

func writeRecord(outPath string, record Record) error {
	// writeRecord is intentionally serialization-only; profile decisions and
	// safety checks must already be complete at this point.
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, append(data, '\n'), 0o644)
}

func Load(path string) (Record, error) {
	var record Record
	// Loading a record parses the JSON shape only. Callers still need verifier
	// logic before treating loaded proof as authority.
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
	// A passing witness must be both successful and CI-witnessed; local or
	// agent-reported records cannot satisfy CI-backed trust on status alone.
	return record.Kind == KindGitHubActions &&
		record.Status == StatusPass &&
		record.TrustScope == TrustScopeCIWitnessed
}

func EnvironmentFromOS() map[string]string {
	return environmentFromEntries(os.Environ())
}

func environmentFromEntries(entries []string) map[string]string {
	// The environment snapshot is copied into a map so tests and callers can
	// replay witness behavior without reading process state.
	env := map[string]string{}
	for _, entry := range entries {
		addEnvironmentEntry(env, entry)
	}
	return env
}

func addEnvironmentEntry(env map[string]string, entry string) {
	key, value, ok := strings.Cut(entry, "=")
	if ok {
		// Ignore malformed entries so environment replay cannot create empty keys
		// that would accidentally satisfy required GitHub identity checks.
		env[key] = value
	}
}

func missingGitHubIdentity(env map[string]string) []string {
	// GitHub identity requires both mandatory provider fields and the explicit
	// GITHUB_ACTIONS=true marker.
	missing := missingEnvKeys(env, githubIdentityEnvKeys)
	if env["GITHUB_ACTIONS"] != "" && env["GITHUB_ACTIONS"] != "true" {
		missing = append(missing, "GITHUB_ACTIONS=true")
	}
	sort.Strings(missing)
	return missing
}

var githubIdentityEnvKeys = []string{
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

func missingEnvKeys(env map[string]string, required []string) []string {
	// Whitespace-only environment values are treated as absent; they cannot bind
	// a witness to a CI run.
	missing := make([]string, 0)
	for _, key := range required {
		if strings.TrimSpace(env[key]) == "" {
			missing = append(missing, key)
		}
	}
	return missing
}

// missingGitHubOIDC identifies whether the environment can fetch OIDC
// tokens at all. Absence of these fields means the CI job lacks OIDC
// token provisioning capability, precluding CI-witnessed trust scope.
func missingGitHubOIDC(env map[string]string) []string {
	// Both URL and request token are required; either missing value prevents live
	// OIDC attestation.
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
	// The live fetcher keeps request construction centralized so tests can
	// inject a deterministic TokenFetcher without touching network behavior.
	token, err := fetchGitHubOIDCToken(
		http.DefaultClient,
		env["ACTIONS_ID_TOKEN_REQUEST_URL"],
		env["ACTIONS_ID_TOKEN_REQUEST_TOKEN"],
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

func fetchGitHubOIDCToken(httpClient *http.Client, requestURL, requestToken string) (string, error) {
	// Request creation, transport, and response parsing are split so each trust
	// boundary can fail without exposing the bearer token.
	req, err := buildOIDCTokenRequest(requestURL, requestToken)
	if err != nil {
		return "", err
	}
	body, err := executeOIDCTokenRequest(httpClient, req)
	if err != nil {
		return "", err
	}
	return parseOIDCTokenResponse(body)
}

func buildOIDCTokenRequest(requestURL, requestToken string) (*http.Request, error) {
	// The GitHub-provided request token is only sent to the expected Actions OIDC
	// host after audience normalization.
	// The caller supplies the URL from the CI environment, so host validation
	// happens before headers are attached.
	// Audience is pinned to sdp-trace to prevent replaying a token minted for a
	// different relying party.
	parsed, err := url.Parse(requestURL)
	if err != nil {
		return nil, err
	}
	if err := validateOIDCTokenHost(parsed.Host); err != nil {
		return nil, err
	}
	setOIDCTokenAudience(parsed)

	req, err := http.NewRequest(http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, err
	}
	// Authorization is attached after host and audience validation so the token
	// cannot be sent to an arbitrary environment-provided URL.
	req.Header.Set("Authorization", "bearer "+requestToken)
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func validateOIDCTokenHost(host string) error {
	if !strings.HasSuffix(host, "actions.githubusercontent.com") {
		// Only GitHub's OIDC endpoint is allowed to receive the request token.
		return fmt.Errorf("unexpected oidc request host: %s", host)
	}
	return nil
}

func setOIDCTokenAudience(parsed *url.URL) {
	// Audience pinning prevents a CI job from replaying a token minted for a
	// different relying party into this witness record.
	query := parsed.Query()
	query.Set("audience", "sdp-trace")
	parsed.RawQuery = query.Encode()
}

func executeOIDCTokenRequest(httpClient *http.Client, req *http.Request) ([]byte, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if !successHTTPStatus(resp.StatusCode) {
		// Non-success responses are reported by status only so response bodies
		// cannot leak provider internals into witness errors.
		return nil, fmt.Errorf("oidc token request returned %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func successHTTPStatus(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func parseOIDCTokenResponse(body []byte) (string, error) {
	var payload struct {
		Value string `json:"value"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", err
	}
	if payload.Value == "" {
		// Empty token responses cannot establish identity and are not replaced
		// with environment-only evidence.
		return "", errors.New("oidc token response missing value")
	}
	return payload.Value, nil
}

func parseOIDCClaims(token string) (OIDCClaims, error) {
	// Claim parsing intentionally reads only the JWT payload; signature
	// validation is outside this portable witness helper.
	// The raw token is never retained in the Record; only selected binding claims
	// are returned.
	// Audience normalization happens after JSON decode because GitHub supports
	// both scalar and array forms.
	payload, err := decodeJWTPayload(token)
	if err != nil {
		return OIDCClaims{}, err
	}
	return oidcClaimsFromPayload(payload)
}

func decodeJWTPayload(token string) ([]byte, error) {
	// Only the payload segment is decoded; the raw JWT is not retained in any
	// output artifact.
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid jwt shape")
	}
	return base64.RawURLEncoding.DecodeString(parts[1])
}

func oidcClaimsFromPayload(payload []byte) (OIDCClaims, error) {
	// Decode into a raw struct first so audience normalization can preserve both
	// scalar and array forms without retaining unused token fields.
	var raw rawOIDCClaims
	if err := json.Unmarshal(payload, &raw); err != nil {
		return OIDCClaims{}, err
	}
	return normalizedOIDCClaims(raw), nil
}

func normalizedOIDCClaims(raw rawOIDCClaims) OIDCClaims {
	// Only claims needed for environment binding are retained; unused JWT
	// material is discarded at the trust boundary.
	return OIDCClaims{
		Issuer:     raw.Issuer,
		Subject:    raw.Subject,
		Audience:   audienceString(raw.Audience),
		Repository: raw.Repository,
		Ref:        raw.Ref,
		SHA:        raw.SHA,
	}
}

func audienceString(value any) string {
	// GitHub may encode audience as either a string or an array; store a compact
	// comparable form for trust-context matching.
	switch typed := value.(type) {
	case string:
		return typed
	case []any:
		return strings.Join(stringItems(typed), ",")
	default:
		return ""
	}
}

func stringItems(values []any) []string {
	// Non-string audience entries are ignored instead of coerced, keeping the
	// comparison strict.
	parts := make([]string, 0, len(values))
	for _, item := range values {
		parts = appendStringItem(parts, item)
	}
	return parts
}

func appendStringItem(parts []string, item any) []string {
	// Only literal string audience values participate in the canonical audience
	// list.
	if text, ok := item.(string); ok {
		return append(parts, text)
	}
	return parts
}

// claimsMatchEnvironment is the composite boundary check: it requires
// both trust context (issuer/audience from claimsTrustContextMatches) and
// git context (source binding from claimsGitContextMatches) to be satisfied.
// Failure on either side means the OIDC evidence does not cross the trust boundary.
func claimsMatchEnvironment(claims OIDCClaims, env map[string]string) bool {
	return claimsTrustContextMatches(claims) &&
		claimsGitContextMatches(claims, env)
}

// claimsTrustContextMatches verifies the OIDC issuer and audience are from the
// trusted GitHub OIDC provider. This establishes the trust origin boundary:
// tokens not issued by githubOIDCIssuer or for a different audience do not
// cross this boundary regardless of other claim validity.
func claimsTrustContextMatches(claims OIDCClaims) bool {
	return claims.Issuer == githubOIDCIssuer && claims.Audience == "sdp-trace"
}

// claimsGitContextMatches binds OIDC claims to the live environment values,
// establishing the provenance boundary. A claim is only valid when its
// repository, ref, and SHA match the running environment's git context.
// Mismatch here means the evidence does not correspond to this execution.
func claimsGitContextMatches(claims OIDCClaims, env map[string]string) bool {
	// Source binding requires repository, ref, and SHA to match the live
	// environment together.
	return claims.Repository == env["GITHUB_REPOSITORY"] &&
		claims.Ref == env["GITHUB_REF"] &&
		claims.SHA == env["GITHUB_SHA"]
}

// hashRunArtifacts establishes evidence for run artifacts by computing
// SHA-256 digests of run.json files discovered under runsRoot. These
// digests bind the witness to specific run outputs, enabling later
// detection of artifact tampering or substitution.
func hashRunArtifacts(runsRoot string) ([]ArtifactDigest, error) {
	// Run discovery is delegated to the demo package so witness binding uses the
	// same run-directory rules as gate/demo evaluation.
	// Each retained digest points at run.json, the manifest that carries run
	// identity and chain heads.
	// Digest paths are stored relative to the discovered run directory name.
	runDirs, err := demo.DiscoverRunDirs(runsRoot)
	if err != nil {
		return nil, err
	}
	artifacts := make([]ArtifactDigest, 0, len(runDirs))
	for _, runDir := range runDirs {
		// Each run contributes exactly the retained run manifest digest; command
		// stdout/stderr bodies are not re-read by witness generation.
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

// hashReportArtifacts establishes evidence for report artifacts by
// computing SHA-256 digests of files under reportDir. Unlike run artifacts,
// these are user-supplied outputs; the digest records what was present at
// witness generation time without asserting correctness.
func hashReportArtifacts(reportDir string) ([]ArtifactDigest, error) {
	// Report artifacts are read from one directory level and then sorted by
	// retained path for deterministic witness output.
	entries, err := os.ReadDir(reportDir)
	if err != nil {
		return nil, err
	}
	return hashReportArtifactEntries(reportDir, entries)
}

func hashReportArtifactEntries(reportDir string, entries []os.DirEntry) ([]ArtifactDigest, error) {
	// The generated ci-witness.json is excluded so rerunning witness generation
	// does not self-reference the prior output.
	// Sorting happens after hashing because skipped entries and read errors must
	// be handled before producing the final evidence list.
	// The result is deterministic even when the filesystem returns directory
	// entries in an arbitrary order.
	artifacts := make([]ArtifactDigest, 0)
	for _, entry := range entries {
		if skipReportArtifactEntry(entry) {
			continue
		}
		artifact, err := hashReportArtifact(reportDir, entry.Name())
		if err != nil {
			return nil, err
		}
		artifacts = append(artifacts, artifact)
	}
	sort.Slice(artifacts, func(i, j int) bool {
		return artifacts[i].Path < artifacts[j].Path
	})
	return artifacts, nil
}

func skipReportArtifactEntry(entry os.DirEntry) bool {
	return entry.IsDir() || entry.Name() == "ci-witness.json"
}

func hashReportArtifact(reportDir, name string) (ArtifactDigest, error) {
	// Paths stored in witness output are report-relative and slash-normalized for
	// cross-platform comparison.
	digest, err := hashFile(filepath.Join(reportDir, name))
	if err != nil {
		return ArtifactDigest{}, err
	}
	return ArtifactDigest{Path: filepath.ToSlash(name), SHA256: digest}, nil
}

func hashFile(path string) (string, error) {
	// File hashing returns only SHA-256 hex; callers decide how the digest is
	// bound to a run or report artifact.
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("hash %s: %w", path, err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
