package harnessobs

import (
	"bufio"
	"crypto/sha256"
	"hash"
	"io"
)

// Event scan loop owns scanner setup and source-hash accumulation for replayed
// JSONL sources.
// The helper split below keeps the loop, scanner limits, and scan finalization
// in one navigable responsibility group without hiding source-hash behavior.
// Blank lines still contribute to the source hash because the hash represents
// the replayed input bytes, not only accepted event records.
// Event append handling stays in this file so the scan loop and accepted-event
// semantics can be reviewed together instead of through a metric-only helper.
// Scanner finalization is kept here for the same reason: it owns the hash
// returned with the accepted event list.
func scanEvents(profile Profile, file io.Reader, maxLine, maxEvents int) ([]Event, string, error) {
	scanner := eventScanner(file, maxLine)
	events := []Event{}
	sourceHash := sha256.New()

	lineNo := 0
	return scanEventLines(profile, scanner, sourceHash, events, lineNo, maxEvents)
}

func scanEventLines(profile Profile, scanner *bufio.Scanner, sourceHash hash.Hash, events []Event, lineNo, maxEvents int) ([]Event, string, error) {
	for scanner.Scan() {
		lineNo++
		if err := scanEventLine(profile, scanner.Bytes(), sourceHash, lineNo, maxEvents, &events); err != nil {
			return nil, "", err
		}
	}
	return scannedEvents(events, sourceHash, scanner.Err())
}

func scanEventLine(profile Profile, line []byte, sourceHash hash.Hash, lineNo, maxEvents int, events *[]Event) error {
	event, ok, err := readEventLine(profile, line, lineNo, len(*events), maxEvents, sourceHash)
	if err != nil {
		return err
	}
	if ok {
		*events = append(*events, event)
	}
	return nil
}
