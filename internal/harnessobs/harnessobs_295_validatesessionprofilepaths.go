package harnessobs

func validateSessionProfilePaths(profile SessionProfile) error {
	// validateSessionProfilePaths keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := validateRequiredSessionPaths(profile); err != nil {
		return err
	}

	return validateRawEventConfig(profile)
}
