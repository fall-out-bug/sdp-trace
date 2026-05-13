package demo

import (
	"os"
	"path/filepath"
)

func collectRunDirs(root string) ([]string, error) {
	// collectRunDirs keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	runDirs := make([]string, 0)
	for _, entry := range entries {

		path := filepath.Join(root, entry.Name())
		if isRunDirCandidate(entry, path) {
			runDirs = append(runDirs, path)
		}
	}
	return runDirs, nil
}
