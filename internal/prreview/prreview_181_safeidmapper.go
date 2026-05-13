package prreview

import (
	"strings"
)

func safeIDMapper(r rune) rune {
	// safeIDMapper keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if r <= 127 && strings.IndexByte(safeIDAllowedChars, byte(r)) >= 0 {
		return r
	}

	return '-'
}
