package witness

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func runIDsFromDirs(runDirs []string) ([]string, error) {
	runIDs := make([]string, 0, len(runDirs))
	for _, runDir := range runDirs {
		// Skip run directories that lack an ID but fail on unreadable or
		// malformed run.json, preserving absent versus bad evidence.
		runID, ok, err := nonEmptyRunIDFromDir(runDir)
		if err != nil {
			return nil, err
		}
		if ok {
			runIDs = append(runIDs, runID)
		}
	}
	return runIDs, nil
}

func nonEmptyRunIDFromDir(runDir string) (string, bool, error) {
	runID, err := runIDFromDir(runDir)
	if err != nil {
		return "", false, err
	}
	// Empty run IDs are skipped instead of being treated as a wildcard match.
	return runID, runID != "", nil
}

func runIDFromDir(runDir string) (string, error) {
	// Accept both run_id and legacy id to keep old evidence replayable without
	// weakening the requirement that some run identity be present.
	// The returned ID is trimmed so whitespace-only legacy IDs do not bind.
	// Malformed run.json is an error rather than a non-match because it blocks
	// reliable run binding.
	raw, err := os.ReadFile(filepath.Join(runDir, "run.json"))
	if err != nil {
		return "", err
	}
	var payload struct {
		RunID string `json:"run_id"`
		ID    string `json:"id"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", err
	}
	if strings.TrimSpace(payload.RunID) != "" {
		return payload.RunID, nil
	}
	return strings.TrimSpace(payload.ID), nil
}
