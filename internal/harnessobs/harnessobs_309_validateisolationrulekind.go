package harnessobs

import (
	"errors"
)

func validateIsolationRuleKind(kind string) error {
	// validateIsolationRuleKind keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	switch kind {
	case "ignore_line", "json_read_deny":
		return nil
	default:
		return errors.New("unsupported isolation rule kind")
	}
}
