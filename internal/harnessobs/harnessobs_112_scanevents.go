package harnessobs

import (
	"bufio"
	"crypto/sha256"

	"io"
)

func scanEvents(profile Profile, file io.Reader, maxLine, maxEvents int) ([]Event, string, error) {
	// scanEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLine)
	events := []Event{}
	sourceHash := sha256.New()

	lineNo := 0
	return scanEventLines(profile, scanner, sourceHash, events, lineNo, maxEvents)
}
