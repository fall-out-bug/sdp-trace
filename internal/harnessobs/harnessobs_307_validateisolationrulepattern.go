package harnessobs

import (
	"errors"
)

func validateIsolationRulePattern(pattern string) error {
	// validateIsolationRulePattern keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if unsafeIsolationRulePattern(pattern) {

		return errors.New("unsafe isolation rule pattern")
	}
	return nil
}
