package releaseproof

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestEvaluateFailsWhenManifestArtifactsAreMissing(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "present.txt"), "present\n")
	manifest := `{
  "id": "test-manifest",
  "signing_profile": "sdp-trace-signature/sigstore-dsse-keyless-v1",
  "trusted_identity_policy_ref": "policy.json",
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

func writeFile(t *testing.T, path string, value string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(value), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}
