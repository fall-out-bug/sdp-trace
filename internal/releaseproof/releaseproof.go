package releaseproof

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	SchemaVersion = "0.1.0"
	TrustScope    = "source_bound_local_release"

	StatePass         = "pass"
	StateFail         = "fail"
	StateCannotVerify = "cannot_verify"
	StatusMatched     = "matched"
	StatusMismatch    = "mismatch"
	StatusMissing     = "missing"
	StatusNotAssessed = "not_assessed"
)

type Manifest struct {
	ID                       string             `json:"id"`
	SigningProfile           string             `json:"signing_profile"`
	TrustedIdentityPolicyRef string             `json:"trusted_identity_policy_ref"`
	SourceCommit             string             `json:"source_commit"`
	Artifacts                []ManifestArtifact `json:"artifacts"`
	Accountability           Accountability     `json:"accountability"`
}

type ManifestArtifact struct {
	Path   string `json:"path"`
	Kind   string `json:"kind"`
	SHA256 string `json:"sha256"`
}

type Accountability struct {
	DRI                 Actor  `json:"dri"`
	Approver            Actor  `json:"approver"`
	Escalation          Actor  `json:"escalation"`
	AuthorityScope      string `json:"authority_scope"`
	AccountabilityClaim string `json:"accountability_claim"`
	ApprovalRef         string `json:"approval_ref"`
	RiskOwner           Actor  `json:"risk_owner"`
	LineOfDefense       string `json:"line_of_defense"`
}

type Actor struct {
	IdentityRef string `json:"identity_ref"`
	ActorType   string `json:"actor_type"`
}

type Verification struct {
	ID                         string            `json:"id"`
	SchemaVersion              string            `json:"schema_version"`
	ArtifactRole               string            `json:"artifact_role"`
	TrustScope                 string            `json:"trust_scope"`
	ReleaseVerificationState   string            `json:"release_verification_state"`
	ManifestRef                string            `json:"manifest_ref"`
	ManifestDigest             string            `json:"manifest_digest"`
	ManifestDigestStatus       string            `json:"manifest_digest_status"`
	ArtifactDigestStatus       string            `json:"artifact_digest_status"`
	SignatureProfile           string            `json:"signature_profile"`
	SignatureStatus            string            `json:"signature_status"`
	IdentityPolicyRef          string            `json:"identity_policy_ref"`
	IdentityPolicyStatus       string            `json:"identity_policy_status"`
	SourceCommit               string            `json:"source_commit"`
	SourceCommitStatus         string            `json:"source_commit_status"`
	SourceCommitArtifactStatus string            `json:"source_commit_artifact_status"`
	SourceCommitArtifactCounts ArtifactCounts    `json:"source_commit_artifact_counts"`
	ExternalTrustProfile       string            `json:"external_trust_profile"`
	ExternalAttestationRef     *string           `json:"external_attestation_ref"`
	TransparencyEvidenceRef    *string           `json:"transparency_evidence_ref"`
	SourceURIStatus            string            `json:"source_uri_status"`
	ProtectedRefStatus         string            `json:"protected_ref_status"`
	WorkflowIdentityStatus     string            `json:"workflow_identity_status"`
	ApprovalStatus             string            `json:"approval_status"`
	ProductionReleaseVerified  ProofStateBoolean `json:"production_release_verified"`
	TransparencyLogStatus      string            `json:"transparency_log_status"`
	FreshnessStatus            string            `json:"freshness_status"`
	VerifiedAt                 string            `json:"verified_at"`
	TrustedContractRelease     bool              `json:"trusted_contract_release"`
	PrivateEquivalentProfile   string            `json:"private_equivalent_profile,omitempty"`
	ProvenanceRefs             []string          `json:"provenance_refs,omitempty"`
	Accountability             Accountability    `json:"accountability"`
	SourceCommitReason         string            `json:"source_commit_reason,omitempty"`
	ExternalTrustReason        string            `json:"external_trust_reason,omitempty"`
	ArtifactIssues             []ArtifactIssue   `json:"artifact_issues,omitempty"`
}

