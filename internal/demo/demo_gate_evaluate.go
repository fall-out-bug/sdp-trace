package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func EvaluateGate(rows []RunRow, contract trace.Contract) GateResult {
	// EvaluateGate keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result := newGateResult(rows, contract)
	observedEvidence := applyRunRows(&result, rows)
	applyRequiredRuns(&result, rows, contract)
	applyRequiredEvidence(&result, contract, observedEvidence)
	finalizeGateResult(&result)
	return result
}
