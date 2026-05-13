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

func TestPolicyValidateRequiredFields(t *testing.T) {
	valid := AuthorityPolicy{
		SchemaVersion:          "v1",
		PolicyID:               "policy-a",
		AllowedSigners:         []SignerAuthorityEntry{{SignerID: "signer-a"}},
		AllowedWitnessProfiles: []string{"github-actions"},
	}
	tests := []struct {
		name    string
		policy  AuthorityPolicy
		wantErr string
	}{
		{name: "valid", policy: valid},
		{
			name:    "schema",
			policy:  AuthorityPolicy{PolicyID: valid.PolicyID, AllowedSigners: valid.AllowedSigners, AllowedWitnessProfiles: valid.AllowedWitnessProfiles},
			wantErr: "schema_version is required",
		},
		{
			name:    "policy",
			policy:  AuthorityPolicy{SchemaVersion: valid.SchemaVersion, AllowedSigners: valid.AllowedSigners, AllowedWitnessProfiles: valid.AllowedWitnessProfiles},
			wantErr: "policy_id is required",
		},
		{
			name:    "signers",
			policy:  AuthorityPolicy{SchemaVersion: valid.SchemaVersion, PolicyID: valid.PolicyID, AllowedWitnessProfiles: valid.AllowedWitnessProfiles},
			wantErr: "at least one allowed signer is required",
		},
		{
			name:    "witness_profiles",
			policy:  AuthorityPolicy{SchemaVersion: valid.SchemaVersion, PolicyID: valid.PolicyID, AllowedSigners: valid.AllowedSigners},
			wantErr: "at least one allowed witness profile is required",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.policy.Validate()
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v", err)
				}
				return
			}
			if err == nil || err.Error() != test.wantErr {
				t.Fatalf("Validate() error = %v, want %q", err, test.wantErr)
			}
		})
	}
}
