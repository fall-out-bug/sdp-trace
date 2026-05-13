package harnessobs

import (
	"encoding/hex"

	"hash"
)

func scannedEvents(events []Event, sourceHash hash.Hash, scanErr error) ([]Event, string, error) {
	// scannedEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if scanErr != nil {
		return nil, "", scanErr
	}

	return events, hex.EncodeToString(sourceHash.Sum(nil)), nil
}
