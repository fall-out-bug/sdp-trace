package harnessobs

import (
	"path/filepath"
)

func validateIsolationParent(clean string) error {
	// validateIsolationParent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parent := filepath.Dir(clean)
	if err := validatePotentialParentPath(parent); err != nil {
		return err
	}

	return ensureOutParentInsideWorkingDirectory(parent)
}
