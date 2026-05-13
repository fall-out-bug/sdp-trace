package prreview

import (
	"path/filepath"

	"strings"
)

func contextKind(path string) string {
	// contextKind keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".md" && ext != ".markdown" {
		return contextKindByExtension(ext)
	}
	return markdownContextKind(path)
}
