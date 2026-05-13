package harnessobs

func validateProfileEventFamilies(requiredEventFamilies []string, optionalEventFamilies []string) error {
	// validateProfileEventFamilies keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := validateFamilyList(requiredEventFamilies); err != nil {
		return err
	}

	return validateFamilyList(optionalEventFamilies)
}
