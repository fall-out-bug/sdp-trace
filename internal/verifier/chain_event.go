package verifier

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func verifyChainEvent(events []trace.Event, index int, event trace.Event) string {
	// Checks are ordered from cheapest structural guard to strongest hash proof.
	// firstChainIssue short-circuits so later evidence cannot mask the earliest
	// chain defect.
	checks := []chainCheck{
		func() string { return verifyChainSequence(index, event) },
		func() string { return verifyChainPrevHash(events, index, event) },
		func() string { return verifyChainPayloadDigest(event) },
		func() string { return verifyChainEventHash(event) },
	}
	return firstChainIssue(checks)
}

type chainCheck func() string

func firstChainIssue(checks []chainCheck) string {
	for _, check := range checks {
		if issue := check(); issue != "" {
			// Return first defect to keep diagnostics deterministic.
			return issue
		}
	}
	return ""
}

func verifyChainSequence(index int, event trace.Event) string {
	if event.Sequence != index {
		// Event sequence is zero-based and must match file replay order.
		return fmt.Sprintf("sequence mismatch at %s", event.EventID)
	}
	return ""
}
