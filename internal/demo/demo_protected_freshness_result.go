package demo

func witnessFreshnessCannotVerify(code, reason, next string) ProtectedCondition {
	return witnessFreshnessCondition(GateCannotVerify, code, reason, next)
}

func witnessFreshnessFail(code, reason, next string) ProtectedCondition {
	return witnessFreshnessCondition(GateFail, code, reason, next)
}
func witnessFreshnessCondition(state, code, reason, next string) ProtectedCondition {
	// witnessFreshnessCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return ProtectedCondition{
		ID:         "witness_freshness_valid",
		State:      state,
		ReasonCode: code,
		Reason:     reason,
		NextAction: next,
	}
}
