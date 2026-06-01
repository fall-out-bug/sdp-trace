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

func expectedFixtureResultFailed(runDir string, result trace.VerifierResult, expectation fixtureExpectation, stderr io.Writer) bool {
	if expectation.ExpectedResult == string(result.Result) {
		return false
	}
	// Mismatches are printed with the run path so fixture corpus drift is
	// actionable.
	fmt.Fprintf(stderr, "%s expected %s, got %s\n", runDir, expectation.ExpectedResult, result.Result)
	return true
}

func unexpectedFixtureResultFailed(result trace.VerifierResult) bool {
	// Without an explicit expected result, only fail/cannot_verify are treated as
	// fixture failures; observed/not_assessed remain inspectable but nonfatal.
	return result.Result == trace.VerdictFail || result.Result == trace.VerdictCannotVerify
}
