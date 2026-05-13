package harnessobs

func validateEvent(profile Profile, event Event) error {
	// validateEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := validateEventIdentity(profile, event); err != nil {
		return err
	}

	if err := validateEventRefs(event); err != nil {
		return err
	}
	if err := validateEventContent(event); err != nil {
		return err
	}
	return validateUnavailableFields(event.UnavailableFields)
}
