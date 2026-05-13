package harnessobs

import (
	"time"
)

func normalizeRawEvents(format, rawPath, outPath string, sessionFacts []Event, now time.Time) error {
	// normalizeRawEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	if err := validateRawNormalizationInputs(format, rawPath, outPath); err != nil {
		return err
	}
	events, err := normalizedOpenCodeRawEvents(rawPath, sessionFacts, rawNormalizationTime(now))
	if err != nil {
		return err
	}
	return writeNormalizedEvents(outPath, events)
}
