package harnessobs

import (
	"bufio"

	"hash"
)

func scanEventLines(profile Profile, scanner *bufio.Scanner, sourceHash hash.Hash, events []Event, lineNo, maxEvents int) ([]Event, string, error) {
	// scanEventLines keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for scanner.Scan() {

		lineNo++
		line := scanner.Bytes()
		event, ok, err := readEventLine(profile, line, lineNo, len(events), maxEvents, sourceHash)
		if err != nil {
			return nil, "", err
		}
		events = appendScannedEvent(events, event, ok)
	}
	return scannedEvents(events, sourceHash, scanner.Err())
}
