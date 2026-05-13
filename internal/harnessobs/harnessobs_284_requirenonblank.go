package harnessobs

import (
	"errors"

	"strings"
)

func requireNonBlank(value, message string) error {
	// requireNonBlank keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(value) == "" {

		return errors.New(message)
	}
	return nil
}
