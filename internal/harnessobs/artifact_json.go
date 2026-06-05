package harnessobs

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// writeJSON writes replayable harness artifacts with stable formatting and the
// same trailing newline used by checked fixtures and generated evidence.
func writeJSON(path string, value any) error {
	data, err := jsonArtifactData(value)
	if err != nil {
		return err
	}
	return writeJSONFile(path, data)
}

// jsonArtifactData isolates serialization from filesystem effects so tests can
// assert artifact shape without depending on directory creation.
func jsonArtifactData(value any) ([]byte, error) {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

// writeJSONFile creates the artifact parent on demand; callers still own path
// safety before handing a destination to this writer.
func writeJSONFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
