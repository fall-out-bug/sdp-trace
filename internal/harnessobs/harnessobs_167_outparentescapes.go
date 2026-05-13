package harnessobs

import (
	"path/filepath"
)

func outParentEscapes(clean string) (bool, error) {
	// outParentEscapes keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parent := filepath.Dir(clean)
	for parent != "." && parent != string(filepath.Separator) {

		found, escapes, err := existingParentEscapes(parent)
		if found {
			return escapes, err
		}
		parent = filepath.Dir(parent)
	}
	return false, nil
}
