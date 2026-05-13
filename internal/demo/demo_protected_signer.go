package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func protectedSignerCondition(input ProtectedGateInput) ProtectedCondition {
	// protectedSignerCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if !input.PolicyProvided {
		return missingSignerPolicyCondition()
	}
	state := mapCheckpointState(input.Checkpoint.SignerAuthorityState)

	if protectedSignerPass(state, input.Checkpoint.TrustScope) {

		return ProtectedCondition{ID: "checkpoint_signer_authorized", State: GatePass, ReasonCode: "checkpoint_signer_authorized", Reason: "checkpoint signer is authorized for CI signed protected profile"}
	}
	if protectedSignerLocalOnly(state, input.Checkpoint.TrustScope) {

		state = GateFail
	}
	return ProtectedCondition{
		ID:         "checkpoint_signer_authorized",
		State:      state,
		ReasonCode: "checkpoint_signer_not_protected",
		Reason:     "checkpoint signer authority does not satisfy protected profile",
		NextAction: "Run checkpoint signing in an authorized CI signer context.",
	}
}

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
func protectedSignerPass(state, trustScope string) bool {
	return state == GatePass && trustScope == checkpoint.TrustScopeCISigned
}

func protectedSignerLocalOnly(state, trustScope string) bool {
	return state == GatePass && trustScope == checkpoint.TrustScopeLocalSigned
}
