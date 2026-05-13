package harnessobs

import (
	"os"
)

func readEvents(profile Profile, sourcePath string) ([]Event, string, error) {
	// readEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	file, err := os.Open(sourcePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	maxLine, maxEvents := effectiveEventLimits(profile.Limits)
	return scanEvents(profile, file, maxLine, maxEvents)
}
