package main

import (
	"errors"
	"fmt"
	"os"
)

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
