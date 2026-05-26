package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func checkRunDirSafe(runDir string) error {
	if filepath.IsAbs(runDir) {
		return fmt.Errorf("run_dir is absolute (possible traversal): %q", runDir)
	}
	if strings.HasPrefix(filepath.Clean(runDir), "..") {
		return fmt.Errorf("run_dir escapes tmpDir (possible traversal): %q", runDir)
	}
	return nil
}
