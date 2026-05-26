package main

import (
	"fmt"
	"path/filepath"
)

func evalSymlinkPath(path, label string) (string, error) {
	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("%s resolution failed: %w", label, err)
	}
	return resolvedPath, nil
}
