package witness

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func WriteGitHubActions(outPath, runsRoot, reportDir string, env map[string]string) (Record, error) {
	// Writing a witness is explicit: no output path means no durable claim is
	// produced.
	if strings.TrimSpace(outPath) == "" {
		return Record{}, errors.New("witness requires --out <file>")
	}
	record, err := BuildGitHubActions(runsRoot, reportDir, env)
	if err != nil {
		return Record{}, err
	}
	// Safety finalization runs after building the provider record and before the
	// durable artifact is written.
	record = finalizeRecordForWrite(record)
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return Record{}, err
	}
	// Persist only after the complete record is known so partial files never
	// stand in for verified witness evidence.
	return record, writeRecord(outPath, record)
}

func writeRecord(outPath string, record Record) error {
	// writeRecord is intentionally serialization-only; profile decisions and
	// safety checks must already be complete at this point.
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, append(data, '\n'), 0o644)
}

func Load(path string) (Record, error) {
	var record Record
	// Loading parses the JSON shape only; verification still belongs to the
	// caller.
	data, err := os.ReadFile(path)
	if err != nil {
		return Record{}, err
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return Record{}, err
	}
	// The caller decides whether the parsed witness is sufficient for a gate.
	return record, nil
}

func IsPassingCI(record Record) bool {
	// A passing witness must be both successful and CI-witnessed; local records
	// cannot satisfy CI-backed trust on status alone.
	return record.Kind == KindGitHubActions &&
		record.Status == StatusPass &&
		record.TrustScope == TrustScopeCIWitnessed
}
