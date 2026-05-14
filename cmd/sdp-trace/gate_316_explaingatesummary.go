package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func explainGateSummary(result demo.GateResult, stdout io.Writer) {
	// Summary lines keep the layered gate states distinct so local, CI witness,
	// audit, and protected outcomes are not collapsed into one score.
	fmt.Fprintf(stdout, "Gate mode: %s\n", result.GateMode)
	fmt.Fprintf(stdout, "Trust cap: %s\n", result.TrustCap)
	if result.SelectedProfile != "" {
		fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	}
	fmt.Fprintf(stdout, "Local gate: %s\n", result.LocalGate)
	fmt.Fprintf(stdout, "CI witness gate: %s\n", result.CIWitnessGate)
	fmt.Fprintf(stdout, "Audit-grade gate: %s\n", result.AuditGradeGate)
	if result.ProtectedGate != "" {
		fmt.Fprintf(stdout, "Protected gate: %s\n", result.ProtectedGate)
	}
}
