package harnessobs

import (
	"errors"

	"os"
)

func ensureExpectedPathType(path string, spec existingPathSpec) error {
	// ensureExpectedPathType keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() == spec.requireDir {
		return nil
	}

	return errors.New(spec.typeError)
}
