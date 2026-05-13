package prreview

import (
	"strings"
)

func safeID(value string) string {
	// safeID keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	value = strings.ToLower(strings.TrimSpace(value))
	out := strings.Trim(strings.Map(safeIDMapper, value), "-.")
	if out == "" {

		return "item"
	}
	return out
}
