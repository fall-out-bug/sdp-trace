package main

import (
	"context"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/verifier"
)

func runValidateFixtures(_ context.Context, args []string, stdout, stderr io.Writer) int {
	fixtureRoot := fixtureRootArg(args)
	// Fixture discovery is rooted explicitly so validation cannot wander into
	// unrelated run artifacts.
	runDirs, err := demo.DiscoverRunDirs(fixtureRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if validateFixtureRuns(fixtureRoot, runDirs, stdout, stderr) {
		return 1
	}
	return 0
}

func validateFixtureRuns(fixtureRoot string, runDirs []string, stdout, stderr io.Writer) bool {
	failed := false
	for _, runDir := range runDirs {
		// Continue through all fixtures so one broken run does not hide other
		// drift in the example corpus.
		if validateFixtureRun(fixtureRoot, runDir, stdout, stderr) {
			failed = true
		}
	}
	return failed
}

func validateFixtureRun(fixtureRoot, runDir string, stdout, stderr io.Writer) bool {
	result, table, audit, verifyErr := verifier.VerifyRun(runDir)
	if err := verifier.WriteVerifierArtifacts(runDir, result, table, audit); err != nil {
		// Verifier artifacts are part of the fixture evidence, even when replay
		// reports semantic verification errors.
		fmt.Fprintf(stderr, "failed writing verifier artifacts for %s: %v\n", runDir, err)
		return true
	}
	fmt.Fprintf(stdout, "%s => %s\n", runDir, result.Result)
	if verifyErr != nil {
		// Surface replay diagnostics but still compare the structured verdict
		// against the fixture expectation.
		fmt.Fprintf(stderr, "%s verification error: %v\n", runDir, verifyErr)
	}
	return fixtureExpectationFailed(fixtureRoot, runDir, result, stderr)
}
