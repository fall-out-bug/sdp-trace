package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

var previewBoundaryRows = []previewBoundary{
	{
		Boundary: string(trace.ObservationBoundaryProcessWrapper),
		State:    "pass",
		Reason:   "preview covers local process-wrapper capture only",
	},
	{
		Boundary: string(trace.ObservationBoundaryAdapterSocket),
		State:    string(trace.ObservationStateNotIntegrated),
		Reason:   "adapter socket/API capture is not configured in Block 13B",
	},
	{
		Boundary: string(trace.ObservationBoundaryToolWrapper),
		State:    string(trace.ObservationStateUnsupported),
		Reason:   "tool-level wrapping is a future observation boundary",
	},
	{
		Boundary: string(trace.ObservationBoundaryVCSPRObserver),
		State:    string(trace.ObservationStateNotIntegrated),
		Reason:   "VCS/PR observer is not configured in Block 13B",
	},
	{
		Boundary: string(trace.ObservationBoundaryCIObserver),
		State:    string(trace.ObservationStateOfflineDev),
		Reason:   "CI witness cannot be produced by local preview",
	},
	{
		// External witness state remains not_integrated until a concrete witness
		// profile exists.
		Boundary: string(trace.ObservationBoundaryExternalWitness),
		State:    string(trace.ObservationStateNotIntegrated),
		Reason:   "external witness profile is not implemented in Block 13B",
	},
}

func previewBoundaries() []previewBoundary {
	// Preview is explicit about which observation boundaries are local,
	// unsupported, or not integrated.
	return append([]previewBoundary(nil), previewBoundaryRows...)
}
