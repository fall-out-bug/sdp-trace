package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// validateRunDirUnderTmp ensures runDir is safely contained within tmpDir.
func validateRunDirUnderTmp(runDir, tmpDir string) error {
	if err := checkRunDirSafe(runDir); err != nil {
		return err
	}
	runJSONPath := filepath.Join(tmpDir, filepath.Clean(runDir), "run.json")
	if err := checkRunJSONUnderTmp(runJSONPath, tmpDir); err != nil {
		return err
	}
	if _, err := os.Stat(runJSONPath); err != nil {
		return fmt.Errorf("run.json not found at expected path %s: %w", runJSONPath, err)
	}
	return nil
}
