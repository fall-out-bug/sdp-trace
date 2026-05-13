package harnessobs

import (
	"bufio"

	"io"

	"time"
)

func scanOpenCodeRawEvents(file io.Reader, sessionFacts []Event, now time.Time) ([]Event, error) {
	// scanOpenCodeRawEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), DefaultMaxLineBytes)
	lineNo := 0

	events := append([]Event{}, sessionFacts...)
	for scanner.Scan() {
		lineNo++
		var err error

		events, err = appendNormalizedRawLine(events, scanner.Bytes(), lineNo, now)
		if err != nil {
			return nil, err
		}
	}

	return events, scanner.Err()
}
