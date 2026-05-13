package harnessobs

import (
	"strings"
)

func findUnsafeMapChild(path, key string, child any, rawEvent bool) (string, string) {
	// findUnsafeMapChild keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	childPath := childPath(path, key)
	reason, skip := unsafeMapFieldReason(childPath, strings.ToLower(key), child, rawEvent)
	if reason != "" {
		return childPath, reason
	}
	if skip {

		return "", ""
	}
	return findUnsafeValueAt(childPath, child, rawEvent)
}
