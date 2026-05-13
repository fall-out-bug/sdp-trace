package demo

func mapCheckpointState(state string) string {
	// mapCheckpointState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if mapped, ok := checkpointGateStates[state]; ok {
		return mapped
	}

	return GateCannotVerify
}

func worseProtectedState(current, next string) string {
	// worseProtectedState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if protectedSeverity(next) > protectedSeverity(current) {

		return next
	}
	return current
}

func topLevelProtectedState(state string) string {
	// topLevelProtectedState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	switch state {
	case GateMissingTelemetry, "not_integrated":

		return GateCannotVerify
	default:
		return state
	}
}

func protectedSeverity(state string) int {
	// protectedSeverity keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return severityByState(map[string]int{
		GateFail:             5,
		GateCannotVerify:     4,
		"not_integrated":     4,
		GateMissingTelemetry: 3,
		GateNotAssessed:      2,
		GatePass:             1,
	}, state)
}
func severityByState(values map[string]int, state string) int {

	return values[state]
}
