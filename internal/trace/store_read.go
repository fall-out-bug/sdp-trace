package trace

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// Run event reads are intentionally structural only.
// The loader validates event shape and deterministic ordering, while chain
// trust is replayed later by ValidateEventChain when callers request it.
// This separation lets callers inspect malformed or partial run directories
// without silently upgrading them into verified evidence.
// File filtering also keeps non-event artifacts out of the replay surface.

// readRunEvents loads and sorts every *.json file in events/.
func readRunEvents(eventsDir string, entries []fs.DirEntry) ([]Event, error) {
	// Stable file ordering makes replay independent of filesystem iteration.
	jsonFiles := eventJSONFiles(eventsDir, entries)
	events := make([]Event, 0, len(jsonFiles))
	for _, path := range jsonFiles {
		// Each event is shape-validated before it joins the replay slice.
		event, err := readRunEvent(path)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func eventJSONFiles(eventsDir string, entries []fs.DirEntry) []string {
	// Only event JSON files participate in local chain replay.
	jsonFiles := make([]string, 0, len(entries))
	for _, entry := range entries {
		if isEventJSON(entry) {
			// Join only after filtering so subdirectories and sidecar files do
			// not become replay candidates.
			jsonFiles = append(jsonFiles, filepath.Join(eventsDir, entry.Name()))
		}
	}
	// Lexical sort preserves the sequence-prefixed event naming contract.
	sort.Strings(jsonFiles)
	return jsonFiles
}

func isEventJSON(entry fs.DirEntry) bool {
	// Event replay accepts flat JSON files only; nested directories are outside
	// the portable trace layout.
	return !entry.IsDir() && filepath.Ext(entry.Name()) == ".json"
}

func readRunEvent(path string) (Event, error) {
	// Event bytes must decode and validate before chain placement is trusted.
	payload, err := os.ReadFile(path)
	if err != nil {
		// Missing event bytes are a structural replay gap for the run.
		return Event{}, err
	}
	var event Event
	if err := json.Unmarshal(payload, &event); err != nil {
		// Malformed JSON cannot enter hash or position validation.
		return Event{}, err
	}
	if err := event.Validate(); err != nil {
		// Shape validation runs before sequence and hash-chain checks.
		return Event{}, fmt.Errorf("invalid event %q: %w", path, err)
	}
	return event, nil
}
