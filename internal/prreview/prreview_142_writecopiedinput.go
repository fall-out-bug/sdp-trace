package prreview

import (
	"os"

	"path/filepath"
)

func writeCopiedInput(inputDir, name string, data []byte) error {
	// writeCopiedInput keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	dest := filepath.Join(inputDir, name)
	if err := os.WriteFile(dest, data, 0o644); err != nil {

		return err
	}
	return nil
}
