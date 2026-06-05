package main

func managedPreviewActions(inputs map[string]string) []string {
	// Managed preview actions name setup gaps for each required evidence class
	// without deriving a managed-harness verdict.
	return previewActionsForInputs(
		inputs,
		[]string{"run", "adapter_registry", "managed_policy", "managed_witness"},
		"Supply %s before managed assessment.",
		"Fix %s so it is readable JSON or a run directory.",
	)
}

func forensicPreviewActions(inputs map[string]string) []string {
	// Forensic preview actions stay limited to setup remediation; the retention
	// assessment itself still requires a full policy/run evaluation.
	return previewActionsForInputs(
		inputs,
		[]string{"run", "redaction_policy"},
		"Supply %s before forensic retention assessment.",
		"Fix %s so it is readable JSON or a run directory.",
	)
}

func adapterCapturePreviewActions(inputs map[string]string) []string {
	return previewActionForInputState(inputs, "run", "Supply run before adapter capture assessment.", "Fix run so it is a readable JSON run directory.")
}

func ciArtifactPreviewActions(inputs map[string]string) []string {
	return previewActionForInputState(inputs, "artifact_manifest", "Supply artifact manifest before CI artifact observation assessment.", "Fix artifact manifest so it is readable JSON.")
}

func authorityPreviewActions(inputs map[string]string) []string {
	return previewActionForInputState(inputs, "authority_package", "Supply authority package before authority envelope assessment.", "Fix authority package so it is readable JSON.")
}
