package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
	"io"
)

func runAdapterCaptureAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Preview records only whether the run path is present, not whether the run
	// satisfies adapter-capture evidence rules.
	inputs := map[string]string{
		"run": managedInputStatus(opts.stringValue("run")),
	}
	writeJSONPayloadUnchecked(stdout, newAdapterCapturePreviewReport(inputs))
	return previewInputExitCode(inputs)
}

func newAdapterCapturePreviewReport(inputs map[string]string) adapterCapturePreviewReport {
	// Adapter preview lists expected vocabulary and safety constraints only; it
	// does not inspect raw adapter payloads or issue a verdict.
	return adapterCapturePreviewReport{
		Command:          "assess preview",
		SelectedProfile:  adaptercapture.ProfileAdapterCapture,
		Inputs:           inputs,
		ExpectedEvidence: adapterCapturePreviewExpectedEvidence,
		Safety:           adapterCapturePreviewSafety,
		NextActions:      adapterCapturePreviewActions(inputs),
		Claim:            "preview is read-only and does not emit an adapter capture verdict",
	}
}

var adapterCapturePreviewExpectedEvidence = map[string]string{
	"binding_modes":        "same_chain,adapter_bundle",
	"test_provenance":      "ci_executed,wrapper_executed,harness_observed,agent_reported,cannot_verify",
	"capture_depth_states": "captured,missing_telemetry,not_integrated,unsupported,retention_limited,not_assessed,cannot_verify",
}

var adapterCapturePreviewSafety = map[string]string{
	"raw_payloads":    "not_rendered",
	"adapter_secrets": "not_rendered",
	"gateway_refs":    "token_free_refs_only",
	"model_payloads":  "digest_or_block18_reference_only",
}
