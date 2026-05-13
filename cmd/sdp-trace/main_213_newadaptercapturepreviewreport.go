package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
)

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
