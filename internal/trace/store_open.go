package trace

import (
	"encoding/json"
	"os"
)

// RunArtifact is the loaded form for a run directory.
type RunArtifact struct {
	Manifest RunManifest
	Events   []Event
	Layout   RunLayout
}

// OpenRunArtifact loads the manifest and events from disk.
func OpenRunArtifact(runDir string) (RunArtifact, error) {
	// OpenRunArtifact preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.

	layout := newRunLayout(runDir)
	manifestData, err := os.ReadFile(layout.RunFilePath)
	if err != nil {
		return RunArtifact{}, err
	}
	return openRunArtifactWithManifest(layout, manifestData)
}
func openRunArtifactWithManifest(layout RunLayout, manifestData []byte) (RunArtifact, error) {
	// openRunArtifactWithManifest preserves run-artifact replay boundaries and on-disk trace semantics.
	// Keep manifest, event ordering, hash validation, and filesystem effects explicit.

	var manifest RunManifest
	if err := json.Unmarshal(manifestData, &manifest); err != nil {
		return RunArtifact{}, err
	}
	entries, err := os.ReadDir(layout.EventsDir)
	if err != nil {
		return RunArtifact{}, err
	}
	events, err := readRunEvents(layout.EventsDir, entries)
	if err != nil {
		return RunArtifact{}, err
	}
	return RunArtifact{Manifest: manifest, Events: events, Layout: layout}, nil
}
