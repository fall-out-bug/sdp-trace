package harnessobs

func evaluate(profile Profile, run Run, events []Event) Validation {
	// Evaluation keeps replayed event evidence separate from the final
	// validation artifact, then applies schema mismatch as an explicit
	// cannot_verify override.
	dimensions := evaluationDimensions(profile, events)
	state, reason := compose(dimensions)
	if run.EventSchemaVersion != profile.EventSchemaVersion {
		state, reason = StateCannotVerify, "schema_version_mismatch"
	}
	return validationFromEvaluation(profile, dimensions, len(events), state, reason)
}

func validationFromEvaluation(profile Profile, dimensions []Dimension, eventCount int, state, reason string) Validation {
	// The digest is calculated after all validation fields are assigned so the
	// artifact remains replay-bound to the exact emitted validation content.
	validation := Validation{
		SchemaVersion:      ValidationSchemaVersion,
		ProfileID:          profile.ProfileID,
		HarnessFamily:      profile.HarnessFamily,
		EventSchemaVersion: profile.EventSchemaVersion,
		ValidationState:    state,
		ReasonCode:         reason,
		Dimensions:         dimensions,
		EventCount:         eventCount,
		NonAuthority:       nonAuthority(),
	}

	validation.ValidationDigest = validationDigest(validation)
	return validation
}
