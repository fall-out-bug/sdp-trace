package demo

import (
	"fmt"
	"os"
	"path/filepath"
)

func ensureRunRootDir(root string) error {
	// ensureRunRootDir keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", root)
	}
	return nil
}

func isRunDirCandidate(entry os.DirEntry, path string) bool {
	return entry.IsDir() && hasRunManifest(path)
}
func hasRunManifest(path string) bool {
	_, err := os.Stat(filepath.Join(path, "run.json"))
	return err == nil
}
