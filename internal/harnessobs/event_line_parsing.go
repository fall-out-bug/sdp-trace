package harnessobs

import (
	"fmt"
	"io"
	"strings"
)

// Event line parsing applies source hashing and event-count limits before
// delegating JSON decoding and semantic validation.
func readEventLine(profile Profile, line []byte, lineNo, eventCount, maxEvents int, sourceHash io.Writer) (Event, bool, error) {
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

func parseEvent(profile Profile, line []byte, lineNo int) (Event, error) {
	event, err := decodeSafeEventLine(line, lineNo)
	if err != nil {
		return Event{}, err
	}
	return event, validateParsedEvent(profile, event, line, lineNo)
}

// blankJSONLLine is shared by normal event replay and raw normalization so
// both source types treat whitespace-only records as non-events.
func blankJSONLLine(line []byte) bool {
	return len(strings.TrimSpace(string(line))) == 0
}
