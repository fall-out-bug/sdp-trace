package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func runProtectedGatePreview(opts *flagSet, stdout io.Writer) int {
	// Protected preview classifies required input files before any checkpoint
	// replay or witness trust evaluation.
	inputs := map[string]string{
		"checkpoint":        protectedInputStatus(opts.stringValue("checkpoint")),
		"checkpoint_policy": protectedInputStatus(opts.stringValue("checkpoint-policy")),
		"witness":           protectedInputStatus(opts.stringValue("witness")),
	}
	report := newProtectedGatePreviewReport(inputs)
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	for _, state := range inputs {
		if state == "present_unreadable" || state == "present_malformed" {
			// Bad setup inputs keep preview cannot_verify; they are not lowered
			// into protected gate failures because no verdict was evaluated.
			return exitCannotVerify
		}
	}
	return 0
}

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
