package harnessobs

import (
	"encoding/json"

	"fmt"

	"time"
)

func normalizeOpenCodeRawLineBytes(line []byte, lineNo int, now time.Time) ([]Event, error) {
	// normalizeOpenCodeRawLineBytes keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if blankJSONLLine(line) {
		return nil, nil
	}
	// Decode as a generic map first so unsafe provider fields can be rejected
	// before typed event construction drops unknown data.
	var raw map[string]any
	if err := json.Unmarshal(line, &raw); err != nil {
		return nil, fmt.Errorf("raw source line %d: malformed_jsonl", lineNo)
	}

	if err := rejectUnsafeRawEvent(raw, lineNo); err != nil {
		return nil, err
	}
	events := normalizeOpenCodeRawLine(raw, lineNo, now)
	return addNormalizedSourceDigests(events)
}
