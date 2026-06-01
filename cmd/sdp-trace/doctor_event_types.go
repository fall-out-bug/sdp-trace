package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func knownEventType(eventType string) bool {
	switch trace.EventType(eventType) {
	case trace.EventRecorderAttached,
		trace.EventRunStarted,
		trace.EventCommandStarted,
		trace.EventCommandFinished,
		trace.EventRunClosed,
		trace.EventPolicyOverrideRequested:
		// Keep doctor scoped to the stable local recorder event model.
		// Unsupported future events remain explicit spec-drift gaps.
		return true
	default:
		return false
	}
}
