package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
)

func newForensicPreviewReport(inputs map[string]string) forensicPreviewReport {
	// Forensic preview documents retention policy effects without executing
	// redaction or exposing matched sensitive values.
	return forensicPreviewReport{
		Command:         "assess preview",
		SelectedProfile: forensic.ProfileForensicRetention,
		Inputs:          inputs,
		PolicyEffects:   forensicPreviewPolicyEffects,
		NextActions:     forensicPreviewActions(inputs),
		Claim:           "preview is read-only and does not emit a forensic verdict",
	}
}
