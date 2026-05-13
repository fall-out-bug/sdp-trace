package verifier

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func loadRunEvents(runDir string) ([]trace.Event, error) {
	eventDir := filepath.Join(runDir, "events")
	eventFiles, err := loadRunEventFiles(eventDir)
	if err != nil {
		// The public error names the events directory as the missing replay
		// component.
		return nil, fmt.Errorf("events directory missing: %w", err)
	}
	sort.Strings(eventFiles)
	return parseRunEventFiles(eventFiles)
}

func parseRunEventFiles(eventFiles []string) ([]trace.Event, error) {
	if err := requireRunEventFiles(eventFiles); err != nil {
		return nil, err
	}
	// Parsing preserves one event per accepted file. Structural trust checks
	// stay in chain verification so parse errors and replay contradictions have
	// distinct verifier outcomes.
	events := make([]trace.Event, 0, len(eventFiles))
	for _, path := range eventFiles {
		var err error
		events, err = appendParsedRunEvent(events, path)
		if err != nil {
			return nil, err
		}
	}
	return events, nil
}

func requireRunEventFiles(eventFiles []string) error {
	if len(eventFiles) == 0 {
		// No valid event files means no replayable telemetry.
		return errors.New("no event files")
	}
	return nil
}
