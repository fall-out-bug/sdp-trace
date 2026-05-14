package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/verifier"
)

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
