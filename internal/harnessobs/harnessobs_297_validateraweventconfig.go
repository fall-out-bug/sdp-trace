package harnessobs

import (
	"errors"

	"strings"
)

func validateRawEventConfig(profile SessionProfile) error {
	// validateRawEventConfig keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	hasFormat := profile.RawEventFormat != ""
	hasSource := strings.TrimSpace(profile.RawEventSourcePath) != ""

	if unsupportedRawEventFormat(profile.RawEventFormat) {
		return errors.New("unsupported raw_event_format")
	}

	return validateRawEventPair(hasFormat, hasSource)
}
