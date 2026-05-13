package query

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fall_out_bug/sdp-trace/internal/verifier"
)

const QueryMissingEvidence = "missing-evidence"

// MissingEvidence returns the missing-evidence table for a run.
//
// If a verifier artifact already exists, that result is reused so query output is
// stable between invocations.
func MissingEvidence(runDir string) ([]byte, error) {
	// Prefer the verifier artifact when present so query output remains a view
	// of recorded verifier state instead of a silent re-verification.
	path := filepath.Join(runDir, "verifier", "missing-evidence-table.json")
	data, err := os.ReadFile(path)
	if err == nil {
		return data, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}
	return computeMissingEvidence(runDir)
}

func computeMissingEvidence(runDir string) ([]byte, error) {
	// Without a persisted artifact, live verifier replay is the only authority
	// for the missing-evidence table.
	_, table, _, err := verifier.VerifyRun(runDir)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(table, "", "  ")
}
