package harnessobs

import (
	"bufio"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

// Event scan input setup owns scanner limits and source-hash finalization for
// replayed JSONL bytes.
func readEvents(profile Profile, sourcePath string) ([]Event, string, error) {
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	maxLine, maxEvents := effectiveEventLimits(profile.Limits)
	return scanEvents(profile, file, maxLine, maxEvents)
}

func eventScanner(file io.Reader, maxLine int) *bufio.Scanner {
	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLine)
	return scanner
}

func scannedEvents(events []Event, sourceHash hash.Hash, scanErr error) ([]Event, string, error) {
	if scanErr != nil {
		return nil, "", scanErr
	}
	return events, hex.EncodeToString(sourceHash.Sum(nil)), nil
}
