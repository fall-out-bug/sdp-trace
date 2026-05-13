package policy

import (
	"encoding/json"
	"os"
)

// LoadPolicy loads and validates a policy file.
func LoadPolicy(path string) (AuthorityPolicy, error) {
	// Policy loading is file-backed; callers must supply a concrete authority
	// artifact before validation can establish the local policy boundary.
	raw, err := os.ReadFile(path)
	if err != nil {
		return AuthorityPolicy{}, err
	}
	// Policy bytes are decoded before validation so callers never receive a
	// partially trusted authority policy.
	var loaded AuthorityPolicy
	if err := json.Unmarshal(raw, &loaded); err != nil {
		return AuthorityPolicy{}, err
	}
	if err := loaded.Validate(); err != nil {
		return AuthorityPolicy{}, err
	}
	return loaded, nil
}

// AuthorityPolicy is the loaded and validated authority policy for adapters, signers, and witness scopes.
type AuthorityPolicy struct {
	SchemaVersion          string                  `json:"schema_version"`
	PolicyID               string                  `json:"policy_id"`
	AuthorityDescription   string                  `json:"authority_description"`
	AllowedAdapters        []AdapterAuthorityEntry `json:"allowed_adapters"`
	AllowedSigners         []SignerAuthorityEntry  `json:"allowed_signers"`
	AllowedWitnessProfiles []string                `json:"allowed_witness_profiles"`
	DemonstrationProfiles  []string                `json:"demonstration_profiles"`
	TrustBoundaryDefaults  TrustBoundaryDefaults   `json:"trust_boundary_defaults"`
}