type ArtifactCounts struct {
	Checked    int `json:"checked"`
	Missing    int `json:"missing"`
	Mismatched int `json:"mismatched"`
}

type ProofStateBoolean struct {
	State  string `json:"state"`
	Value  *bool  `json:"value"`
	Reason string `json:"reason,omitempty"`
}

type ArtifactIssue struct {
	Path     string `json:"path"`
	Issue    string `json:"issue"`
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
}

type manifestData struct {
	manifest Manifest
	ref      string
	digest   [sha256.Size]byte
}

type verificationState struct {
	sourceCommit   string
	state          string
	commitStatus   string
	artifactStatus string
	artifactCounts ArtifactCounts
	artifactIssues []ArtifactIssue
	sourceReason   string
}

type verificationInput struct {
	manifestData
	verificationState
	verificationTime time.Time
}

func Evaluate(repoRoot, manifestPath string, now time.Time) (Verification, error) {
	// Evaluation has three trust phases: load bounded manifest evidence,
	// compare it to the source commit, then render a conservative verdict.
	manifestData, err := loadManifest(repoRoot, manifestPath)
	if err != nil {
		return Verification{}, err
	}
	state := evaluateManifestState(repoRoot, manifestData.manifest)
	return buildVerification(verificationInput{
		manifestData:      manifestData,
		verificationState: state,
		verificationTime:  now,
	}), nil
}

func loadManifest(repoRoot, manifestPath string) (manifestData, error) {
	// Resolve the manifest through the repository boundary before hashing so
	// the verification record names the same relative file that was read.
	manifestRel, manifestBytes, err := loadManifestBytes(repoRoot, manifestPath)
	if err != nil {
		return manifestData{}, err
	}
	var manifest Manifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return manifestData{}, err
	}
	// The digest is computed over the exact bytes accepted for parsing.
	return manifestData{
		manifest: manifest,
		ref:      manifestRel,
		digest:   sha256.Sum256(manifestBytes),
	}, nil
}

func evaluateManifestState(repoRoot string, manifest Manifest) verificationState {
	// Source-commit status is the root local proof; artifact checks depend on
	// that immutable source boundary being available.
	manifestSourceCommit := strings.TrimSpace(manifest.SourceCommit)
	commitStatus, reason := sourceCommitState(repoRoot, manifestSourceCommit)
	state := initialReleaseState(commitStatus)
	// Artifact results and dirty-check state can only lower confidence.
	counts, issues, artifactStatus, artifactReason := artifactVerificationState(repoRoot, manifestSourceCommit, manifest.Artifacts, state)
	state, reason = combineState(state, reason, artifactStatus, artifactReason)
	// Dirty checkout evidence is local structural evidence, not source proof.
	state, commitStatus, reason = applyDirtyState(repoRoot, state, commitStatus, reason)
	// Keep the rendered verification fields separate from decision logic.
	return verificationState{
		sourceCommit:   manifestSourceCommit,
		state:          state,
		commitStatus:   commitStatus,
		artifactStatus: artifactStatus,
		artifactCounts: counts,
		artifactIssues: issues,
		sourceReason:   reason,
	}
}

func initialReleaseState(commitStatus string) string {
	// A missing source commit breaks the immutable source boundary before any
	// artifact digest can be trusted as release proof.
	if commitStatus == StatusMissing {
		return StateCannotVerify
	}
	return StatePass
}

func buildVerification(input verificationInput) Verification {
	// Build the public proof in groups matching the trust boundaries above.
	result := baseVerification(input)
	applyManifestEvidence(&result, input)
	applySourceEvidence(&result, input.verificationState)
	applyExternalTrustDefaults(&result)
	applyReleaseTrustDefaults(&result)
	return result
}

