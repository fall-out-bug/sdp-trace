package harnessobs

import (
	"encoding/json"
	"fmt"
	"time"
)

func appendNormalizedRawLine(events []Event, line []byte, lineNo int, now time.Time) ([]Event, error) {
	// Empty raw lines contribute no events, while malformed or unsafe lines stop
	// normalization before any partial output file is written.
	lineEvents, err := normalizeOpenCodeRawLineBytes(line, lineNo, now)
	if err != nil {
		return nil, err
	}

	return append(events, lineEvents...), nil
}

func normalizeOpenCodeRawLineBytes(line []byte, lineNo int, now time.Time) ([]Event, error) {
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

func rejectUnsafeRawEvent(raw map[string]any, lineNo int) error {
	if unsafeField, reason := findUnsafeRawEvent(raw); unsafeField != "" {
		// Include the raw line number because provider JSONL is the evidence the
		// user must inspect when normalization refuses to continue.
		return fmt.Errorf("raw source line %d: unsafe_input:%s:%s", lineNo, unsafeField, reason)
	}
	return nil
}
