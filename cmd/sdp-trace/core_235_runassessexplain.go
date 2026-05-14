package main

import (
	"fmt"
	"io"
)

func runAssessExplain(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "assess explain"}
	opts.setString("assessment-result", "")
	// Explanation renders an existing assessment artifact only; it never
	// re-evaluates evidence or upgrades a verdict.
	path, err := parseAssessExplainArgs(opts, args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	return explainAssessmentResult(path, stdout, stderr)
}
