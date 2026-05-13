package repoobserver

import (
	"os"
	"path/filepath"
)

func writeNewTarget(path string, data []byte, mode os.FileMode) ([]DiffSummary, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	// New generated files do not need a force diff summary because there is no
	// previous repository content to overwrite.
	return nil, os.WriteFile(path, data, mode)
}
