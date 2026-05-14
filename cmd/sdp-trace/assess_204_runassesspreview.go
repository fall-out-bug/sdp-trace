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
