package harnessobs

import (
	"errors"

	"os"
)

func ignorableMissingOutDir(err error) error {
	// ignorableMissingOutDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err == nil || errors.Is(err, os.ErrNotExist) {

		return nil
	}
	return err
}
