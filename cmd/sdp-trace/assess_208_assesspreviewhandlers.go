package main

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
