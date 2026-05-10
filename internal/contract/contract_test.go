package contract

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadValidateAndDigestCanonicalizesSliceOrder(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "contract.json")
	payload := []byte(`{
  "schema_version": "v1",
  "contract_id": "contract-a",
  "version": "2026-05-10",
  "contract_source": "spec",
  "lock_required_before": "implementation",
  "required_observers": ["zeta", "alpha"],
  "optional_observers": ["beta"],
  "required_events": ["run_closed", "run_started"],
  "gate_events": ["gate_checked"],
  "minimum_gate_trust_scope": "local_observed",
  "retention_profile": "digest_only",
  "redaction_profile": "safe"
}`)
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	digest, err := loaded.Digest()
	if err != nil {
		t.Fatalf("digest: %v", err)
	}

	reordered := loaded
	reordered.RequiredObservers = []string{"alpha", "zeta"}
	reordered.RequiredEvents = []string{"run_started", "run_closed"}
	reorderedDigest, err := reordered.Digest()
	if err != nil {
		t.Fatalf("reordered digest: %v", err)
	}
	if digest != reorderedDigest {
		t.Fatalf("digest should be order-stable: %s != %s", digest, reorderedDigest)
	}
}

func TestValidateRejectsMissingRequiredFields(t *testing.T) {
	valid := ExpectedEvidenceContract{
		SchemaVersion:         "v1",
		ContractID:            "contract-a",
		Version:               "2026-05-10",
		ContractSource:        "spec",
		LockRequiredBefore:    "implementation",
		RequiredObservers:     []string{"observer"},
		OptionalObservers:     []string{"optional"},
		RequiredEvents:        []string{"run_started"},
		GateEvents:            []string{"gate_checked"},
		MinimumGateTrustScope: "local_observed",
		RetentionProfile:      "digest_only",
		RedactionProfile:      "safe",
	}

	tests := []struct {
		name      string
		mutate    func(*ExpectedEvidenceContract)
		wantError string
	}{
		{
			name:      "schema_version required",
			mutate:    func(c *ExpectedEvidenceContract) { c.SchemaVersion = "" },
			wantError: "schema_version is required",
		},
		{
			name:      "contract_id required",
			mutate:    func(c *ExpectedEvidenceContract) { c.ContractID = "" },
			wantError: "contract_id is required",
		},
		{
			name:      "version required",
			mutate:    func(c *ExpectedEvidenceContract) { c.Version = "" },
			wantError: "version is required",
		},
		{
			name:      "contract_source required",
			mutate:    func(c *ExpectedEvidenceContract) { c.ContractSource = "" },
			wantError: "contract_source is required",
		},
		{
			name:      "lock_required_before required",
			mutate:    func(c *ExpectedEvidenceContract) { c.LockRequiredBefore = "" },
			wantError: "lock_required_before is required",
		},
		{
			name:      "required_observer required",
			mutate:    func(c *ExpectedEvidenceContract) { c.RequiredObservers = nil },
			wantError: "at least one required_observer is required",
		},
		{
			name:      "required_event required",
			mutate:    func(c *ExpectedEvidenceContract) { c.RequiredEvents = nil },
			wantError: "at least one required_event is required",
		},
		{
			name:      "gate_event required",
			mutate:    func(c *ExpectedEvidenceContract) { c.GateEvents = nil },
			wantError: "at least one gate_event is required",
		},
		{
			name:      "minimum_gate_trust_scope required",
			mutate:    func(c *ExpectedEvidenceContract) { c.MinimumGateTrustScope = "" },
			wantError: "minimum_gate_trust_scope is required",
		},
		{
			name:      "retention_profile required",
			mutate:    func(c *ExpectedEvidenceContract) { c.RetentionProfile = "" },
			wantError: "retention_profile is required",
		},
		{
			name:      "redaction_profile required",
			mutate:    func(c *ExpectedEvidenceContract) { c.RedactionProfile = "" },
			wantError: "redaction_profile is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contract := valid
			tt.mutate(&contract)
			if err := contract.Validate(); err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantError)
			} else if err.Error() != tt.wantError {
				t.Fatalf("expected %q, got %q", tt.wantError, err.Error())
			}
		})
	}
}

func TestValidateAllowsValidContract(t *testing.T) {
	contract := ExpectedEvidenceContract{
		SchemaVersion:         "v1",
		ContractID:            "contract-a",
		Version:               "2026-05-10",
		ContractSource:        "spec",
		LockRequiredBefore:    "implementation",
		RequiredObservers:     []string{"observer"},
		RequiredEvents:        []string{"run_started"},
		GateEvents:            []string{"gate_checked"},
		MinimumGateTrustScope: "local_observed",
		RetentionProfile:      "digest_only",
		RedactionProfile:      "safe",
	}
	if err := contract.Validate(); err != nil {
		t.Fatalf("expected valid contract, got %v", err)
	}
}
