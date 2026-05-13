package ciartifact

func evaluateSafety(input OutputSafetyInput) OutputSafetyResult {
	// Safety evaluation scans source and run identity fields before they are copied
	// into human-facing reasons or machine-readable refs.

	state := defaultString(input.State, StateNotAssessed)
	outcome, ok := safetyOutcomes[state]
	if !ok {

		state = StateCannotVerify
		outcome = familyOutcome{StateCannotVerify, "output_safety_cannot_verify", "observation output safety state is unrecognized under selected profile", ""}
	}
	return OutputSafetyResult{State: state, UnsafeClasses: safeClasses(input.UnsafeClasses), ReasonCode: outcome.code, Reason: outcome.reason}
}

var safetyOutcomes = map[string]familyOutcome{
	StatePass:         {StatePass, "output_safety_pass", "observation output safety classes are absent", ""},
	StateFail:         {StateFail, "unsafe_artifact_output", "observation detected forbidden output-safety classes", ""},
	StateCannotVerify: {StateCannotVerify, "output_safety_cannot_verify", "observation output safety could not be verified", ""},
	StateNotAssessed:  {StateNotAssessed, "output_safety_not_assessed", "output safety was outside selected profile scope", ""},
}
