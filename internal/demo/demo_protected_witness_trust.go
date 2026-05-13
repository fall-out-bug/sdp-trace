package demo

func protectedWitnessTrustScopeCondition(input ProtectedGateInput) ProtectedCondition {
	// protectedWitnessTrustScopeCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	witnessState, _ := witnessBindingState(*input.Witness, input.WitnessExpectation)
	freshness := protectedWitnessFreshnessCondition(input)
	if witnessState == GatePass && freshness.State == GatePass {

		return ProtectedCondition{ID: "protected_trust_scope_satisfied", State: GatePass, ReasonCode: "protected_trust_scope_satisfied", Reason: "CI signed checkpoint and CI witness binding satisfy protected profile"}
	}
	state := worseProtectedState(witnessState, freshness.State)
	return ProtectedCondition{
		ID:         "protected_trust_scope_satisfied",
		State:      state,
		ReasonCode: "protected_trust_scope_not_satisfied",
		Reason:     "CI signed checkpoint does not have passing CI witness binding and freshness",
		NextAction: "Provide fresh CI witness binding for the selected run.",
	}
}
