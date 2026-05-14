package main

// refs.go validates that example file references in the index point to existing paths.

import (
	"errors"
	"os"
	"path/filepath"
)

// validateExampleRefs checks every example path for a schema entry.
// It reports broken refs or filesystem access errors.
func validateExampleRefs(root string, entry *SchemaEntry) []string {
	var issues []string
	for _, ex := range entry.Examples {
		p := filepath.Join(root, ex)
		if _, err := os.Stat(p); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				issues = append(issues, entry.Name+": broken example ref: "+ex)
			} else {
				issues = append(issues, entry.Name+": cannot access example "+ex+": "+err.Error())
			}
		}
	}
	return issues
}
