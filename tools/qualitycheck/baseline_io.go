package main

import (
	"encoding/json"
	"os"
)

// readJSONFile reads typed baseline JSON without giving callers partial state.
func readJSONFile[T any](path string) (T, error) {
	var value T
	// Baseline reads preserve the zero value on failure so callers never compare
	// against partially decoded ratchet data.
	payload, err := os.ReadFile(path)
	if err != nil {
		return value, err
	}
	// Decode into the typed baseline selected by the caller; schema validation
	// happens after this structural parse succeeds.
	if err := json.Unmarshal(payload, &value); err != nil {
		return value, err
	}
	return value, nil
}
