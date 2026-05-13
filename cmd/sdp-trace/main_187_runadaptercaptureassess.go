package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
)

func runAdapterCaptureAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireAdapterCaptureAssessInputs(opts, stderr) {
		// Missing durable input/output flags are usage failures before
		// adaptercapture can evaluate run evidence.
		return exitUsage
	}
	// Adapter-capture assessment is run-bound only; missing run evidence is a
	// usage error before any verdict artifact exists.
	input, err := loadAdapterCaptureInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := adaptercapture.Evaluate(input)
	return writeAssessmentArtifact(opts.stringValue("out"), result, stdout, stderr, adapterCaptureExitCode)
}
