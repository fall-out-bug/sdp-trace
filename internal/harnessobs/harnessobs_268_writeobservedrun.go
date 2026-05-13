package harnessobs

import (
	"path/filepath"
)

func writeObservedRun(observedDir string, events []Event, observed Run) error {
	// writeObservedRun keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	if err := writeObservedEvents(observedDir, events); err != nil {
		return err
	}
	return writeJSON(filepath.Join(observedDir, "run.json"), observed)
}
