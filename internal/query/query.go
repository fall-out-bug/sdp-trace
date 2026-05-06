package query

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/verifier"
)

const (
	QueryMissingEvidence = "missing-evidence"
)

// MissingEvidence returns the missing-evidence table for a run.
//
// If a verifier artifact already exists, that result is reused so query output is
// stable between invocations.
func MissingEvidence(runDir string) ([]byte, error) {
	artifactPath := filepath.Join(runDir, "verifier", "missing-evidence-table.json")
	if _, err := os.Stat(artifactPath); err == nil {
		return os.ReadFile(artifactPath)
	}
	_, table, _, err := verifier.VerifyRun(runDir)
	if err != nil {
		return nil, err
	}
	payload, err := json.MarshalIndent(table, "", "  ")
	if err != nil {
		return nil, err
	}
	return payload, nil
}
