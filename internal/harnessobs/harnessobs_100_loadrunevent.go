package harnessobs

import (
	"encoding/json"

	"fmt"

	"os"
	"path/filepath"
)

func loadRunEvent(dir, ref string) (Event, error) {
	// loadRunEvent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !safeEventRef(ref) {

		return Event{}, fmt.Errorf("unsafe event ref: %s", ref)
	}

	data, err := os.ReadFile(filepath.Join(dir, ref))
	if err != nil {
		return Event{}, err
	}
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return Event{}, err
	}

	return event, nil
}
