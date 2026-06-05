package main

import (
	"errors"
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

func parseAssessExplainArgs(opts *flagSet, args []string) (string, error) {
	if err := opts.parse(args); err != nil {
		return "", err
	}
	if len(opts.rest()) != 0 {
		return "", errors.New("assess explain accepts only flags")
	}
	path := opts.stringValue("assessment-result")
	if path == "" {
		// A concrete artifact path is required so operators cannot ask the
		// explainer to interpret ad hoc text as evidence.
		return "", errors.New("assess explain requires --assessment-result <file>")
	}
	return path, nil
}
