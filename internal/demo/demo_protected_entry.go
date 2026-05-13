package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func EvaluateProtectedGate(rows []RunRow, contract trace.Contract, input ProtectedGateInput) GateResult {
	// EvaluateProtectedGate keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result := EvaluateGate(rows, contract)
	applyProtectedGateContext(&result, input)
	applyProtectedConditionResults(&result, input)

	return result
}

func applyProtectedGateContext(result *GateResult, input ProtectedGateInput) {
	// applyProtectedGateContext keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result.SchemaVersion = GateSchemaVersionBlock16
	result.SelectedProfile = GateProfileProtected
	result.ProtectedGate = GatePass
	result.GateMode = GateProfileProtected
	result.Witness = input.Witness
	result.CIWitnessGate = protectedCIWitnessGate(input)
	result.TrustCap = protectedTrustCap(input, result.CIWitnessGate)

	result.GateConditions = gateConditions(*result)
	result.CheckpointVerification = &input.Checkpoint
}
