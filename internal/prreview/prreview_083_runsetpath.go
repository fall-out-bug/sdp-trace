package prreview

import (
	"os"

	"path/filepath"
)

func runSetPath(path string) string {
	// runSetPath keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {

		return filepath.Join(path, "results.json")
	}
	return path
}
