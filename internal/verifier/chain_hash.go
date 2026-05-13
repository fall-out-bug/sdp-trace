package verifier

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func verifyChainPayloadDigest(event trace.Event) string {
	if err := event.VerifyPayloadDigest(); err != nil {
		// Payload-digest failure blocks chain trust even if event_hash matches a
		// malformed payload copy.
		return fmt.Sprintf("invalid payload digest for %s: %s", event.EventID, err)
	}
	return ""
}

func verifyChainEventHash(event trace.Event) string {
	recomputed, err := trace.EventHash(event)
	if err != nil {
		// Hash recomputation can fail when canonical event shape is invalid.
		return fmt.Sprintf("invalid event hash for %s", event.EventID)
	}
	if event.EventHash != recomputed {
		// Event hash mismatch proves retained event contents no longer match the
		// recorded chain.
		return fmt.Sprintf("hash mismatch for %s", event.EventID)
	}
	return ""
}
