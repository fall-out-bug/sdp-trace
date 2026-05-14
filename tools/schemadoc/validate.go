package main

// validate.go orchestrates per-entry validation: name suffix, file presence,
// status, purpose, example coverage, and example references.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// validateEntry runs all validation rules for a single schema index entry.
func validateEntry(root string, entry *SchemaEntry) []string {
	var issues []string
	if !strings.HasSuffix(entry.Name, ".schema.json") {
		issues = append(issues, fmt.Sprintf("%s: name must end with .schema.json", entry.Name))
	}
	issues = append(issues, validateFilePresence(root, entry)...)
	issues = append(issues, validateEntryStatus(entry)...)
	issues = append(issues, validateEntryPurpose(entry)...)
	issues = append(issues, validateExampleCoverage(entry)...)
	issues = append(issues, validateExampleRefs(root, entry)...)
	return issues
}

// validateFilePresence confirms that the indexed schema file exists on disk.
func validateFilePresence(root string, entry *SchemaEntry) []string {
	path := filepath.Join(root, "schema", entry.Name)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return []string{fmt.Sprintf("extra index entry (file missing): %s", entry.Name)}
		}
		return []string{fmt.Sprintf("%s: cannot access file: %v", entry.Name, err)}
	}
	return nil
}
