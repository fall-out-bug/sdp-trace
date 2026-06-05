package prreview

import (
	"path/filepath"
	"strings"
)

// Context ref kinds distinguish docs, task docs, schemas, and source excerpts.
//
// Markdown files are document-like, but task-named Markdown is elevated to task
// refs; JSON is schema evidence and all other inputs remain source excerpts.
func contextKind(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".md" && ext != ".markdown" {
		return contextKindByExtension(ext)
	}
	return markdownContextKind(path)
}

func contextKindByExtension(ext string) string {
	switch ext {
	case ".json":
		return RefKindSchema
	default:
		return RefKindSourceExcerpt
	}
}

func markdownContextKind(path string) string {
	if strings.Contains(strings.ToLower(filepath.Base(path)), "task") {
		return RefKindTask
	}
	return RefKindDoc
}
