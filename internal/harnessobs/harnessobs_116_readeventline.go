package harnessobs

import (
	"fmt"

	"io"
)

func readEventLine(profile Profile, line []byte, lineNo, eventCount, maxEvents int, sourceHash io.Writer) (Event, bool, error) {
	// readEventLine keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if _, err := sourceHash.Write(line); err != nil {
		return Event{}, false, err
	}
	if blankJSONLLine(line) {

		return Event{}, false, nil
	}
	if eventCount >= maxEvents {

		return Event{}, false, fmt.Errorf("source line %d: event limit exceeded", lineNo)
	}
	event, err := parseEvent(profile, line, lineNo)
	return event, err == nil, err
}
