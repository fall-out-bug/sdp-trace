package demo

import (
	"fmt"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func applyRunRows(result *GateResult, rows []RunRow) map[string]bool {
	// applyRunRows keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	observedEvidence := map[string]bool{}
	for _, row := range rows {
		applyRunRow(result, observedEvidence, row)
	}
	return observedEvidence
}
func applyRunRow(result *GateResult, observedEvidence map[string]bool, row RunRow) {
	// applyRunRow keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	result.OverrideRequests = append(result.OverrideRequests, row.OverrideRequests...)
	markObservedEvidence(observedEvidence, row)
	applyRowResult(result, row)
	applyRowClosure(result, row)
}

func markObservedEvidence(observedEvidence map[string]bool, row RunRow) {
	// markObservedEvidence keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if row.Kind != "" && row.Kind != "unmatched" && row.Result == trace.VerdictObserved {

		observedEvidence[row.Kind] = true
	}
}
func applyRowResult(result *GateResult, row RunRow) {
	// applyRowResult keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if row.Result != trace.VerdictObserved {

		result.LocalGate = GateFail
		result.Reasons = append(result.Reasons, fmt.Sprintf("%s result is %s, expected observed", row.Name, row.Result))
	}
}

func applyRowClosure(result *GateResult, row RunRow) {
	// applyRowClosure keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if row.ClosureState != trace.ClosureStateCompleted {

		result.LocalGate = GateFail
		result.Reasons = append(result.Reasons, fmt.Sprintf("%s closure_state is %s", row.Name, row.ClosureState))
	}
}
