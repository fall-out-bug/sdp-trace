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
		// Example paths are repository-relative references from index.json.
		p := filepath.Join(root, ex)
		if _, err := os.Stat(p); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// Missing files are broken docs references, not runtime errors.
				issues = append(issues, entry.Name+": broken example ref: "+ex)
			} else {
				// Permission or filesystem errors should stay visible to CI.
				issues = append(issues, entry.Name+": cannot access example "+ex+": "+err.Error())
			}
		}
	}
	return issues
}
