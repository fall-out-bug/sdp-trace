package harnessobs

import (
	"os"
	"path/filepath"
)

func createNormalizedEventsFile(outPath string) (*os.File, error) {
	// createNormalizedEventsFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return nil, err
	}
	return os.Create(outPath)
}
