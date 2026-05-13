package demo

func protectedConditionFromGateCondition(conditions []GateCondition, id string) ProtectedCondition {
	// protectedConditionFromGateCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, condition := range conditions {
		if condition.ID == id {
			return protectedConditionFromLocalGate(id, condition)
		}
	}
	return ProtectedCondition{
		ID:         id,
		State:      GateCannotVerify,
		ReasonCode: id + "_missing",
		Reason:     "required gate condition is missing",
		NextAction: "Regenerate the gate result with current sdp-trace.",
	}
}

func protectedConditionFromLocalGate(id string, condition GateCondition) ProtectedCondition {
	// protectedConditionFromLocalGate keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	code := "condition_pass"
	next := ""
	if condition.State != GatePass {

		code = id + "_not_satisfied"
		next = "Supply the required run and evidence before evaluating protected profile."
	}

	return ProtectedCondition{ID: id, State: condition.State, ReasonCode: code, Reason: condition.Reason, NextAction: next}
}
