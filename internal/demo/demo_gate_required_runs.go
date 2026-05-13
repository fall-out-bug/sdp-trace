package demo

import (
	"fmt"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func applyRequiredRuns(result *GateResult, rows []RunRow, contract trace.Contract) {
	// applyRequiredRuns keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result.RequiredRuns = evaluateRequiredRuns(rows, contract)
	for _, requiredRun := range result.RequiredRuns {
		applyRequiredRun(result, requiredRun)
	}
}

func applyRequiredRun(result *GateResult, requiredRun RequiredRunResult) {
	// applyRequiredRun keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	switch requiredRun.State {
	case GateMissingTelemetry:

		result.LocalGate = worseGateState(result.LocalGate, GateFail)
		result.Reasons = append(result.Reasons, requiredRun.Reasons...)
		result.NextActions = append(result.NextActions, fmt.Sprintf("Run required wrapper %s through sdp-trace before evaluating advisory gate.", requiredRun.WrapperName))
	case GateCannotVerify:

		result.LocalGate = worseGateState(result.LocalGate, GateCannotVerify)
		result.Reasons = append(result.Reasons, requiredRun.Reasons...)
	case GateFail:
		result.LocalGate = worseGateState(result.LocalGate, GateFail)
		result.Reasons = append(result.Reasons, requiredRun.Reasons...)
	}
}
