package main

import (
	"io"
)

func runManagedAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Managed preview reports local artifact readiness without evaluating
	// capability, registry, run, or witness state.
	inputs := map[string]string{
		"run":              managedInputStatus(opts.stringValue("run")),
		"adapter_registry": managedInputStatus(opts.stringValue("adapter-registry")),
		"managed_policy":   managedInputStatus(opts.stringValue("managed-policy")),
		"managed_witness":  managedInputStatus(opts.stringValue("managed-witness")),
	}
	writeJSONPayloadUnchecked(stdout, newManagedPreviewReport(inputs))
	return previewInputExitCode(inputs)
}
