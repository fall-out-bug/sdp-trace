package demo

import (
	"fmt"
)

func witnessBindingState(record WitnessSummary, expected WitnessExpectation) (string, []string) {
	// witnessBindingState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, binding := range witnessScalarBindings(record, expected) {
		if state, reasons := validateWitnessScalarBinding(binding); state != GatePass {
			return state, reasons
		}
	}
	return witnessArtifactBindingState(record.RunArtifacts, expected.RunArtifacts)
}

func witnessScalarBindings(record WitnessSummary, expected WitnessExpectation) []witnessScalarBinding {
	// witnessScalarBindings keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	return []witnessScalarBinding{
		{label: "repository", expected: expected.Repository, actual: record.Source.Repository},
		{label: "ref", expected: expected.Ref, actual: record.Source.Ref},
		{label: "commit", expected: expected.CommitSHA, actual: record.Source.CommitSHA},
		{label: "run id", expected: expected.RunID, actual: record.CIIdentity.RunID},
	}
}

func validateWitnessScalarBinding(binding witnessScalarBinding) (string, []string) {
	// validateWitnessScalarBinding keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if binding.expected == "" {
		return GatePass, nil
	}
	if binding.actual == "" {
		return GateCannotVerify, []string{fmt.Sprintf("ci witness %s binding is missing", binding.label)}
	}
	if binding.actual != binding.expected {
		return GateFail, []string{fmt.Sprintf("ci witness %s mismatch: expected %s got %s", binding.label, binding.expected, binding.actual)}
	}
	return GatePass, nil
}
