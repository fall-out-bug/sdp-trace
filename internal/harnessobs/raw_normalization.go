package harnessobs

import (
	"errors"
	"path/filepath"
	"time"
)

// Raw normalization turns provider-specific raw JSONL into portable event JSONL
// only after the profile format and source/output relationship are validated.
func normalizeRawEvents(format, rawPath, outPath string, sessionFacts []Event, now time.Time) error {
	if err := validateRawNormalizationInputs(format, rawPath, outPath); err != nil {
		return err
	}
	events, err := normalizedOpenCodeRawEvents(rawPath, sessionFacts, rawNormalizationTime(now))
	if err != nil {
		return err
	}
	return writeNormalizedEvents(outPath, events)
}

func validateRawNormalizationInputs(format, rawPath, outPath string) error {
	if format != OpenCodeJSONLRawFormat {
		return errors.New("unsupported raw_event_format")
	}

	// Raw capture and normalized output must be distinct replay artifacts; using
	// the same path would overwrite the source before later evidence can be
	// inspected or re-normalized.
	if filepath.Clean(rawPath) == filepath.Clean(outPath) {
		return errors.New("raw_event_source_path and event_source_path must be different files")
	}
	return nil
}

func rawNormalizationTime(now time.Time) time.Time {
	if now.IsZero() {
		// A zero collection time is allowed for direct helper use; normalize to
		// UTC so generated observed_at values remain portable.
		return time.Now().UTC()
	}
	return now
}
