package policy

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func TestLoadPolicyAndValidator(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policy.json")
	payload := []byte(`{
  "schema_version": "v1",
  "policy_id": "policy-a",
  "allowed_adapters": [
    {
      "adapter_id": "adapter-a",
      "provider": "generic",
      "identity_state": "verified",
      "allowed_event_types": ["run_started"],
      "allowed_by_policy": true
    }
  ],
  "allowed_signers": [
    {
      "signer_id": "signer-a",
      "profile_id": "github-actions",
      "allowed_scopes": ["signer"],
      "independence_required": true,
      "environment": "ci"
    }
  ],
  "allowed_witness_profiles": ["github-actions"],
  "demonstration_profiles": ["local-dev"]
}`)
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatal(err)
	}

	loaded, err := LoadPolicy(path)
	if err != nil {
		t.Fatalf("load policy: %v", err)
	}
	validator := NewAuthorityPolicyValidator(loaded)

	if !validator.CanAdapterEmit("adapter-a", trace.EventType("run_started")) {
		t.Fatalf("expected adapter event to be allowed")
	}
	if validator.CanAdapterEmit("adapter-a", trace.EventType("run_closed")) {
		t.Fatalf("unexpected adapter event allowed")
	}
	if !validator.SignerAllowed("signer-a", "github-actions") {
		t.Fatalf("expected signer to be allowed")
	}
	if !validator.WitnessProfileAllowed("github-actions") {
		t.Fatalf("expected witness profile to be allowed")
	}
	if !validator.IsDemonstrationProfile("local-dev") {
		t.Fatalf("expected demonstration profile")
	}
}

func TestPolicyValidateRequiresWitnessProfiles(t *testing.T) {
	policy := AuthorityPolicy{
		SchemaVersion:  "v1",
		PolicyID:       "policy-a",
		AllowedSigners: []SignerAuthorityEntry{{SignerID: "signer-a"}},
	}
	if err := policy.Validate(); err == nil || err.Error() != "at least one allowed witness profile is required" {
		t.Fatalf("expected witness profile validation error, got %v", err)
	}
}
