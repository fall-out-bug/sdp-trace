package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func runStandardGate(target, outPath string, opts *flagSet, stdout, stderr io.Writer) int {
	// Standard gates derive local/CI/audit verdicts from demo run evidence and
	// optional witness data; protected-profile trust is handled by runProtectedGate.
	result, err := demo.WriteGate(target, outPath, opts.stringValue("contract"), opts.stringValue("witness"))
	if err != nil {
		// WriteGate returns user-facing diagnostics for missing output paths,
		// unreadable targets, contract errors, and artifact write failures.
		fmt.Fprintln(stderr, err)
		return 1
	}
	// Use the shared unchecked writer because GateResult is a concrete JSON-safe
	// value; the helper preserves the existing indented payload and final newline.
	writeJSONPayloadUnchecked(stdout, result)
	return gateExitCode(result)
}

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
