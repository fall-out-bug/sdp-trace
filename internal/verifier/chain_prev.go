package verifier

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func verifyChainPrevHash(events []trace.Event, index int, event trace.Event) string {
	if event.PrevEventHash != chainExpectedPrevHash(events, index) {
		if index == 0 {
			// The first event must link to the null sentinel only.
			return "first event has non-empty prev_event_hash"
		}
		return fmt.Sprintf("broken chain at %d (%s)", index+1, event.EventID)
	}
	return ""
}

func chainExpectedPrevHash(events []trace.Event, index int) string {
	if index == 0 {
		// Genesis events link to the canonical null hash sentinel.
		return trace.NullEventHash
	}
	return events[index-1].EventHash
}
