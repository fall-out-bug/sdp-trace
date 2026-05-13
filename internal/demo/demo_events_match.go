package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func evidenceRequirementMatches(events []trace.Event, requirement trace.EvidenceRequirement) bool {
	// evidenceRequirementMatches keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if requirement.ID == "" {
		return false
	}
	for _, event := range events {
		if eventMatchesRequirement(event, requirement) {
			return true
		}
	}
	return false
}

func eventMatchesRequirement(event trace.Event, requirement trace.EvidenceRequirement) bool {
	// eventMatchesRequirement keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	if requirement.EventType != "" && event.EventType != trace.EventType(requirement.EventType) {
		return false
	}
	return payloadString(event, requirement.PayloadField) == requirement.PayloadEquals
}
