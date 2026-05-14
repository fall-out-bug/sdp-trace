package main

import (
	"strings"
)

func writeOptionalJSONFile(path string, value any) error {
	if strings.TrimSpace(path) == "" {
		// Optional JSON outputs are side effects; an omitted path must not change
		// the command verdict.
		return nil
	}
	return writeJSONFile(path, value)
}
