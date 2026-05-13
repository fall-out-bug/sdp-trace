package demo

import (
	"fmt"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func overrideRequestRequiredRunState(event trace.Event, contract trace.Contract, state string, reason string) (string, string) {
	// overrideRequestRequiredRunState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, id := range payloadStringSlice(event, "affected_required_runs") {
		if !contractHasRequiredRun(contract, id) {
			state = GateCannotVerify
			reason = fmt.Sprintf("override request references unknown required run %s", id)
		}
	}
	return state, reason
}
