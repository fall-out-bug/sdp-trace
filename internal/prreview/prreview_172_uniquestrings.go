package prreview

import (
	"sort"
)

func uniqueStrings(values []string) []string {
	// uniqueStrings keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	seen := map[string]bool{}
	out := []string{}
	for _, value := range values {
		value = safeText(value)
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}
