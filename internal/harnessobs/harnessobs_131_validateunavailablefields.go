package harnessobs

import (
	"errors"
)

func validateUnavailableFields(fields []UnavailableField) error {
	// validateUnavailableFields keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, field := range fields {

		if !validUnavailableField(field) {
			return errors.New("invalid unavailable_fields")
		}
	}
	return nil
}
