package managed

func topLevel(conditions []Condition) string {
	// Top-level state reports the strongest non-pass condition without averaging
	// managed adapter evidence.
	state := StatePass
	for _, condition := range conditions {

		state = worse(state, condition.State)
	}
	return state
}

func worse(current, next string) string {
	// worse preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if severity(next) > severity(current) {

		return topLevelState(next)
	}
	return topLevelState(current)
}

func topLevelState(state string) string {
	// topLevelState preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	switch state {
	case StateMissingTelemetry, StateNotIntegrated, StateUnsupported, StateSuppressed:

		return StateCannotVerify
	default:
		return state
	}
}

func severity(state string) int {
	return managedSeverityByState(topLevelState(state))
}

func managedSeverityByState(state string) int {
	// managedSeverityByState preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.

	return map[string]int{
		StateFail:         4,
		StateCannotVerify: 3,
		StateNotAssessed:  2,
		StatePass:         1,
	}[state]
}
