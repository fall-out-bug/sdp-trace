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
	return firstValidationError(
		func() error { return validateRequiredString(c.SchemaVersion, "schema_version") },
		func() error { return validateRequiredString(c.ContractID, "contract_id") },
		func() error { return validateRequiredString(c.Version, "version") },
		func() error { return validateRequiredString(c.ContractSource, "contract_source") },
		func() error { return validateRequiredString(c.LockRequiredBefore, "lock_required_before") },
		func() error { return validateNonEmptyList(c.RequiredObservers, "required_observer") },
		func() error { return validateNonEmptyList(c.RequiredEvents, "required_event") },
		func() error { return validateNonEmptyList(c.GateEvents, "gate_event") },
		func() error { return validateRequiredString(c.MinimumGateTrustScope, "minimum_gate_trust_scope") },
		func() error { return validateRequiredString(c.RetentionProfile, "retention_profile") },
		func() error { return validateRequiredString(c.RedactionProfile, "redaction_profile") },
	)
}

func validateRequiredString(value, name string) error {
	if value == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

func validateNonEmptyList(values []string, name string) error {
	if len(values) == 0 {
		return fmt.Errorf("at least one %s is required", name)
	}
	return nil
}

func firstValidationError(checks ...func() error) error {
	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
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
