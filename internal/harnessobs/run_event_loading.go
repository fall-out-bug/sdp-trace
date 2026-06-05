package harnessobs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Run event loading resolves only event references that pass the package's
// safe-reference rules before decoding the event payload.
// This preserves the run manifest's event order and fails on the first
// unreadable or unsafe reference.
func loadRunEvents(dir string, refs []string) ([]Event, error) {
	events := make([]Event, 0, len(refs))
	for _, ref := range refs {
		event, err := loadRunEvent(dir, ref)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func loadRunEvent(dir, ref string) (Event, error) {
	// Event references are manifest-controlled inputs and must not escape the
	// run directory before the file is read.
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
