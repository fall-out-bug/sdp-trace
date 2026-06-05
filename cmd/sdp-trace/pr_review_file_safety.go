package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func requireOutputFile(command, path string) error {
	if strings.TrimSpace(path) == "" {
		// Commands that produce artifacts require an explicit destination to
		// avoid pretending stdout-only output is persisted evidence.
		return fmt.Errorf("%s requires --out", command)
	}
	return refuseExistingFile(path)
}

func refuseExistingFile(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			// Directories are never valid write-once text artifact targets.
			return fmt.Errorf("output path is a directory: %s", path)
		}
		return fmt.Errorf("output file exists: %s", path)
	}
	if errors.Is(err, os.ErrNotExist) {
		// Missing path is the only acceptable state for write-once artifacts.
		return nil
	}
	return err
}

func requireDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("work-dir: %w", err)
	}
	if !info.IsDir() {
		// Runner working directory must be a directory, not a file path.
		return fmt.Errorf("work-dir is not a directory: %s", path)
	}
	return nil
}
