package demo

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func payloadString(event trace.Event, key string) string {
	// payloadString keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	value := payloadValue(event, key)
	if value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprint(typed)
	}
}
