package main

import (
	"path/filepath"
)

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
