package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func overrideRequestReferenceState(event trace.Event, contract trace.Contract, state string, reason string) (string, string) {
	// overrideRequestReferenceState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	state, reason = overrideRequestRequiredRunState(event, contract, state, reason)
	return overrideRequestEvidenceState(event, contract, state, reason)
}
