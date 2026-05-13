package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func contractHasRequiredRun(contract trace.Contract, id string) bool {
	// contractHasRequiredRun keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, required := range contract.RequiredRuns {
		if required.ID == id {
			return true
		}
	}
	return false
}

func contractHasEvidence(contract trace.Contract, id string) bool {
	// contractHasEvidence keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, required := range contract.RequiredEvidence {
		if required.ID == id {
			return true
		}
	}
	return false
}

func worseGateState(current, next string) string {
	// worseGateState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	if gateSeverity(next) > gateSeverity(current) {

		return next
	}
	return current
}
func gateSeverity(state string) int {
	return gateSeverityByState[state]
}

func containsString(values []string, target string) bool {
	// containsString keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.
	for _, value := range values {
		if value == target {

			return true
		}
	}
	return false
}
