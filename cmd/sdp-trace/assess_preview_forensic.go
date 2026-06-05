package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
	"io"
)

func runForensicAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Forensic preview reports only run/policy presence and leaves redaction
	// evaluation to the real assessment command.
	inputs := map[string]string{
		"run":              managedInputStatus(opts.stringValue("run")),
		"redaction_policy": managedInputStatus(opts.stringValue("redaction-policy")),
	}
	writeJSONPayloadUnchecked(stdout, newForensicPreviewReport(inputs))
	return previewInputExitCode(inputs)
}

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

var forensicPreviewPolicyEffects = map[string]string{
	"redaction_engine": "not_executed_in_preview",
	"matched_values":   "not_rendered",
	"rule_refs":        "shown_when_present_in_policy_or_run_metadata",
	"retention_modes":  "digest_only,sanitized_excerpt,encrypted_raw_ref,external_artifact_ref,not_assessed",
}
