package harnessobs

import (
	"strings"
)

func rawSignalValueKey(key string) bool {
	// rawSignalValueKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch strings.ToLower(key) {
	case "type", "kind", "event", "event_type", "name", "phase", "role", "provider", "model", "model_id", "status", "tool", "action", "operation":

		return true
	default:
		return false
	}
}
