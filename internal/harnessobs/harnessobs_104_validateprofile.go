package harnessobs

func validateProfile(profile Profile) error {
	// validateProfile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := validateProfileMetadata(profile); err != nil {
		return err
	}

	if err := validateProfileEventFamilies(profile.RequiredEventFamilies, profile.OptionalEventFamilies); err != nil {
		return err
	}
	return validateProfileDegradationRules(profile.DegradationRules)
}
