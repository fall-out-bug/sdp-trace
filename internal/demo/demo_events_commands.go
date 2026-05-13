package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func commandEvents(events []trace.Event) (trace.Event, trace.Event) {
	// Last-seen command events win, matching replay behavior for append-only
	// traces that may contain retries.
	var started trace.Event
	var finished trace.Event
	for _, event := range events {
		switch event.EventType {
		case trace.EventCommandStarted:
			started = event
		case trace.EventCommandFinished:
			finished = event
		}
	}
	return started, finished
}

func classify(events []trace.Event, requirements []trace.EvidenceRequirement) (string, string) {
	// classify keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, requirement := range requirements {
		if evidenceRequirementMatches(events, requirement) {
			return requirement.ID, "matched contract evidence requirement"
		}
	}
	return "unmatched", "no contract evidence requirement matched"
}
