package harnessobs

import (
	"time"
)

func appendNormalizedRawLine(events []Event, line []byte, lineNo int, now time.Time) ([]Event, error) {
	// appendNormalizedRawLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	lineEvents, err := normalizeOpenCodeRawLineBytes(line, lineNo, now)
	if err != nil {
		return nil, err
	}

	return append(events, lineEvents...), nil
}
