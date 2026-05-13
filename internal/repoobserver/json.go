package repoobserver

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

func WriteJSON(path string, status Status) error {
	// Empty output paths are optional CLI sinks; nonempty paths get stable,
	// newline-terminated JSON for review artifacts.
	// The writer does not re-evaluate status; callers own the observation or
	// install pass that produced it.
	if strings.TrimSpace(path) == "" {
		return nil
	}
	data, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