func baseVerification(input verificationInput) Verification {
	// These fields identify the verifier output before any evidence-specific
	// status is attached.
	return Verification{
		ID:                       "contract-release-verification-block-18-19-source-bound",
		SchemaVersion:            SchemaVersion,
		ArtifactRole:             "verifier_output",
		TrustScope:               TrustScope,
		ReleaseVerificationState: input.state,
		VerifiedAt:               input.verificationTime.UTC().Format(time.RFC3339),
	}
}

func applyManifestEvidence(result *Verification, input verificationInput) {
	// Manifest evidence is local and byte-bound: reference, digest, signature
	// profile, identity policy, and accountability all come from one payload.
	result.ManifestRef = input.ref
	result.ManifestDigest = hex.EncodeToString(input.digest[:])
	result.ManifestDigestStatus = StatusMatched
	result.SignatureProfile = input.manifest.SigningProfile
	result.SignatureStatus = StatusNotAssessed
	result.IdentityPolicyRef = input.manifest.TrustedIdentityPolicyRef
	result.IdentityPolicyStatus = StatusNotAssessed
	result.Accountability = input.manifest.Accountability
}

func applySourceEvidence(result *Verification, state verificationState) {
	// Source evidence records only what was checked against the manifest's
	// source commit; it does not make external production claims.
	result.ArtifactDigestStatus = state.artifactStatus
	result.SourceCommit = state.sourceCommit
	result.SourceCommitStatus = state.commitStatus
	result.SourceCommitArtifactStatus = state.artifactStatus
	result.SourceCommitArtifactCounts = state.artifactCounts
	result.SourceCommitReason = state.sourceReason
	result.ArtifactIssues = state.artifactIssues
}

func applyExternalTrustDefaults(result *Verification) {
	// External production controls are intentionally outside this local source
	// profile, so each related field stays explicitly not_assessed.
	result.ExternalTrustProfile = StatusNotAssessed
	result.ExternalAttestationRef = nil
	result.TransparencyEvidenceRef = nil
	result.SourceURIStatus = StatusNotAssessed
	result.ProtectedRefStatus = StatusNotAssessed
	result.WorkflowIdentityStatus = StatusNotAssessed
	result.ApprovalStatus = StatusNotAssessed
}

func applyReleaseTrustDefaults(result *Verification) {
	// Local source proof is not equivalent to release authorization.
	result.ProductionReleaseVerified = productionReleaseNotAssessed()
	result.TransparencyLogStatus = StatusNotAssessed
	result.FreshnessStatus = StatusNotAssessed
	result.TrustedContractRelease = false
	result.PrivateEquivalentProfile = "not_assessed"
	result.ExternalTrustReason = "external production trust is not assessed by the local source-bound profile"
}

func productionReleaseNotAssessed() ProofStateBoolean {
	// The boolean value is absent because this profile does not inspect
	// production attestation evidence.
	return ProofStateBoolean{
		State:  StatusNotAssessed,
		Value:  nil,
		Reason: "Production release verification requires external attestation in addition to source-bound local checks.",
	}
}

func loadManifestBytes(repoRoot, manifestPath string) (string, []byte, error) {
	// Normalize and resolve before reading so ManifestRef and ManifestDigest
	// describe the same repository-contained file.
	manifestRel, err := cleanRepoRelativePath(manifestPath)
	if err != nil {
		return "", nil, err
	}
	manifestAbs, err := resolveRepoFile(repoRoot, manifestRel)
	if err != nil {
		return "", nil, err
	}
	manifestBytes, err := os.ReadFile(manifestAbs)
	return manifestRel, manifestBytes, err
}

func resolveRepoFile(repoRoot, relPath string) (string, error) {
	// Resolve symlinks before the containment check so a repository-relative
	// path cannot point at a file outside the repository.
	root, target, err := resolvedRepoAndTarget(repoRoot, relPath)
	if err != nil {
		return "", err
	}
	relative, err := filepath.Rel(root, target)
	if err != nil {
		return "", err
	}
	if repoRelativePathInside(relative) {
		return target, nil
	}
	return "", fmt.Errorf("manifest path %q resolves outside repository", relPath)
}

