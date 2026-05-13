package harnessobs

import (
	"errors"

	"time"
)

func validateObservedAt(value string) error {
	// validateObservedAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if value == "" {

		return nil
	}
	if _, err := time.Parse(time.RFC3339, value); err != nil {
		return errors.New("invalid observed_at")
	}
	return nil
}
