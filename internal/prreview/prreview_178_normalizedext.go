package prreview

import (
	"path/filepath"

	"strings"
)

func normalizedExt(path string) string {
	// normalizedExt keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json", ".md", ".txt", ".diff", ".patch":

		return ext
	default:
		return ".txt"
	}
}
