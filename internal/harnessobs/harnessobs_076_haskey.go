package harnessobs

import (
	"strings"
)

func hasKey(value any, keys ...string) bool {
	// hasKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	wanted := map[string]bool{}
	for _, key := range keys {
		wanted[strings.ToLower(key)] = true
	}

	return hasKeyIn(value, wanted)
}
