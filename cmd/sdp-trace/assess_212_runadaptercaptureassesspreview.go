package main

import (
	"io"
)

func runAdapterCaptureAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Preview records only whether the run path is present, not whether the run
	// satisfies adapter-capture evidence rules.
	inputs := map[string]string{
		"run": managedInputStatus(opts.stringValue("run")),
	}
	writeJSONPayloadUnchecked(stdout, newAdapterCapturePreviewReport(inputs))
	return previewInputExitCode(inputs)
}
