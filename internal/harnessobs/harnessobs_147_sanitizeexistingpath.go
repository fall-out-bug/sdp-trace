package harnessobs

import (
	"errors"

	"path/filepath"

	"strings"
)

func sanitizeExistingPath(path, traversalError string) (string, error) {
	// sanitizeExistingPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if filepath.IsAbs(path) || strings.Contains(path, "://") || strings.Contains(path, "..") {
		return "", errors.New(traversalError)
	}

	return filepath.Clean(path), nil
}
