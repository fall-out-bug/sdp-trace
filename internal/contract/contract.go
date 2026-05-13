package contract

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"os"
	"sort"
	"strings"
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

// Validate checks required fields and basic cardinality constraints.
func (c ExpectedEvidenceContract) Validate() error {

	return firstValidationError(
		func() error { return validateStringFields(c, contractHeaderFields) },
		func() error { return validateListFields(c, contractEventSetFields) },
		func() error { return validateStringFields(c, contractPolicyFields) },
	)
}

// Digest computes a deterministic digest for fixture signing and lock checks.
func (c ExpectedEvidenceContract) Digest() (string, error) {
	canonical, err := canonicalizeContract(c)
	if err != nil {
		return "", err
	}

	return canonicalDigest(canonical), nil
}

type contractStringField struct {
	name  string
	value func(ExpectedEvidenceContract) string
}

type contractListField struct {
	name  string
	value func(ExpectedEvidenceContract) []string
}

var contractHeaderFields = []contractStringField{
	{"schema_version", func(c ExpectedEvidenceContract) string { return c.SchemaVersion }},
	{"contract_id", func(c ExpectedEvidenceContract) string { return c.ContractID }},
	{"version", func(c ExpectedEvidenceContract) string { return c.Version }},
	{"contract_source", func(c ExpectedEvidenceContract) string { return c.ContractSource }},
	{"lock_required_before", func(c ExpectedEvidenceContract) string { return c.LockRequiredBefore }},
}

var contractEventSetFields = []contractListField{
	{"required_observer", func(c ExpectedEvidenceContract) []string { return c.RequiredObservers }},
	{"required_event", func(c ExpectedEvidenceContract) []string { return c.RequiredEvents }},
	{"gate_event", func(c ExpectedEvidenceContract) []string { return c.GateEvents }},
}

var contractPolicyFields = []contractStringField{
	{"minimum_gate_trust_scope", func(c ExpectedEvidenceContract) string { return c.MinimumGateTrustScope }},
	{"retention_profile", func(c ExpectedEvidenceContract) string { return c.RetentionProfile }},
	{"redaction_profile", func(c ExpectedEvidenceContract) string { return c.RedactionProfile }},
}

func validateStringFields(contract ExpectedEvidenceContract, fields []contractStringField) error {

	for _, field := range fields {
		if err := validateRequiredString(field.value(contract), field.name); err != nil {

			return err
		}
	}
	return nil
}

func validateListFields(contract ExpectedEvidenceContract, fields []contractListField) error {

	for _, field := range fields {
		if err := validateNonEmptyList(field.value(contract), field.name); err != nil {

			return err
		}
	}
	return nil
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
