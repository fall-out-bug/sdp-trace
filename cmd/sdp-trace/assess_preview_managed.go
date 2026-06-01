package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/managed"
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

func newManagedPreviewReport(inputs map[string]string) managedPreviewReport {
	// Managed preview exposes which local artifacts are ready to be assessed,
	// while keeping capability and witness verdicts uncomputed.
	return managedPreviewReport{
		Command:         "assess preview",
		SelectedProfile: managed.ProfileManagedHarness,
		Inputs:          inputs,
		NextActions:     managedPreviewActions(inputs),
		Claim:           "preview is read-only and does not emit a managed verdict",
	}
}
