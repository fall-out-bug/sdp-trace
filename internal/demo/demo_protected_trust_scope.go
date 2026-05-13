package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func protectedTrustScopeCondition(input ProtectedGateInput) ProtectedCondition {
	// protectedTrustScopeCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if !input.PolicyProvided {
		return ProtectedCondition{
			ID:         "protected_trust_scope_satisfied",
			State:      GateCannotVerify,
			ReasonCode: "missing_policy",
			Reason:     "protected trust scope cannot be verified without trusted-checkpoint policy",
			NextAction: "Supply a trusted-checkpoint policy for the protected signer.",
		}
	}
	if protectedCheckpointCanUseWitness(input) {
		return protectedWitnessTrustScopeCondition(input)
	}
	return protectedInsufficientTrustScopeCondition(input.Checkpoint.TrustScope)
}

func protectedCheckpointCanUseWitness(input ProtectedGateInput) bool {
	// protectedCheckpointCanUseWitness keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return input.Checkpoint.Result == checkpoint.StatePass &&
		input.Checkpoint.TrustScope == checkpoint.TrustScopeCISigned &&
		input.Checkpoint.SignerAuthorityState == checkpoint.StatePass &&
		input.Witness != nil
}
