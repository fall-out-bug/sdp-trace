package prreview

import (
	"path/filepath"
	"strings"
)

// Copied input file types and extensions are normalized for portable refs.
//
// Only a narrow extension allow-list is retained in packet input names; unknown
// source extensions become text so downstream tools do not infer unsupported
// semantics from local filenames.
func contentType(path string) string {
	if kind, ok := contentTypeByExt[normalizedContentExt(path)]; ok {
		return kind
	}
	return ContentText
}

var contentTypeByExt = map[string]string{
	".json": ContentJSON,
	".md":   ContentMarkdown,
}

func normalizedContentExt(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if ext == ".markdown" {
		return ".md"
	}
	return ext
}

func normalizedExt(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".json", ".md", ".txt", ".diff", ".patch":
		return ext
	default:
		return ".txt"
	}
}
