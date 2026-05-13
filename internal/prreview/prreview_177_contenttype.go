package prreview

import (
	"path/filepath"

	"strings"
)

func contentType(path string) string {
	// contentType keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":

		return ContentJSON
	case ".md", ".markdown":
		return ContentMarkdown
	default:
		return ContentText
	}
}