func repoRelativePathInside(relative string) bool {
	// filepath.Rel may return "." for the repository root; everything else
	// must stay below the root without a leading parent traversal.
	return relative == "." || (!strings.HasPrefix(relative, ".."+string(filepath.Separator)) && relative != "..")
}

func resolvedRepoAndTarget(repoRoot, relPath string) (string, string, error) {
	// Compare canonical paths so symlinked roots and symlinked manifests use
	// the same filesystem view during containment checks.
	root, err := filepath.EvalSymlinks(repoRoot)
	if err != nil {
		return "", "", err
	}
	target, err := filepath.EvalSymlinks(filepath.Join(repoRoot, relPath))
	if err != nil {
		return "", "", err
	}
	return root, target, nil
}

func cleanRepoRelativePath(path string) (string, error) {
	// The manifest reference is stored as portable slash-separated repository
	// relative data; absolute or parent paths are never accepted as evidence.
	clean := filepath.Clean(strings.TrimSpace(path))
	if clean == "." || clean == "" {
		return "", errors.New("manifest path is required")
	}
	if unsafeRepoRelativePath(clean) {
		return "", fmt.Errorf("manifest path must be repository-relative: %s", path)
	}
	return filepath.ToSlash(clean), nil
}

func unsafeRepoRelativePath(clean string) bool {
	return filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator))
}

func sourceCommitState(repoRoot, sourceCommit string) (string, string) {
	// The source commit is the anchor for every source-bound artifact check; if
	// it cannot resolve locally, artifact claims stay unverified.
	if sourceCommit == "" {
		return StatusMissing, "manifest source_commit is missing"
	}
	if !sourceCommitExists(repoRoot, sourceCommit) {
		return StatusMissing, "manifest source_commit could not be resolved from git"
	}
	return StatusMatched, "source commit contains every manifest artifact path with matching digest"
}

func artifactVerificationState(repoRoot, sourceCommit string, artifacts []ManifestArtifact, state string) (ArtifactCounts, []ArtifactIssue, string, string) {
	// Do not inspect artifacts when the source anchor is absent; that would turn
	// missing source proof into misleading artifact evidence.
	if state == StateCannotVerify {
		return ArtifactCounts{}, nil, StatusNotAssessed, "manifest artifacts were not checked because source_commit cannot be verified"
	}
	counts, issues := artifactCounts(repoRoot, sourceCommit, artifacts)
	artifactStatus, artifactReason := artifactState(counts)
	return counts, issues, artifactStatus, artifactReason
}

func artifactState(counts ArtifactCounts) (string, string) {
	// Missing artifact paths are reported before digest mismatches because the
	// verifier cannot compare bytes that were not present in the source commit.
	if counts.Missing > 0 {
		return StatusMissing, "manifest artifact paths are missing from the current source checkout"
	}
	if counts.Mismatched > 0 {
		return StatusMismatch, "manifest artifact digests do not match the current source checkout"
	}
	return StatusMatched, ""
}

func combineState(state, reason, artifactStatus, artifactReason string) (string, string) {
	// Artifact verification can only lower confidence; it never upgrades a
	// missing source commit or previously unverified release proof.
	if state == StateCannotVerify || artifactStatus == StatusMatched {
		return state, reason
	}
	return StateFail, artifactReason
}

func applyDirtyState(repoRoot, state, commitStatus, reason string) (string, string, string) {
	// A dirty checkout is local structural evidence only, so it blocks a source
	// match without turning external trust green.
	if state == StateCannotVerify || !worktreeDirty(repoRoot) {
		return state, commitStatus, reason
	}
	return StateFail, StatusMismatch, "dirty checkout cannot support source-bound local release proof"
}

