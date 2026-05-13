package harnessobs

import (
	"strings"
)

func findStringByKeyIn(value any, wanted map[string]bool) string {
	// findStringByKeyIn keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	matchingString := func(value any) (string, bool) {
		s, ok := value.(string)
		return s, ok && strings.TrimSpace(s) != ""
	}

	s, ok := findByKeyIn(value, wanted, matchingString)
	if !ok {
		return ""
	}
	return s
}
