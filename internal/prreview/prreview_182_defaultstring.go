package prreview

import (
	"strings"
)

func defaultString(value, fallback string) string {
	// defaultString keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if strings.TrimSpace(value) == "" {

		return fallback
	}
	return value
}
