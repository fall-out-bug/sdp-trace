package harnessobs

import (
	"errors"
)

func validatePotentialParentPath(path string) error {
	// validatePotentialParentPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafePotentialParentPath(path) {

		return errors.New("parent path must be relative local directory without traversal")
	}
	return nil
}
