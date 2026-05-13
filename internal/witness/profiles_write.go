package witness

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func WriteProfile(kind, outPath, runsRoot, reportDir string, opts ProfileOptions) (Record, error) {
	if strings.TrimSpace(outPath) == "" {
		return Record{}, errors.New("witness requires --out <file>")
	}
	// Profile files are publication artifacts. Build and safety finalization
	// happen before the write so callers never persist an unscanned trust claim.
	record, err := buildProfileForWrite(kind, runsRoot, reportDir, opts)
	if err != nil {
		return Record{}, err
	}
	if err := writeProfileRecord(outPath, record); err != nil {
		return Record{}, err
	}
	return record, nil
}

func buildProfileForWrite(kind, runsRoot, reportDir string, opts ProfileOptions) (Record, error) {
	// BuildProfile may produce pass, fail, cannot_verify, or not_assessed
	// records. finalizeRecordForWrite is the final output-safety gate that can
	// still replace an otherwise useful record with a redacted failure record.
	record, err := BuildProfile(kind, runsRoot, reportDir, opts)
	if err != nil {
		return Record{}, err
	}
	return finalizeRecordForWrite(record), nil
}

func writeProfileRecord(outPath string, record Record) error {
	// The filesystem write is deliberately dumb: all trust decisions have
	// already been made, and this layer only materializes the exact record.
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(outPath, append(data, '\n'), 0o644); err != nil {
		return err
	}
	return nil
}
