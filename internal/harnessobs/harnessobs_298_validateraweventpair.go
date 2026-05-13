package harnessobs

import (
	"errors"
)

func validateRawEventPair(hasFormat, hasSource bool) error {
	// validateRawEventPair keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	switch {
	case hasFormat == hasSource:
		return nil
	case hasFormat:
		return errors.New("raw_event_source_path required for raw_event_format")
	default:
		return errors.New("raw_event_format required for raw_event_source_path")
	}
}
