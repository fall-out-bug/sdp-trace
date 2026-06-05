package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/verifier"
)

func runVerify(_ context.Context, args []string, stdout, stderr io.Writer) int {
	runDir, code, ok := parseVerifyArgs(args, stderr)
	if !ok {
		return code
	}
	// VerifyRun computes the verdict and derived artifacts from retained run
	// evidence; artifact writing happens even when verification reports a
	// semantic error so failures remain inspectable.
	result, table, audit, err := verifier.VerifyRun(runDir)
	if writeErr := verifier.WriteVerifierArtifacts(runDir, result, table, audit); writeErr != nil {
		fmt.Fprintf(stderr, "failed writing verifier artifacts for %s: %v\n", runDir, writeErr)
		return 1
	}
	// Stdout carries the structured verifier result after artifacts are written
	// so terminal consumers cannot observe a result that was not retained.
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", data)
	if err != nil {
		// The JSON result is still emitted before the diagnostic so automation
		// can capture structured state.
		fmt.Fprintf(stderr, "%v\n", err)
	}
	return verifierResultExitCode(result.Result)
}
