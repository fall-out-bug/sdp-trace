package harnessobs

import (
	"path/filepath"
)

func resolveExistingAbsolutePath(path string) (string, error) {
	// resolveExistingAbsolutePath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	return absolutePath(resolved)
}
