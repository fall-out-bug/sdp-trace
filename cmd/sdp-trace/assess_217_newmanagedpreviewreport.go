package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

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
