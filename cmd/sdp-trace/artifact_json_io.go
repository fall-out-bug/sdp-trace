package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func readJSONFile(path string, dst any) error {
	// Shared JSON reads are local artifact loads; callers decide whether a
	// failure is usage, cannot_verify, or ordinary command failure.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func writeJSONFile(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Pretty JSON keeps generated evidence reviewable and stable in fixtures.
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
