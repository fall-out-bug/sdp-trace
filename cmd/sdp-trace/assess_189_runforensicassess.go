package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/forensic"
)

func runForensicAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireForensicAssessInputs(opts, stderr) {
		// A forensic verdict without a redaction policy would overclaim retention
		// coverage.
		return exitUsage
	}
	// Forensic retention assessment requires both policy and run evidence so
	// missing redaction rules cannot be treated as passing defaults.
	input, err := loadForensicInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := forensic.Evaluate(input)
	return writeAssessmentArtifact(opts.stringValue("out"), result, stdout, stderr, forensicExitCode)
}
