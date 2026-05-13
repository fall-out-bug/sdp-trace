package harnessobs

import (
	"strings"
)

func safeEventRef(ref string) bool {
	// safeEventRef keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeEventRefPath(ref) {
		return false
	}
	if !strings.HasPrefix(ref, "events/") || !strings.HasSuffix(ref, ".json") {
		return false
	}

	id := strings.TrimSuffix(strings.TrimPrefix(ref, "events/"), ".json")
	return safeFileIDPattern.MatchString(id)
}
