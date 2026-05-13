package harnessobs

import (
	"fmt"
)

func validateParsedEvent(profile Profile, event Event, line []byte, lineNo int) error {
	// validateParsedEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	expected := digestLine(line)
	if event.SourceDigest != expected {
		return fmt.Errorf("source line %d: source_digest_mismatch:%s", lineNo, safeEvent(event.EventID))
	}

	if err := validateEvent(profile, event); err != nil {
		return fmt.Errorf("source line %d: %w", lineNo, err)
	}
	return nil
}
