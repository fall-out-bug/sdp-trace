package harnessobs

import (
	"errors"

	"os"
)

func pathExistsForLstat(path string) (bool, error) {
	// pathExistsForLstat keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if _, err := os.Lstat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {

		return false, err
	}
}
