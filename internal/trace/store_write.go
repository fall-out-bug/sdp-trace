package trace

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// This file owns persisted run and event JSON writes.

func (layout RunLayout) WriteEvent(event Event) error {
	// WriteEvent preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.
	event = event.EnsureDefaults()
	filename := EventFileName(event.Sequence, event.EventType)
	path := filepath.Join(layout.EventsDir, filename)
	encoded, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, encoded, 0o644)
}

// WriteRun writes run manifest under run.json.
func (layout RunLayout) WriteRun(manifest RunManifest) error {
	// WriteRun preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.
	encoded, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(layout.RunFilePath, encoded, 0o644)
}
