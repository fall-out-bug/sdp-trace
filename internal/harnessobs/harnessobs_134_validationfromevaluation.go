package harnessobs

func validationFromEvaluation(profile Profile, dimensions []Dimension, eventCount int, state, reason string) Validation {
	// validationFromEvaluation keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	validation := Validation{
		SchemaVersion:      ValidationSchemaVersion,
		ProfileID:          profile.ProfileID,
		HarnessFamily:      profile.HarnessFamily,
		EventSchemaVersion: profile.EventSchemaVersion,

		ValidationState: state,
		ReasonCode:      reason,
		Dimensions:      dimensions,
		EventCount:      eventCount,
		NonAuthority:    nonAuthority(),
	}

	validation.ValidationDigest = validationDigest(validation)
	return validation
}
