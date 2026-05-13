package contract

import (
	"encoding/json"
	"os"
)

// ExpectedEvidenceContract is the parsed form of schema/expected-evidence-contract.schema.json.
type ExpectedEvidenceContract struct {
	SchemaVersion         string   `json:"schema_version"`
	ContractID            string   `json:"contract_id"`
	Version               string   `json:"version"`
	ContractSource        string   `json:"contract_source"`
	LockRequiredBefore    string   `json:"lock_required_before"`
	RequiredObservers     []string `json:"required_observers"`
	OptionalObservers     []string `json:"optional_observers"`
	RequiredEvents        []string `json:"required_events"`
	GateEvents            []string `json:"gate_events"`
	MinimumGateTrustScope string   `json:"minimum_gate_trust_scope"`
	RetentionProfile      string   `json:"retention_profile"`
	RedactionProfile      string   `json:"redaction_profile"`
}

// Load parses a contract from JSON and validates required fields.
func Load(path string) (ExpectedEvidenceContract, error) {
	// Contracts are loaded from a concrete JSON artifact before validation; an
	// absent or unreadable file remains a structural evidence gap for callers.
	raw, err := os.ReadFile(path)
	if err != nil {
		return ExpectedEvidenceContract{}, err
	}
	// Decode and validate before returning so callers never receive a partially
	// trusted expected-evidence contract.
	var contract ExpectedEvidenceContract
	if err := json.Unmarshal(raw, &contract); err != nil {
		return ExpectedEvidenceContract{}, err
	}
	if err := contract.Validate(); err != nil {
		return ExpectedEvidenceContract{}, err
	}
	return contract, nil
}
