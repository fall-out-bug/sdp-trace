package main

import (
	"fmt"
	"os"
)

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
