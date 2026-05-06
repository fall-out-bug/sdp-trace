package contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
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
	raw, err := os.ReadFile(path)
	if err != nil {
		return ExpectedEvidenceContract{}, err
	}
	var contract ExpectedEvidenceContract
	if err := json.Unmarshal(raw, &contract); err != nil {
		return ExpectedEvidenceContract{}, err
	}
	if err := contract.Validate(); err != nil {
		return ExpectedEvidenceContract{}, err
	}
	return contract, nil
}

// Validate checks required fields and basic cardinality constraints.
func (c ExpectedEvidenceContract) Validate() error {
	if c.SchemaVersion == "" {
		return fmt.Errorf("schema_version is required")
	}
	if c.ContractID == "" {
		return fmt.Errorf("contract_id is required")
	}
	if c.Version == "" {
		return fmt.Errorf("version is required")
	}
	if c.ContractSource == "" {
		return fmt.Errorf("contract_source is required")
	}
	if c.LockRequiredBefore == "" {
		return fmt.Errorf("lock_required_before is required")
	}
	if len(c.RequiredObservers) == 0 {
		return fmt.Errorf("at least one required_observer is required")
	}
	if len(c.RequiredEvents) == 0 {
		return fmt.Errorf("at least one required_event is required")
	}
	if len(c.GateEvents) == 0 {
		return fmt.Errorf("at least one gate_event is required")
	}
	if c.MinimumGateTrustScope == "" {
		return fmt.Errorf("minimum_gate_trust_scope is required")
	}
	if c.RetentionProfile == "" {
		return fmt.Errorf("retention_profile is required")
	}
	if c.RedactionProfile == "" {
		return fmt.Errorf("redaction_profile is required")
	}
	return nil
}

// Digest computes a deterministic digest for fixture signing and lock checks.
func (c ExpectedEvidenceContract) Digest() (string, error) {
	canonical, err := canonicalizeContract(c)
	if err != nil {
		return "", err
	}
	return canonicalDigest(canonical), nil
}

func canonicalizeContract(contract ExpectedEvidenceContract) ([]byte, error) {
	stable := contract
	sort.Strings(stable.RequiredObservers)
	sort.Strings(stable.OptionalObservers)
	sort.Strings(stable.RequiredEvents)
	sort.Strings(stable.GateEvents)
	canonical, err := trace.CanonicalJSON(stable)
	if err != nil {
		return nil, err
	}
	return canonical, nil
}

func canonicalDigest(data []byte) string {
	sum := sha256.Sum256(data)
	return strings.ToLower(hex.EncodeToString(sum[:]))
}

