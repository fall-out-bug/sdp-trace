package policy

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// AuthorityPolicy is the loaded and validated authority policy for adapters, signers, and witness scopes.
type AuthorityPolicy struct {
	SchemaVersion          string                    `json:"schema_version"`
	PolicyID               string                    `json:"policy_id"`
	AuthorityDescription   string                    `json:"authority_description"`
	AllowedAdapters        []AdapterAuthorityEntry   `json:"allowed_adapters"`
	AllowedSigners         []SignerAuthorityEntry    `json:"allowed_signers"`
	AllowedWitnessProfiles []string                  `json:"allowed_witness_profiles"`
	DemonstrationProfiles  []string                  `json:"demonstration_profiles"`
	TrustBoundaryDefaults  TrustBoundaryDefaults     `json:"trust_boundary_defaults"`
}

type AdapterAuthorityEntry struct {
	AdapterID       string   `json:"adapter_id"`
	Provider        string   `json:"provider"`
	IdentityState   string   `json:"identity_state"`
	AllowedEventTypes []string `json:"allowed_event_types"`
	AllowedByPolicy bool     `json:"allowed_by_policy"`
}

type SignerAuthorityEntry struct {
	SignerID             string   `json:"signer_id"`
	ProfileID            string   `json:"profile_id"`
	AllowedScopes        []string `json:"allowed_scopes"`
	IndependenceRequired bool     `json:"independence_required"`
	Environment          string   `json:"environment"`
}

type TrustBoundaryDefaults struct {
	DefaultWitnessIndependence string `json:"default_witness_independence"`
	LocalProfileLabel         string `json:"local_profile_label"`
}

type SigningProfileSpec struct {
	SchemaVersion string `json:"schema_version"`
	ProfileID    string `json:"profile_id"`
	Format       string `json:"format"`
	Description  string `json:"description"`
	VerifierHost string `json:"verifier_host"`
}

type RedactionProfileSpec struct {
	SchemaVersion string `json:"schema_version"`
	ProfileID    string `json:"profile_id"`
	Description  string `json:"description"`
	DefaultMode  string `json:"default_retention_mode"`
}

type AuthorityPolicyValidator struct {
	policy AuthorityPolicy
}

// LoadPolicy loads and validates a policy file.
func LoadPolicy(path string) (AuthorityPolicy, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return AuthorityPolicy{}, err
	}
	var loaded AuthorityPolicy
	if err := json.Unmarshal(raw, &loaded); err != nil {
		return AuthorityPolicy{}, err
	}
	if err := loaded.Validate(); err != nil {
		return AuthorityPolicy{}, err
	}
	return loaded, nil
}

// Validate checks policy invariants required before verifier use.
func (policy AuthorityPolicy) Validate() error {
	if policy.SchemaVersion == "" {
		return fmt.Errorf("schema_version is required")
	}
	if policy.PolicyID == "" {
		return fmt.Errorf("policy_id is required")
	}
	if len(policy.AllowedSigners) == 0 {
		return fmt.Errorf("at least one allowed signer is required")
	}
	if len(policy.AllowedWitnessProfiles) == 0 {
		return fmt.Errorf("at least one allowed witness profile is required")
	}
	return nil
}

// NewAuthorityPolicyValidator builds a read-only validator.
func NewAuthorityPolicyValidator(policy AuthorityPolicy) AuthorityPolicyValidator {
	return AuthorityPolicyValidator{policy: policy}
}

// CanAdapterEmit checks whether an adapter id can emit a given event type.
func (v AuthorityPolicyValidator) CanAdapterEmit(adapterID string, eventType trace.EventType) bool {
	for _, adapter := range v.policy.AllowedAdapters {
		if adapter.AdapterID != adapterID {
			continue
		}
		for _, allowedType := range adapter.AllowedEventTypes {
			if allowedType == string(eventType) && adapter.AllowedByPolicy {
				return true
			}
		}
	}
	return false
}

// SignerAllowed checks whether a signer/profile tuple is permitted.
func (v AuthorityPolicyValidator) SignerAllowed(signerID string, profileID string) bool {
	for _, signer := range v.policy.AllowedSigners {
		if signer.SignerID != signerID {
			continue
		}
		if !signer.ScopeAllowed("signer") {
			continue
		}
		if signer.ProfileID == profileID {
			return true
		}
	}
	return false
}

// WitnessProfileAllowed checks whether witness profile is trusted by policy.
func (v AuthorityPolicyValidator) WitnessProfileAllowed(profileID string) bool {
	for _, profile := range v.policy.AllowedWitnessProfiles {
		if profile == profileID {
			return true
		}
	}
	return false
}

// IsDemonstrationProfile marks a local-only profile as not production-grade.
func (v AuthorityPolicyValidator) IsDemonstrationProfile(profileID string) bool {
	for _, profile := range v.policy.DemonstrationProfiles {
		if profile == profileID {
			return true
		}
	}
	return false
}

func (signer SignerAuthorityEntry) ScopeAllowed(scope string) bool {
	for _, item := range signer.AllowedScopes {
		if item == scope {
			return true
		}
	}
	return false
}

