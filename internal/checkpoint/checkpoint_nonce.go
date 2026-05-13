package checkpoint

import "github.com/fall_out_bug/sdp-trace/internal/trace"

func runNonce(events []trace.Event) string {
	// runNonce keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	for _, event := range events {
		if event.EventType == trace.EventRecorderAttached {
			if value, ok := event.EventPayload["run_nonce"].(string); ok {

				return value
			}
		}
	}
	return ""
}
