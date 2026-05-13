package main

import (
	"context"
	"fmt"
	"io"
)

func runAssess(_ context.Context, args []string, stdout, stderr io.Writer) int {
	// Assess accepts either a documented subcommand or a profile-bound verdict
	// run; no positional assessment payloads are accepted.
	if code, ok := runAssessSubcommand(args, stdout, stderr); ok {
		return code
	}
	opts, ok := parseAssessOptions(args, stderr)
	if !ok {
		return exitUsage
	}
	// The profile flag is the assessment boundary: each profile has a distinct
	// evidence shape and exit-code policy.
	handler, ok := assessHandlers()[opts.stringValue("profile")]
	if !ok {
		fmt.Fprintln(stderr, "assess requires --profile adapter-capture, managed-harness, forensic-retention, ci-artifact-observation, or authority-envelope")
		return exitUsage
	}
	return handler(opts, stdout, stderr)
}
