package releaseproof

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestEvaluateFailsWhenManifestArtifactsAreMissing(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	runGit(t, root, "config", "user.email", "test@example.invalid")
	runGit(t, root, "config", "user.name", "Test")
	writeFile(t, filepath.Join(root, "present.txt"), "present\n")
	runGit(t, root, "add", "present.txt")
	runGit(t, root, "commit", "-m", "source")
	head := strings.TrimSpace(runGit(t, root, "rev-parse", "HEAD"))
	manifest := `{
  "id": "test-manifest",
  "signing_profile": "sdp-trace-signature/sigstore-dsse-keyless-v1",
  "trusted_identity_policy_ref": "policy.json",
  "source_commit": "` + head + `",
  "artifacts": [
    {"path": "present.txt", "kind": "doc", "sha256": "0000000000000000000000000000000000000000000000000000000000000000"},
    {"path": "missing.txt", "kind": "doc", "sha256": "1111111111111111111111111111111111111111111111111111111111111111"}
  ],
  "accountability": {
    "dri": {"identity_ref": "role:dri", "actor_type": "human_role"},
    "approver": {"identity_ref": "role:approver", "actor_type": "human_role"},
    "escalation": {"identity_ref": "role:cto", "actor_type": "human_role"},
    "authority_scope": "contract_release",
    "accountability_claim": "release_approval",
    "approval_ref": "approval",
    "risk_owner": {"identity_ref": "role:risk", "actor_type": "human_role"},
    "line_of_defense": "second"
  }
	}`
	writeFile(t, filepath.Join(root, "manifest.json"), manifest)

	result, err := Evaluate(root, "manifest.json", time.Date(2026, 5, 7, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if result.ReleaseVerificationState != StateFail {
		t.Fatalf("state = %s", result.ReleaseVerificationState)
	}
	if result.SourceCommitArtifactStatus != StatusMissing {
		t.Fatalf("artifact status = %s", result.SourceCommitArtifactStatus)
	}
	if result.SourceCommitArtifactCounts.Checked != 2 || result.SourceCommitArtifactCounts.Missing != 1 || result.SourceCommitArtifactCounts.Mismatched != 1 {
		t.Fatalf("counts = %+v", result.SourceCommitArtifactCounts)
	}
	if result.TrustedContractRelease {
		t.Fatalf("local source-bound failure must not be trusted")
	}
}

func TestEvaluateCannotVerifyWhenManifestSourceCommitMissing(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "present.txt"), "present\n")
	manifest := `{
  "id": "test-manifest",
  "signing_profile": "sdp-trace-signature/sigstore-dsse-keyless-v1",
  "trusted_identity_policy_ref": "policy.json",
  "artifacts": [
    {"path": "present.txt", "kind": "doc", "sha256": "0000000000000000000000000000000000000000000000000000000000000000"}
  ],
  "accountability": {
    "dri": {"identity_ref": "role:dri", "actor_type": "human_role"},
    "approver": {"identity_ref": "role:approver", "actor_type": "human_role"},
    "escalation": {"identity_ref": "role:cto", "actor_type": "human_role"},
    "authority_scope": "contract_release",
    "accountability_claim": "release_approval",
    "approval_ref": "approval",
    "risk_owner": {"identity_ref": "role:risk", "actor_type": "human_role"},
    "line_of_defense": "second"
  }
	}`
	writeFile(t, filepath.Join(root, "manifest.json"), manifest)

	result, err := Evaluate(root, "manifest.json", time.Date(2026, 5, 7, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if result.ReleaseVerificationState != StateCannotVerify {
		t.Fatalf("state = %s", result.ReleaseVerificationState)
	}
	if result.SourceCommitStatus != StatusMissing {
		t.Fatalf("source commit status = %s", result.SourceCommitStatus)
	}
}

func TestEvaluateReadsArtifactsFromManifestSourceCommit(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	runGit(t, root, "config", "user.email", "test@example.invalid")
	runGit(t, root, "config", "user.name", "Test")
	writeFile(t, filepath.Join(root, "present.txt"), "present\n")
	runGit(t, root, "add", "present.txt")
	runGit(t, root, "commit", "-m", "source")
	sourceCommit := strings.TrimSpace(runGit(t, root, "rev-parse", "HEAD"))
	if len(sourceCommit) != 40 {
		t.Fatalf("source commit = %q", sourceCommit)
	}
	writeFile(t, filepath.Join(root, "present.txt"), "changed after source\n")
	runGit(t, root, "add", "present.txt")
	runGit(t, root, "commit", "-m", "proof-artifact-commit")
	manifest := `{
  "id": "test-manifest",
  "signing_profile": "sdp-trace-signature/sigstore-dsse-keyless-v1",
  "trusted_identity_policy_ref": "policy.json",
  "source_commit": "` + sourceCommit + `",
  "artifacts": [
    {"path": "present.txt", "kind": "doc", "sha256": "` + sha256String("present\n") + `"}
  ],
  "accountability": {
    "dri": {"identity_ref": "role:dri", "actor_type": "human_role"},
    "approver": {"identity_ref": "role:approver", "actor_type": "human_role"},
    "escalation": {"identity_ref": "role:cto", "actor_type": "human_role"},
    "authority_scope": "contract_release",
    "accountability_claim": "release_approval",
    "approval_ref": "approval",
    "risk_owner": {"identity_ref": "role:risk", "actor_type": "human_role"},
    "line_of_defense": "second"
  }
	}`
	writeFile(t, filepath.Join(root, "manifest.json"), manifest)
	runGit(t, root, "add", "manifest.json")
	runGit(t, root, "commit", "-m", "manifest")

	result, err := Evaluate(root, "manifest.json", time.Date(2026, 5, 7, 0, 0, 0, 0, time.UTC))
	if err != nil {
		t.Fatalf("evaluate: %v", err)
	}
	if result.ReleaseVerificationState != StatePass {
		t.Fatalf("state = %s reason=%s status=%s issues=%+v", result.ReleaseVerificationState, result.SourceCommitReason, result.SourceCommitArtifactStatus, result.ArtifactIssues)
	}
	if result.SourceCommit != sourceCommit {
		t.Fatalf("source commit = %s", result.SourceCommit)
	}
	if result.SourceCommitStatus != StatusMatched {
		t.Fatalf("source commit status = %s", result.SourceCommitStatus)
	}
	if result.SourceCommitArtifactStatus != StatusMatched {
		t.Fatalf("source artifact status = %s", result.SourceCommitArtifactStatus)
	}
}

func TestWriteReadAndRepoRoot(t *testing.T) {
	root := t.TempDir()
	runGit(t, root, "init")
	result := Verification{
		ID:                       "verification",
		SchemaVersion:            SchemaVersion,
		TrustScope:               TrustScope,
		ReleaseVerificationState: StatePass,
		ManifestRef:              "manifest.json",
		SourceCommit:             "1111111111111111111111111111111111111111",
		Accountability: Accountability{
			DRI:            Actor{IdentityRef: "role:dri", ActorType: "human_role"},
			Approver:       Actor{IdentityRef: "role:approver", ActorType: "human_role"},
			Escalation:     Actor{IdentityRef: "role:cto", ActorType: "human_role"},
			RiskOwner:      Actor{IdentityRef: "role:risk", ActorType: "human_role"},
			AuthorityScope: "contract_release",
			LineOfDefense:  "second",
		},
	}
	path := filepath.Join(root, "nested", "release-proof.json")
	if err := Write(path, result); err != nil {
		t.Fatalf("write: %v", err)
	}
	loaded, err := Read(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if loaded.ID != result.ID || loaded.SourceCommit != result.SourceCommit {
		t.Fatalf("loaded = %+v", loaded)
	}
	nested := filepath.Join(root, "nested")
	repoRoot, err := RepoRoot(nested)
	if err != nil {
		t.Fatalf("repo root: %v", err)
	}
	canonicalRepoRoot, err := filepath.EvalSymlinks(repoRoot)
	if err != nil {
		t.Fatalf("canonical repo root: %v", err)
	}
	canonicalRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("canonical root: %v", err)
	}
	if canonicalRepoRoot != canonicalRoot {
		t.Fatalf("repo root = %s want %s", repoRoot, root)
	}
}

func writeFile(t *testing.T, path string, value string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(value), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func sha256String(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
