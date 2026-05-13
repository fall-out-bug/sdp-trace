package harnessobs

import (
	"errors"

	"path/filepath"
)

func safeOutFile(path string) (string, error) {
	// safeOutFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeOutputPath(path) {
		return "", errors.New("out must be a relative local file without traversal")
	}

	parent, err := safeParentDir(filepath.Dir(path))
	if err != nil {
		return "", err
	}
	base := filepath.Base(path)
	if !safeOutputBaseName(base) {
		return "", errors.New("unsafe output filename")
	}

	return filepath.Join(parent, base), nil
}
