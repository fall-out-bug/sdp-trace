package main

import (
	"io"
)

func runForensicAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Forensic preview reports only run/policy presence and leaves redaction
	// evaluation to the real assessment command.
	inputs := map[string]string{
		"run":              managedInputStatus(opts.stringValue("run")),
		"redaction_policy": managedInputStatus(opts.stringValue("redaction-policy")),
	}
	writeJSONPayloadUnchecked(stdout, newForensicPreviewReport(inputs))
	return previewInputExitCode(inputs)
}
