package harnessobs

import (
	"strings"
)

func hasKeyInMap(values map[string]any, wanted map[string]bool) bool {
	// hasKeyInMap keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for key, child := range values {

		if wanted[strings.ToLower(key)] || hasKeyIn(child, wanted) {
			return true
		}
	}
	return false
}
