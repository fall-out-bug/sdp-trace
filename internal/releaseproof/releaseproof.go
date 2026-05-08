package releaseproof

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

func Evaluate(repoRoot, manifestPath string, now time.Time) (Verification, error) {
	manifestBytes, err := os.ReadFile(filepath.Join(repoRoot, manifestPath))
	if err != nil {
		return Verification{}, err
	}
	var manifest Manifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return Verification{}, err
	}
	manifestSourceCommit := strings.TrimSpace(manifest.SourceCommit)
	commitStatus := StatusMatched
	dirty := worktreeDirty(repoRoot)
	counts := ArtifactCounts{Checked: len(manifest.Artifacts)}
	var issues []ArtifactIssue
	artifactStatus := StatusMatched
	state := StatePass
	reason := "source commit contains every manifest artifact path with matching digest"
	if manifestSourceCommit == "" {
		commitStatus = StatusMissing
		state = StateCannotVerify
		reason = "manifest source_commit is missing"
	} else if !sourceCommitExists(repoRoot, manifestSourceCommit) {
		commitStatus = StatusMissing
		state = StateCannotVerify
		reason = "manifest source_commit could not be resolved from git"
	} else {
		counts, issues = artifactCounts(repoRoot, manifestSourceCommit, manifest.Artifacts)
	}
	if counts.Missing > 0 {
		artifactStatus = StatusMissing
		if state != StateCannotVerify {
			state = StateFail
			reason = "manifest artifact paths are missing from the current source checkout"
		}
	} else if counts.Mismatched > 0 {
		artifactStatus = StatusMismatch
		if state != StateCannotVerify {
			state = StateFail
			reason = "manifest artifact digests do not match the current source checkout"
		}
	}
	if dirty && state != StateCannotVerify {
		commitStatus = StatusMismatch
		state = StateFail
		reason = "dirty checkout cannot support source-bound local release proof"
	}
	manifestDigest := sha256.Sum256(manifestBytes)
	return Verification{
		ID:                         "contract-release-verification-block-18-19-source-bound",
		SchemaVersion:              SchemaVersion,
		ArtifactRole:               "verifier_output",
		TrustScope:                 TrustScope,
		ReleaseVerificationState:   state,
		ManifestRef:                manifestPath,
		ManifestDigest:             hex.EncodeToString(manifestDigest[:]),
		ManifestDigestStatus:       StatusMatched,
		ArtifactDigestStatus:       artifactStatus,
		SignatureProfile:           manifest.SigningProfile,
		SignatureStatus:            StatusNotAssessed,
		IdentityPolicyRef:          manifest.TrustedIdentityPolicyRef,
		IdentityPolicyStatus:       StatusNotAssessed,
		SourceCommit:               manifestSourceCommit,
		SourceCommitStatus:         commitStatus,
		SourceCommitArtifactStatus: artifactStatus,
		SourceCommitArtifactCounts: counts,
		ExternalTrustProfile:       StatusNotAssessed,
		ExternalAttestationRef:     nil,
		TransparencyEvidenceRef:    nil,
		SourceURIStatus:            StatusNotAssessed,
		ProtectedRefStatus:         StatusNotAssessed,
		WorkflowIdentityStatus:     StatusNotAssessed,
		ApprovalStatus:             StatusNotAssessed,
		ProductionReleaseVerified: ProofStateBoolean{
			State:  StatusNotAssessed,
			Value:  nil,
			Reason: "Production release verification requires external attestation in addition to source-bound local checks.",
		},
		TransparencyLogStatus:    StatusNotAssessed,
		FreshnessStatus:          StatusNotAssessed,
		VerifiedAt:               now.UTC().Format(time.RFC3339),
		TrustedContractRelease:   false,
		PrivateEquivalentProfile: "not_assessed",
		Accountability:           manifest.Accountability,
		SourceCommitReason:       reason,
		ExternalTrustReason:      "external production trust is not assessed by the local source-bound profile",
		ArtifactIssues:           issues,
	}, nil
}

func artifactCounts(repoRoot, sourceCommit string, artifacts []ManifestArtifact) (ArtifactCounts, []ArtifactIssue) {
	counts := ArtifactCounts{Checked: len(artifacts)}
	issues := []ArtifactIssue{}
	for _, artifact := range artifacts {
		path := filepath.Clean(artifact.Path)
		data, err := artifactBytes(repoRoot, sourceCommit, path)
		if err != nil {
			counts.Missing++
			issues = append(issues, ArtifactIssue{Path: artifact.Path, Issue: StatusMissing, Expected: artifact.SHA256})
			continue
		}
		sum := sha256.Sum256(data)
		actual := hex.EncodeToString(sum[:])
		if !strings.EqualFold(actual, artifact.SHA256) {
			counts.Mismatched++
			issues = append(issues, ArtifactIssue{Path: artifact.Path, Issue: StatusMismatch, Expected: artifact.SHA256, Actual: actual})
		}
	}
	return counts, issues
}

func sourceCommitExists(repoRoot, sourceCommit string) bool {
	cmd := exec.Command("git", "cat-file", "-e", sourceCommit+"^{commit}")
	cmd.Dir = repoRoot
	return cmd.Run() == nil
}

func artifactBytes(repoRoot, sourceCommit, path string) ([]byte, error) {
	cmd := exec.Command("git", "show", sourceCommit+":"+path)
	cmd.Dir = repoRoot
	return cmd.Output()
}

func worktreeDirty(repoRoot string) bool {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		return true
	}
	return strings.TrimSpace(string(out)) != ""
}

func Write(path string, result Verification) error {
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
