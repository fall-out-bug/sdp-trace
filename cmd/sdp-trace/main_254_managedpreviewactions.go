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
