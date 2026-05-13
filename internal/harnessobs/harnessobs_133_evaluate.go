package harnessobs

func evaluate(profile Profile, run Run, events []Event) Validation {
	// evaluate keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	dimensions := evaluationDimensions(profile, events)
	state, reason := compose(dimensions)
	if run.EventSchemaVersion != profile.EventSchemaVersion {

		state, reason = StateCannotVerify, "schema_version_mismatch"
	}
	return validationFromEvaluation(profile, dimensions, len(events), state, reason)
}
