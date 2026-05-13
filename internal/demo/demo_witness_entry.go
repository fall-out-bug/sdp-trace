package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func EvaluateGateWithWitness(rows []RunRow, contract trace.Contract, witnessPath string) GateResult {

	return applyWitness(EvaluateGate(rows, contract), witnessPath)
}

func EvaluateGateWithWitnessContext(rows []RunRow, contract trace.Contract, witnessPath string, expected WitnessExpectation) GateResult {

	return applyWitnessWithExpectation(EvaluateGate(rows, contract), witnessPath, expected)
}

func PreviewWitnessBinding(witnessPath, target string) (bool, []string) {
	// PreviewWitnessBinding keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	record, err := loadWitnessSummary(witnessPath)
	if err != nil {
		return false, []string{err.Error()}
	}
	expected, err := witnessExpectationFromTarget(target)
	if err != nil {

		return true, []string{err.Error()}
	}
	state, reasons := witnessBindingState(record, expected)
	if state == GatePass {
		return true, []string{}
	}
	return true, reasons
}
