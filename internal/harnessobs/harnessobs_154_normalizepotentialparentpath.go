package harnessobs

import (
	"path/filepath"
)

func normalizePotentialParentPath(path string) (string, error) {
	// normalizePotentialParentPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if path == "" {

		path = "."
	}
	if err := validatePotentialParentPath(path); err != nil {
		return "", err
	}
	return filepath.Clean(path), nil
}
