package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type fixtureExpectation struct {
	ExpectedResult string `json:"expected_result"`
}

func readFixtureExpectation(root, runDir string) (fixtureExpectation, error) {
	// Expectations are optional corpus metadata; absence leaves default verifier
	// failure handling in place.
	expectations, err := readFixtureExpectations(root)
	if err != nil {
		return fixtureExpectation{}, err
	}
	if len(expectations) == 0 {
		return fixtureExpectation{}, nil
	}
	name := filepath.Base(runDir)
	// Fixture expectations are keyed by run directory basename so the corpus can
	// move as a whole.
	return fixtureExpectation{ExpectedResult: expectations[name]}, nil
}

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
