package main

// check.go validates that index.json matches the schema directory and that every
// entry has required metadata and valid example references.
//
// Steps:
//   1. List schema files on disk.
//   2. Build a map from index entries, detecting duplicates.
//   3. Find schema files missing from the index.
//   4. Validate every indexed entry (presence, status, purpose, examples).
//   5. Return a sorted drift report if any issue is found.
//
// The checker is intentionally small: it reads JSON, walks two slices, and
// returns a multi-line error so CI can print the drift list directly.
//
// Design note: each validation step is a pure function that returns a slice of
// issue strings. This keeps cyclomatic complexity low and makes testing easy.

import (
	"fmt"
	"sort"
	"strings"
)

// check runs the full schema documentation drift check.
func check(root string, idx *Index) error {
	files, err := listSchemaFiles(root)
	if err != nil {
		return err
	}

	indexed, issues := buildIndexedMap(idx.Schemas)
	issues = append(issues, findMissingFromIndex(files, indexed)...)
	issues = append(issues, validateAllEntries(root, idx.Schemas)...)

	if len(issues) > 0 {
		sort.Strings(issues)
		return fmt.Errorf("schema doc drift:\n- %s", strings.Join(issues, "\n- "))
	}
	return nil
}

// buildIndexedMap creates a name-to-entry map and detects duplicate names.
func buildIndexedMap(entries []SchemaEntry) (map[string]*SchemaEntry, []string) {
	indexed := make(map[string]*SchemaEntry, len(entries))
	var issues []string
	for i := range entries {
		// Take the address from the slice element, not the loop variable.
		entry := &entries[i]
		if _, exists := indexed[entry.Name]; exists {
			// Duplicate names make one schema unreachable from the generated map.
			issues = append(issues, "duplicate index entry: "+entry.Name)
			continue
		}
		indexed[entry.Name] = entry
	}
	return indexed, issues
}

// findMissingFromIndex reports schema files on disk that have no index entry.
func findMissingFromIndex(files []string, indexed map[string]*SchemaEntry) []string {
	var issues []string
	for _, f := range files {
		if _, ok := indexed[f]; !ok {
			issues = append(issues, "missing index entry: "+f)
		}
	}
	return issues
}

// validateAllEntries runs per-entry validation for every schema in the index.
func validateAllEntries(root string, entries []SchemaEntry) []string {
	var issues []string
	for i := range entries {
		issues = append(issues, validateEntry(root, &entries[i])...)
	}
	return issues
}
