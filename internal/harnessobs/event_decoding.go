package harnessobs

import (
	"encoding/json"
	"fmt"
)

// Event decoding first inspects raw fields for unsafe input before decoding
// into the typed event contract used by downstream validation.
// Raw and typed decode errors intentionally keep separate labels so malformed
// JSONL and malformed event envelopes remain distinguishable in replay output.
// Unsafe-field rejection stays between those passes so forbidden raw fields are
// caught before Go struct decoding can ignore unknown data.
// Keep these labels stable; tests and user-facing diagnostics depend on them.
func decodeSafeEventLine(line []byte, lineNo int) (Event, error) {
	raw, err := decodeRawEventLine(line, lineNo)
	if err != nil {
		return Event{}, err
	}
	if err := rejectUnsafeEvent(raw, lineNo); err != nil {
		return Event{}, err
	}
	return decodeEventLine(line, lineNo)
}

func decodeRawEventLine(line []byte, lineNo int) (map[string]any, error) {
	raw := map[string]any{}
	if err := decodeJSONLine(line, &raw); err != nil {
		return nil, fmt.Errorf("source line %d: malformed_jsonl", lineNo)
	}
	return raw, nil
}

func decodeJSONLine(line []byte, target any) error {
	return json.Unmarshal(line, target)
}

func decodeEventLine(line []byte, lineNo int) (Event, error) {
	var event Event
	if err := decodeJSONLine(line, &event); err != nil {
		return Event{}, fmt.Errorf("source line %d: malformed_event", lineNo)
	}
	return event, nil
}

func rejectUnsafeEvent(raw map[string]any, lineNo int) error {
	if unsafeField, reason := findUnsafe(raw); unsafeField != "" {
		return fmt.Errorf("source line %d: unsafe_input:%s:%s", lineNo, unsafeField, reason)
	}
	return nil
}
