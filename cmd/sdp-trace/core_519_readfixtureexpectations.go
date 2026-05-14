package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

func readFixtureExpectations(root string) (map[string]string, error) {
	path := filepath.Join(root, "fixture-expectations.json")
	// Fixture expectations are optional metadata outside the verifier result.
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Missing expectations file means defaults, not a broken corpus.
			return nil, nil
		}
		return nil, err
	}
	var expectations map[string]string
	if err := json.Unmarshal(data, &expectations); err != nil {
		// Malformed expectation metadata is reported to the fixture validator.
		return nil, err
	}
	return expectations, nil
}

// flagSet is a tiny local parser for limited CLI needs.
