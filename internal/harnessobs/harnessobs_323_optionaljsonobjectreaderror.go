package harnessobs

import (
	"errors"

	"os"
)

func optionalJSONObjectReadError(err error) (map[string]any, error) {
	// optionalJSONObjectReadError keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if errors.Is(err, os.ErrNotExist) {

		return map[string]any{}, nil
	}
	return nil, err
}
