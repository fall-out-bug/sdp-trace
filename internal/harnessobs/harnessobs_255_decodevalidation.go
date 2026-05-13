package harnessobs

import (
	"encoding/json"

	"io"
)

func DecodeValidation(r io.Reader) (Validation, error) {
	// DecodeValidation keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var validation Validation
	if err := json.NewDecoder(r).Decode(&validation); err != nil {

		return Validation{}, err
	}
	return validation, nil
}
