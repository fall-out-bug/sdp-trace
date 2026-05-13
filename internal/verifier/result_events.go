package verifier

import "github.com/fall_out_bug/sdp-trace/internal/trace"

func observedEventSet(events []trace.Event) map[string]bool {
	observedEvents := map[string]bool{}
	for _, event := range events {
		// Contract coverage is event-type based; individual event payloads remain
		// chain evidence.
		observedEvents[string(event.EventType)] = true
	}
	return observedEvents
}
