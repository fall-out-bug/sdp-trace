package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func gateMode(contract trace.Contract) string {
	// gateMode keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	mode := GateModeObservation
	for _, required := range contract.RequiredRuns {
		switch required.Profile {
		case GateModeProtectedFuture:
			return GateModeProtectedFuture
		case GateModeAdvisoryCI:
			mode = GateModeAdvisoryCI
		}
	}
	return mode
}
