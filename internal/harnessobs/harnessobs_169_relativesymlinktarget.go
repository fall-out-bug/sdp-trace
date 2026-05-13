package harnessobs

import (
	"os"
	"path/filepath"
)

func relativeSymlinkTarget(path string) (string, error) {
	// relativeSymlinkTarget keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(resolved)
	if err != nil {
		return "", err
	}
	return filepath.Rel(cwd, abs)
}
