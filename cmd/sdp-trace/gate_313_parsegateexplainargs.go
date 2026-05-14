package main

import (
	"fmt"
	"io"
)

func parseGateExplainArgs(args []string, stderr io.Writer) (string, int, bool) {
	opts := &flagSet{name: "gate explain"}
	// Explanation is keyed by one existing gate-result artifact.
	opts.setString("gate-result", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return "", exitUsage, false
	}
	if len(opts.rest()) != 0 {
		// Gate explanations accept only an artifact path to avoid mixing
		// positional prose with result evidence.
		// Positional text would not be replayable as a source-bound result.
		fmt.Fprintln(stderr, "gate explain accepts only flags")
		return "", exitUsage, false
	}
	path := opts.stringValue("gate-result")
	if path == "" {
		// A persisted gate result is required because explanation does not
		// synthesize verdicts from loose fields.
		fmt.Fprintln(stderr, "gate explain requires --gate-result <file>")
		return "", exitUsage, false
	}
	// The caller reads and validates the artifact before rendering.
	return path, 0, true
}
