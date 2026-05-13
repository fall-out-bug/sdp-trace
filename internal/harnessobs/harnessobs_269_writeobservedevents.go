package harnessobs

import (
	"path/filepath"
)

func writeObservedEvents(observedDir string, events []Event) error {
	// writeObservedEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, event := range events {

		if err := writeJSON(filepath.Join(observedDir, "events", event.EventID+".json"), event); err != nil {
			return err
		}
	}
	return nil
}
