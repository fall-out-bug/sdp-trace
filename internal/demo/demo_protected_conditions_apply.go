package demo

func applyProtectedConditionResults(result *GateResult, input ProtectedGateInput) {
	// applyProtectedConditionResults keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	result.ProtectedConditions = protectedConditions(*result, input)
	for _, condition := range result.ProtectedConditions {

		if condition.ID == "override_does_not_upgrade_profile" {
			continue
		}
		result.ProtectedGate = worseProtectedState(result.ProtectedGate, topLevelProtectedState(condition.State))
	}
	result.Reasons = append(result.Reasons, protectedReasons(result.ProtectedConditions)...)
	result.NextActions = append(result.NextActions, protectedNextActions(result.ProtectedConditions)...)
}
