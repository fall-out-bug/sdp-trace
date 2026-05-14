package main

import (
	"fmt"
	"io"
)

func runPRReviewRun(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePRReviewRunArgs(args, stderr)
	if !ok {
		return code
	}
	// Reviewer execution can only produce usable evidence when packet, profile,
	// runner allow-list, and work directory are all replayable.
	runs, preview, err := executePRReviewRun(opts, args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	writePRReviewRunOutput(stdout, runs, preview)
	return 0
}
