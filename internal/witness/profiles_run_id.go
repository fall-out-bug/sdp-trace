package witness

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

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
	payload, err := decodeRunIdentity(raw)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(payload.RunID) != "" {
		return payload.RunID, nil
	}
	return strings.TrimSpace(payload.ID), nil
}

func decodeRunIdentity(raw []byte) (runIdentityPayload, error) {
	// Keep run-id decoding isolated so legacy id fallback stays reviewable.
	var payload runIdentityPayload
	err := json.Unmarshal(raw, &payload)
	return payload, err
}

type runIdentityPayload struct {
	RunID string `json:"run_id"`
	ID    string `json:"id"`
}
