package interaction

import (
	"encoding/json"
	"fmt"
	"strings"
)

func appendJSONLEventLine(events []Event, text string) ([]Event, error) {
	// Blank JSONL records are skipped without advancing source sequence because
	// they are not evidence rows.
	var event Event
	keep, err := parseJSONLEventLine(text, &event)
	if err != nil || !keep {
		return events, err
	}
	return append(events, event), nil
}

func parseJSONLEventLine(text string, event *Event) (bool, error) {
	// parseJSONLEventLine keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	line := strings.TrimSpace(text)
	if line == "" {
		return false, nil
	}
	if err := json.Unmarshal([]byte(line), event); err != nil {
		return false, err
	}
	return true, nil
}

func validateOrdering(events []Event) error {
	// validateOrdering keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	last := map[string]int{}
	seen := map[string]bool{}
	for _, event := range events {
		key := event.SourceID
		if seen[key] && event.SourceSequence <= last[key] {
			return fmt.Errorf("non-monotonic source_sequence for source %s", key)
		}
		seen[key] = true
		last[key] = event.SourceSequence
	}
	return nil
}
