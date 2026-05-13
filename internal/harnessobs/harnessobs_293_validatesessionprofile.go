package harnessobs

func validateSessionProfile(profile *SessionProfile) error {
	// validateSessionProfile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	if err := validateSessionProfileIdentity(*profile); err != nil {
		return err
	}

	if err := normalizeSessionStreamCapture(profile); err != nil {
		return err
	}
	if err := validateSessionSetupActions(profile.SetupActions); err != nil {
		return err
	}
	return validateSessionIsolationRules(profile.IsolationRules)
}
