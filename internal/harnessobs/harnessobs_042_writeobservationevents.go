package harnessobs

import (
	"path/filepath"
)

func writeObservationEvents(outDir string, events []Event) error {
	// writeObservationEvents keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, event := range events {

		path := filepath.Join(outDir, "events", event.EventID+".json")
		if err := writeJSON(path, event); err != nil {
			return err
		}
	}
	return nil
}
