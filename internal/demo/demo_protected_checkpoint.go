package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func protectedCheckpointSignatureCondition(result checkpoint.VerificationResult) ProtectedCondition {
	// protectedCheckpointSignatureCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if result.SignatureState == checkpoint.StatePass && result.PayloadDigestState != checkpoint.StateFail {
		return ProtectedCondition{ID: "checkpoint_signature_valid", State: GatePass, ReasonCode: "checkpoint_signature_valid", Reason: "checkpoint signature verification passed"}
	}
	return ProtectedCondition{
		ID:         "checkpoint_signature_valid",
		State:      mapCheckpointState(result.SignatureState),
		ReasonCode: "checkpoint_signature_invalid",
		Reason:     "checkpoint signature verification did not pass",
		NextAction: "Regenerate the signed checkpoint for the selected run.",
	}
}

func protectedCheckpointBindingCondition(result checkpoint.VerificationResult) ProtectedCondition {
	// protectedCheckpointBindingCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	state := GatePass
	for _, candidate := range []string{result.RunBindingState, result.ChainBindingState, result.SourceBindingState, result.NonceBindingState} {

		state = worseProtectedState(state, mapCheckpointState(candidate))
	}
	if state == GatePass {
		return ProtectedCondition{ID: "checkpoint_run_binding_valid", State: GatePass, ReasonCode: "checkpoint_binding_valid", Reason: "checkpoint binding matches the selected run context"}
	}
	return ProtectedCondition{
		ID:         "checkpoint_run_binding_valid",
		State:      state,
		ReasonCode: "checkpoint_binding_invalid",
		Reason:     "checkpoint binding does not satisfy the selected run context",
		NextAction: "Regenerate checkpoint evidence from the selected run context.",
	}
}
