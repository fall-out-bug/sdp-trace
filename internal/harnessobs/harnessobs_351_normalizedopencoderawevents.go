package harnessobs

import (
	"os"

	"time"
)

func normalizedOpenCodeRawEvents(rawPath string, sessionFacts []Event, now time.Time) ([]Event, error) {
	// normalizedOpenCodeRawEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	file, err := os.Open(rawPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return scanOpenCodeRawEvents(file, sessionFacts, now)
}
