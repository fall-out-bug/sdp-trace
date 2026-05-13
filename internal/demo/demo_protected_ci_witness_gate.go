package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func protectedCIWitnessGate(input ProtectedGateInput) string {
	// protectedCIWitnessGate keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if input.Witness == nil {
		return GateCannotVerify
	}
	state, _ := witnessBindingState(*input.Witness, input.WitnessExpectation)
	return state
}

func protectedTrustCap(input ProtectedGateInput, ciWitnessGate string) string {
	// protectedTrustCap keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if trustCap := protectedCheckpointTrustCap(input.Checkpoint.TrustScope); trustCap != "" {
		return trustCap
	}
	return nonCheckpointProtectedTrustCap(input.Checkpoint.TrustScope, ciWitnessGate)
}
func nonCheckpointProtectedTrustCap(checkpointTrustScope, ciWitnessGate string) string {
	// nonCheckpointProtectedTrustCap keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if ciWitnessGate == GatePass {
		return "ci_witnessed"
	}
	if checkpointTrustScope != "" {
		return checkpointTrustScope
	}
	return string(trace.TrustScopeLocalObserved)
}

func protectedCheckpointTrustCap(trustScope string) string {
	// protectedCheckpointTrustCap keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, candidate := range []string{checkpoint.TrustScopeCISigned, checkpoint.TrustScopeLocalSigned} {
		if trustScope == candidate {
			return candidate
		}
	}
	return ""
}
