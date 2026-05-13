package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func newProtectedGatePreviewReport(inputs map[string]string) protectedGatePreviewReport {
	// The report deliberately mirrors the gate command vocabulary while keeping
	// the claim explicit that no protected verdict was produced.
	return protectedGatePreviewReport{
		Command:         "gate preview",
		SelectedProfile: demo.GateProfileProtected,
		TrustCap:        string(trace.TrustScopeLocalObserved),
		Inputs:          inputs,
		NextActions:     protectedPreviewActions(inputs),
		Claim:           "preview is read-only and does not emit a protected verdict",
	}
}
