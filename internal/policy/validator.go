package policy

type AuthorityPolicyValidator struct {
	policy AuthorityPolicy
}

// NewAuthorityPolicyValidator builds a read-only validator.
func NewAuthorityPolicyValidator(policy AuthorityPolicy) AuthorityPolicyValidator {
	return AuthorityPolicyValidator{policy: policy}
}

// WitnessProfileAllowed checks whether witness profile is trusted by policy.
func (v AuthorityPolicyValidator) WitnessProfileAllowed(profileID string) bool {
	return stringInList(v.policy.AllowedWitnessProfiles, profileID)
}

// IsDemonstrationProfile marks a local-only profile as not production-grade.
func (v AuthorityPolicyValidator) IsDemonstrationProfile(profileID string) bool {
	return stringInList(v.policy.DemonstrationProfiles, profileID)
}

func (signer SignerAuthorityEntry) ScopeAllowed(scope string) bool {
	return stringInList(signer.AllowedScopes, scope)
}

func stringInList(values []string, needle string) bool {
	// Authority policy values are closed vocabulary strings; substring or glob
	// matching would silently widen the trust boundary.
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}
