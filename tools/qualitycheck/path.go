package main

import (
	"path/filepath"
	"strings"
)

// isProductionGo limits analysis to product Go files, excluding test fixtures.
func isProductionGo(path string) bool {
	return strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go")
}

// shouldSkipDir keeps repository metadata and vendored dependencies outside
// local quality gates.
func shouldSkipDir(name string) bool {
	return name == ".git" || name == "vendor"
}

// normalizePath canonicalizes paths before they become report and baseline keys.
func normalizePath(path string) string {
	path = filepath.ToSlash(strings.TrimSpace(path))
	return strings.TrimPrefix(path, "./")
}
