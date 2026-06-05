package prreview

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// WriteJSON persists portable review artifacts as stable, indented JSON. A
// blank path is a deliberate no-op for optional artifact outputs.
func WriteJSON(path string, value any) error {
	if strings.TrimSpace(path) == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

// readJSON is the shared decode boundary for typed artifact readers. It does
// not validate artifact semantics; callers apply type-specific checks.
func readJSON(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}
