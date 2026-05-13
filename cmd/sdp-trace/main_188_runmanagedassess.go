package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

func runManagedAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireManagedAssessInputs(opts, stderr) {
		// Managed assessment has no implicit defaults for registry, policy, or
		// witness authority.
		return exitUsage
	}
	// Managed-harness assessment joins contract, policy, registry, run, and
	// witness evidence before deriving a trust state.
	input, err := loadManagedInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := managed.Evaluate(input)
	return writeAssessmentArtifact(opts.stringValue("out"), result, stdout, stderr, managedExitCode)
}
