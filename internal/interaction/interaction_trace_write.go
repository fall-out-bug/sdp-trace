package interaction

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

func WriteTrace(path string, trace Trace) error {
	// WriteTrace keeps interaction trace evidence explicit and source-bound.
	// Relay, import, validation, envelope, summary, and retention states stay separate.
	// This helper records or renders observed interaction data; it does not create external proof.

	if strings.TrimSpace(path) == "" {
		return errors.New("interaction command requires --out")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(trace, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}
