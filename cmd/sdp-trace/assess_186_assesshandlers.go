package main

func assessHandlers() map[string]assessHandler {
	// Handler keys are product profile names, not implementation package names,
	// so CLI output remains aligned with documented assessment profiles.
	return map[string]assessHandler{
		"adapter-capture":         runAdapterCaptureAssess,
		"managed-harness":         runManagedAssess,
		"forensic-retention":      runForensicAssess,
		"ci-artifact-observation": runCIArtifactAssess,
		"authority-envelope":      runAuthorityAssess,
	}
}
