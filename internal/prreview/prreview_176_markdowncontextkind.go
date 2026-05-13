package prreview

import (
	"path/filepath"

	"strings"
)

func markdownContextKind(path string) string {
	// markdownContextKind keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if strings.Contains(strings.ToLower(filepath.Base(path)), "task") {

		return RefKindTask
	}
	return RefKindDoc
}