func artifactCounts(repoRoot, sourceCommit string, artifacts []ManifestArtifact) (ArtifactCounts, []ArtifactIssue) {
	// Count every manifest artifact against the immutable source commit; each
	// issue keeps the manifest path so reports stay source-bound and auditable.
	// Checked counts manifest obligations, not successful reads, so missing
	// artifacts are visible in both the denominator and issue list.
	counts := ArtifactCounts{Checked: len(artifacts)}
	// Issues are sparse: matched artifacts stay represented by Checked, while
	// only missing or mismatched obligations get explicit rows.
	issues := []ArtifactIssue{}
	for _, artifact := range artifacts {
		issue, ok := artifactIssue(repoRoot, sourceCommit, artifact)
		if !ok {
			continue
		}
		countArtifactIssue(&counts, issue)
		issues = append(issues, issue)
	}
	return counts, issues
}

func countArtifactIssue(counts *ArtifactCounts, issue ArtifactIssue) {
	// Only replay failures change the aggregate counters; matched artifacts are
	// already represented by Checked.
	switch issue.Issue {
	case StatusMissing:
		counts.Missing++
	case StatusMismatch:
		counts.Mismatched++
	}
}

func artifactIssue(repoRoot, sourceCommit string, artifact ManifestArtifact) (ArtifactIssue, bool) {
	// Clean the path used for git object lookup, but keep the original
	// manifest path in issues so reports match the signed obligation text.
	path := filepath.Clean(artifact.Path)
	data, err := artifactBytes(repoRoot, sourceCommit, path)
	if err != nil {
		// Missing source-commit bytes are stronger than digest mismatch because
		// there is no artifact content to compare.
		return ArtifactIssue{Path: artifact.Path, Issue: StatusMissing, Expected: artifact.SHA256}, true
	}
	sum := sha256.Sum256(data)
	actual := hex.EncodeToString(sum[:])
	// Hex case is formatting, not proof content; compare digest values
	// case-insensitively while reporting the canonical lowercase actual.
	if strings.EqualFold(actual, artifact.SHA256) {
		return ArtifactIssue{}, false
	}
	return ArtifactIssue{Path: artifact.Path, Issue: StatusMismatch, Expected: artifact.SHA256, Actual: actual}, true
}

func sourceCommitExists(repoRoot, sourceCommit string) bool {
	// Git object resolution is the immutable-source boundary for this local
	// verifier; a missing object keeps the release verdict at cannot_verify.
	cmd := exec.Command("git", "cat-file", "-e", sourceCommit+"^{commit}")
	cmd.Dir = repoRoot
	return cmd.Run() == nil
}

func artifactBytes(repoRoot, sourceCommit, path string) ([]byte, error) {
	// Read artifacts from the manifest source commit, not the dirty checkout,
	// so local edits cannot satisfy source-bound release proof.
	cmd := exec.Command("git", "show", sourceCommit+":"+path)
	cmd.Dir = repoRoot
	return cmd.Output()
}

func worktreeDirty(repoRoot string) bool {
	// Treat git status failures as dirty so command failures cannot accidentally
	// promote a source-bound proof.
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		return true
	}
	return strings.TrimSpace(string(out)) != ""
}

func Write(path string, result Verification) error {
	// Persist verifier output as stable pretty JSON because downstream evidence
	// checks compare the proof artifact as a reviewable file.
	payload, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	payload = append(payload, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, payload, 0o644)
}

func Read(path string) (Verification, error) {
	// Include the proof path in parse errors so malformed evidence can be traced
	// back to the exact artifact under review.
	data, err := os.ReadFile(path)
	if err != nil {
		return Verification{}, err
	}
	var result Verification
	if err := json.Unmarshal(data, &result); err != nil {
		return Verification{}, fmt.Errorf("release proof %s: %w", path, err)
	}
	return result, nil
}

func RepoRoot(dir string) (string, error) {
	// Git decides the repository root for source-bound proof; callers should not
	// infer it from process working-directory assumptions.
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("repo root cannot be determined from %s: %w", dir, err)
	}
	root := strings.TrimSpace(string(out))
	if root == "" {
		return "", fmt.Errorf("repo root cannot be determined from %s", dir)
	}
	return root, nil
}
