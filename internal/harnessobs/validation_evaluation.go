package harnessobs

// evaluationFromRun converts a readable observed run into a validation result;
// unreadable run evidence remains cannot_verify rather than becoming failure.
func evaluationFromRun(profile Profile, runDir string) Validation {
	run, events, err := LoadRun(runDir)
	if err != nil {
		return fallbackSourceUnavailable(profile)
	}
	return evaluate(profile, run, events)
}

// fallbackSourceUnavailable builds the evidence-only cannot_verify validation
// used when the observed run cannot be replayed.
func fallbackSourceUnavailable(profile Profile) Validation {
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
