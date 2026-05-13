package harnessobs

import (
	"strings"
)

func findUnsafeStringAt(path, value string, rawEvent bool) (string, string) {
	// findUnsafeStringAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(value) == "" {

		return "", ""
	}
	if reason := unsafeStringReason(path, value, rawEvent); reason != "" {
		return path, reason
	}
	return "", ""
}
