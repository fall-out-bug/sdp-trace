package harnessobs

import (
	"errors"
)

func validateEventRefs(event Event) error {
	// validateEventRefs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, check := range eventRefChecks(event) {

		if !check.ok {
			return errors.New(check.err)
		}
	}
	return nil
}
