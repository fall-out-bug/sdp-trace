package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func expectedFixtureResultFailed(runDir string, result trace.VerifierResult, expectation fixtureExpectation, stderr io.Writer) bool {
	if expectation.ExpectedResult == string(result.Result) {
		return false
	}
	// Mismatches are printed with the run path so fixture corpus drift is
	// actionable.
	fmt.Fprintf(stderr, "%s expected %s, got %s\n", runDir, expectation.ExpectedResult, result.Result)
	return true
}
