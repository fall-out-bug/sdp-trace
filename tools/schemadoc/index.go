package main

// index.go handles reading the schema index JSON and listing schema files on disk.
// The index is the canonical source of schema metadata (status, purpose, examples).

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// readIndex parses the schema index from JSON and validates its version.
func readIndex(path string) (*Index, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var idx Index
	if err := json.Unmarshal(b, &idx); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	if idx.Version != indexVersion {
		return nil, fmt.Errorf("unsupported %s version %q, want %q", path, idx.Version, indexVersion)
	}
	return &idx, nil
}

// listSchemaFiles returns sorted basenames of all *.schema.json files under the schema directory.
func listSchemaFiles(root string) ([]string, error) {
	matches, err := filepath.Glob(filepath.Join(root, schemaGlob))
	if err != nil {
		return nil, err
	}
	var names []string
	for _, m := range matches {
		names = append(names, filepath.Base(m))
	}
	sort.Strings(names)
	return names, nil
}
