package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func checkQuickstart(registry []string) error {
	qsPath := filepath.Join(repoRoot(), contributorQuickstart)
	qs, err := os.ReadFile(qsPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", contributorQuickstart, err)
	}
	return compareQuickstartWithRegistry(string(qs), registry)
}
