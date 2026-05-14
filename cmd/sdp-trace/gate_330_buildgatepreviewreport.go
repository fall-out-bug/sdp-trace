package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func buildGatePreviewReport(contract trace.Contract, witnessPath, target string) gatePreviewReport {
	// The preview report is a planning artifact: it names required runs and
	// evidence IDs without claiming the gate will pass.
	report := gatePreviewReport{
		Command:          "gate preview",
		GateMode:         previewGateMode(contract),
		TrustCap:         string(trace.TrustScopeLocalObserved),
		RequiredRuns:     requiredRunIDs(contract),
		RequiredEvidence: requiredEvidenceIDsForCLI(contract),
		Claim:            "preview is read-only and does not claim the gate will pass",
	}
	if witnessPath != "" {
		// Optional witness inspection checks binding shape only; it does not
		// produce a CI-witness gate verdict.
		report.WitnessInspectable, report.WitnessMismatches = demo.PreviewWitnessBinding(witnessPath, target)
	}
	return report
}
