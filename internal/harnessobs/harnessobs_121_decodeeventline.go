package harnessobs

import (
	"encoding/json"

	"fmt"
)

func decodeEventLine(line []byte, lineNo int) (Event, error) {
	// decodeEventLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var event Event
	if err := json.Unmarshal(line, &event); err != nil {

		return Event{}, fmt.Errorf("source line %d: malformed_event", lineNo)
	}
	return event, nil
}
