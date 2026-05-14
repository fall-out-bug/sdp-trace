package main

import (
	"errors"
)

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
