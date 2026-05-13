package harnessobs

import (
	"errors"

	"path/filepath"
)

func safeOutDir(path string) (string, error) {
	// safeOutDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeOutputPath(path) {
		return "", errors.New("out must be a relative local directory without traversal")
	}

	clean := filepath.Clean(path)
	return safeCleanOutDir(clean)
}
