package main

import (
	"fmt"
	"io"
)

func runAssessPreview(args []string, stdout, stderr io.Writer) int {
	opts, ok := parseAssessPreviewOptions(args, stderr)
	if !ok {
		return exitUsage
	}
	// Preview is advisory setup output. It must not imply that assessment
	// evidence has been evaluated.
	handler, ok := assessPreviewHandlers()[opts.stringValue("profile")]
	if !ok {
		fmt.Fprintln(stderr, "assess preview requires --profile adapter-capture, managed-harness, forensic-retention, ci-artifact-observation, or authority-envelope")
		return exitUsage
	}
	return handler(opts, stdout)
}

func parseAssessPreviewOptions(args []string, stderr io.Writer) (*flagSet, bool) {
	opts := newStringFlagSet("assess preview", assessPreviewStringFlags)
	if err := opts.parse(args); err != nil {
		// Preview parse failures are usage errors, not assessment verdicts.
		fmt.Fprintln(stderr, err)
		return nil, false
	}
	if len(opts.rest()) != 0 {
		// Preview output is generated from named local paths only; positional
		// prose is not evidence.
		fmt.Fprintln(stderr, "assess preview accepts only flags")
		return nil, false
	}
	return opts, true
}
