package harnessobs

import (
	"bufio"
	"encoding/hex"
	"hash"
	"io"
	"os"
)

// Event scan input setup owns the replay boundary for source JSONL bytes:
// profile policy is loaded first, source paths are already resolved by callers,
// and the returned digest is derived from the bytes actually scanned.
func readEventsFromPath(profilePath, sourcePath string) ([]Event, string, error) {
	profile, err := LoadProfile(profilePath)
	if err != nil {
		return nil, "", err
	}
	return readEvents(profile, sourcePath)
}

// readEvents keeps source opening inside the same boundary as scan-limit
// enforcement, so callers cannot observe a digest for bytes that were not
// accepted by the profile's replay limits.
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

// scannedEvents withholds the digest when scanning fails; a partial replay
// cannot become evidence for the source.
func scannedEvents(events []Event, sourceHash hash.Hash, scanErr error) ([]Event, string, error) {
	if scanErr != nil {
		return nil, "", scanErr
	}
	return events, hex.EncodeToString(sourceHash.Sum(nil)), nil
}
