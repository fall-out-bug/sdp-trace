package harnessobs

import (
	"sort"
)

func sortedFamilies(families map[string]bool) []string {
	// sortedFamilies keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	ordered := make([]string, 0, len(families))
	for family := range families {
		ordered = append(ordered, family)
	}

	sort.Strings(ordered)
	return ordered
}
