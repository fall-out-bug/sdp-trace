package policy

import "errors"

// Validate checks policy invariants required before verifier use.
func (policy AuthorityPolicy) Validate() error {
	// Validation only checks the policy shape needed before verifier use. It
	// does not authenticate signers or upgrade local policy into external proof.
	if err := policy.validateIdentity(); err != nil {
		return err
	}
	if len(policy.AllowedSigners) == 0 {
		return errors.New("at least one allowed signer is required")
	}
	if len(policy.AllowedWitnessProfiles) == 0 {
		return errors.New("at least one allowed witness profile is required")
	}

	return nil
}

func (policy AuthorityPolicy) validateIdentity() error {
	// Keep identity diagnostics before authority-list diagnostics; fixture
	// failures should identify malformed policy identity before trust scope.
	if policy.SchemaVersion == "" {
		return errors.New("schema_version is required")
	}
	if policy.PolicyID == "" {
		return errors.New("policy_id is required")
	}
	return nil
}
