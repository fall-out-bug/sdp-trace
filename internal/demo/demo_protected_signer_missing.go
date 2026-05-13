package demo

func missingSignerPolicyCondition() ProtectedCondition {
	// missingSignerPolicyCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return ProtectedCondition{
		ID:         "checkpoint_signer_authorized",
		State:      GateCannotVerify,
		ReasonCode: "missing_policy",
		Reason:     "trusted-checkpoint policy is required for protected profile",
		NextAction: "Supply a trusted-checkpoint policy for the protected signer.",
	}
}
