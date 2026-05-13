package harnessobs

import (
	"path/filepath"
)

func cleanProfileRelativePath(profilePath, relPath string) string {
	// cleanProfileRelativePath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	baseDir := filepath.Dir(profilePath)
	if baseDir == "." {
		baseDir = ""
	}
	return filepath.Clean(filepath.Join(baseDir, relPath))
}
