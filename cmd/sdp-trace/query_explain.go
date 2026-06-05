package main

import (
	"context"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/verifier"
)

func runExplain(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		// Explanation is anchored to a retained run directory, not free-form
		// evidence text.
		fmt.Fprintln(stderr, "explain requires <run-dir>")
		return exitUsage
	}
	runDir := args[0]
	// ExplainRun renders a derived human view and does not change verifier
	// artifacts or verdict state.
	explanation, err := verifier.ExplainRun(runDir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintln(stdout, explanation)
	return 0
}
