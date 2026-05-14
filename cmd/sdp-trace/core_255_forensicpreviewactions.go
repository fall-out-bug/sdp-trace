package main

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
