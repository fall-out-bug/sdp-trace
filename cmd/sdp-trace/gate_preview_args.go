package main

import (
	"fmt"
	"io"
)

var gatePreviewStringFlags = []string{"contract", "witness", "profile", "checkpoint", "checkpoint-policy"}

func parseGatePreviewArgs(args []string, stderr io.Writer) (*flagSet, []string, int, bool) {
	// Preview accepts both standard and protected flags because it reports setup
	// readiness without committing to a verdict mode.
	opts := newStringFlagSet("gate preview", gatePreviewStringFlags)
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, nil, exitUsage, false
	}
	targets := opts.rest()
	if len(targets) != 1 {
		// Preview is still target-scoped so witness-binding checks, when
		// requested, compare against one run root.
		fmt.Fprintln(stderr, "gate preview requires <runs-root-or-run-dir>")
		return nil, nil, exitUsage, false
	}
	return opts, targets, 0, true
}
