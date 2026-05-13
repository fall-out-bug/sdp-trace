package prreview

import (
	"os"

	"path/filepath"
)

func prepareRunDirectories(outDir string) (string, error) {
	// prepareRunDirectories keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if err := ensureNewDir(outDir); err != nil {
		return "", err
	}
	rawDir := filepath.Join(outDir, "raw")

	if err := os.MkdirAll(rawDir, 0o755); err != nil {
		return "", err
	}
	return rawDir, nil
}
