package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func previewGateMode(contract trace.Contract) string {
	mode := demo.GateModeObservation
	for _, required := range contract.RequiredRuns {
		switch required.Profile {
		case demo.GateModeProtectedFuture:
			// Protected-future requirements dominate advisory CI requirements in
			// preview because they imply a stricter future gate path.
			return demo.GateModeProtectedFuture
		case demo.GateModeAdvisoryCI:
			// Advisory CI is retained unless a protected-future requirement is
			// found later.
			mode = demo.GateModeAdvisoryCI
		}
	}
	return mode
}
