package main

// status.go validates the status and purpose fields of a schema index entry.

import (
	"fmt"
	"strings"
)

// validateEntryStatus checks that the status is present and in the allowed set.
func validateEntryStatus(entry *SchemaEntry) []string {
	if entry.Status == "" {
		return []string{fmt.Sprintf("%s: missing status", entry.Name)}
	}
	if !validStatuses[entry.Status] {
		return []string{fmt.Sprintf("%s: invalid status %q", entry.Name, entry.Status)}
	}
	return nil
}

// validateEntryPurpose checks that the purpose string is non-empty.
func validateEntryPurpose(entry *SchemaEntry) []string {
	if strings.TrimSpace(entry.Purpose) == "" {
		return []string{fmt.Sprintf("%s: missing purpose", entry.Name)}
	}
	return nil
}
