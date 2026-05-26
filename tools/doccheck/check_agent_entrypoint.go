package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func checkAgentEntrypoint(help string, registry []string) error {
	docPath := filepath.Join(repoRoot(), agentEntrypoint)
	doc, err := os.ReadFile(docPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", agentEntrypoint, err)
	}
	return compareAgentEntrypoint(help, string(doc), registry)
}
