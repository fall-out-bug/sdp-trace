package main

import (
	"encoding/json"
	"os"
)

// writeJSONFile persists generated baselines in a review-stable JSON shape.
func writeJSONFile(path string, value any) error {
	// Baseline files are deterministic JSON artifacts so review diffs show only
	// real ratchet movement.
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	// A trailing newline keeps generated baseline files friendly to POSIX tools.
	payload = append(payload, '\n')
	return os.WriteFile(path, payload, 0o600)
}
