package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func fixtureExpectationFailed(fixtureRoot, runDir string, result trace.VerifierResult, stderr io.Writer) bool {
	expectation, err := readFixtureExpectation(fixtureRoot, runDir)
	if err != nil {
		// Bad expectation metadata is fixture drift, not a verifier pass.
		fmt.Fprintf(stderr, "invalid fixture expectation for %s: %v\n", runDir, err)
		return true
	}
	if expectation.ExpectedResult != "" {
		// Explicit fixture expectations define the authoritative verdict.
		return expectedFixtureResultFailed(runDir, result, expectation, stderr)
	}
	// Fixtures without explicit expectations may still fail if replay proves a
	// hard verifier failure or cannot-verify state.
	return unexpectedFixtureResultFailed(result)
}
