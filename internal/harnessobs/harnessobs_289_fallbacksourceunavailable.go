package harnessobs

func fallbackSourceUnavailable(profile Profile) Validation {
	// fallbackSourceUnavailable keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	validation := Validation{
		SchemaVersion:      ValidationSchemaVersion,
		ProfileID:          profile.ProfileID,
		HarnessFamily:      profile.HarnessFamily,
		EventSchemaVersion: profile.EventSchemaVersion,
		ValidationState:    StateCannotVerify,
		ReasonCode:         "source_unavailable",
		NonAuthority:       nonAuthority(),
	}
	validation.ValidationDigest = validationDigest(validation)

	return validation
}
