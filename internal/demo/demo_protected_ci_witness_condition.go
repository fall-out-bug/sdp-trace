package demo

import (
	"strings"
)

func protectedCIWitnessCondition(input ProtectedGateInput) ProtectedCondition {
	// protectedCIWitnessCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if input.Witness == nil {
		return missingCIWitnessCondition()
	}
	state, reasons := witnessBindingState(*input.Witness, input.WitnessExpectation)
	code, reason, next := protectedCIWitnessFields(state, reasons)
	return ProtectedCondition{ID: "ci_witness_bound", State: state, ReasonCode: code, Reason: reason, NextAction: next}
}

func protectedCIWitnessFields(state string, reasons []string) (string, string, string) {
	// protectedCIWitnessFields keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if state == GatePass {

		return "ci_witness_bound", "CI witness source and artifact bindings match protected profile input", ""
	}
	return protectedCIWitnessNonPassFields(state, reasons)
}

func protectedCIWitnessNonPassFields(state string, reasons []string) (string, string, string) {
	// protectedCIWitnessNonPassFields keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	reason := strings.Join(reasons, "; ")
	if state == GateFail {
		return "ci_witness_mismatch", reason, "Fix the CI witness source or artifact binding mismatch."
	}
	return "ci_witness_incomplete", reason, "Supply complete CI witness source and artifact bindings."
}

func missingCIWitnessCondition() ProtectedCondition {
	// missingCIWitnessCondition keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return ProtectedCondition{
		ID:         "ci_witness_bound",
		State:      GateCannotVerify,
		ReasonCode: "missing_ci_witness",
		Reason:     "CI witness evidence is required for protected profile",
		NextAction: "Supply a CI witness bound to the selected run.",
	}
}
