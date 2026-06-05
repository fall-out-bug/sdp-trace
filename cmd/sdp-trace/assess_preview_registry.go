package main

import (
	"io"
)

var assessPreviewStringFlags = []string{
	"profile",
	"out",
	"run",
	"adapter-registry",
	"managed-policy",
	"managed-witness",
	"redaction-policy",
	"artifact-manifest",
	"authority-package",
}

type assessPreviewHandler func(*flagSet, io.Writer) int

func assessPreviewHandlers() map[string]assessPreviewHandler {
	// Preview handlers mirror assess profile names for command-surface parity.
	return map[string]assessPreviewHandler{
		"adapter-capture":         runAdapterCaptureAssessPreview,
		"managed-harness":         runManagedAssessPreview,
		"forensic-retention":      runForensicAssessPreview,
		"ci-artifact-observation": runCIArtifactAssessPreview,
		"authority-envelope":      runAuthorityAssessPreview,
	}
}

func previewInputExitCode(inputs map[string]string) int {
	for _, state := range inputs {
		if previewInputCannotVerify(state) {
			// Bad preview inputs block setup confidence without emitting a profile
			// assessment verdict.
			return exitCannotVerify
		}
	}
	return 0
}
