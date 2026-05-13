package harnessobs

import (
	"errors"
)

func validateEventIdentity(profile Profile, event Event) error {
	// validateEventIdentity keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, check := range eventIdentityChecks(profile, event) {

		if !check.ok {
			return errors.New(check.err)
		}
	}
	return validateObservedAt(event.ObservedAt)
}
