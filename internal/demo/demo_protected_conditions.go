package demo

func protectedConditions(result GateResult, input ProtectedGateInput) []ProtectedCondition {
	// protectedConditions keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return []ProtectedCondition{
		protectedProfileSelectedCondition(),

		protectedConditionFromGateCondition(result.GateConditions, "all_required_runs_present"),
		protectedConditionFromGateCondition(result.GateConditions, "all_required_evidence_observed"),

		protectedCIWitnessCondition(input),
		protectedWitnessFreshnessCondition(input),

		protectedCheckpointSignatureCondition(input.Checkpoint),
		protectedCheckpointBindingCondition(input.Checkpoint),
		protectedSignerCondition(input),
		protectedTrustScopeCondition(input),
		protectedOverrideCondition(result.OverrideRequests),
	}
}

func protectedProfileSelectedCondition() ProtectedCondition {
	// protectedProfileSelectedCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return ProtectedCondition{
		ID:         "protected_profile_explicitly_selected",
		State:      GatePass,
		ReasonCode: "protected_profile_selected",
		Reason:     "protected profile was explicitly selected",
	}
}
