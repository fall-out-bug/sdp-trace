package harnessobs

import (
	"errors"
)

func ensureOutParentInsideWorkingDirectory(clean string) error {
	// ensureOutParentInsideWorkingDirectory keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parentEscapes, err := outParentEscapes(clean)
	if err != nil {
		return err
	}
	if parentEscapes {

		return errors.New("out parent path escapes working directory")
	}
	return nil
}
