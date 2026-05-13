package harnessobs

import (
	"os"
	"path/filepath"
)

func relativePathFromWorkingDirectory(path string) (string, error) {
	// relativePathFromWorkingDirectory keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Rel(cwd, path)
}
