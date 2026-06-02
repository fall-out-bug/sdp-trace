package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func explainGateResult(result demo.GateResult, stdout io.Writer) {
	if result.SchemaVersion == demo.GateSchemaVersion {
		// Legacy gate results have no protected-profile fields; state that
		// absence explicitly rather than implying not_assessed conditions.
		fmt.Fprintln(stdout, "Protected profile fields: absent")
	}
	explainGateSummary(result, stdout)
	explainProtectedGateDetails(result, stdout)
	explainGateCollections(result, stdout)
}

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

func explainProtectedGateDetails(result demo.GateResult, stdout io.Writer) {
	if result.CheckpointVerification != nil {
		// Checkpoint verification is shown separately from protected conditions
		// because replay/signature failure and policy failure are different
		// evidence gaps.
		fmt.Fprintf(stdout, "Checkpoint result: %s\n", result.CheckpointVerification.Result)
		fmt.Fprintf(stdout, "Checkpoint trust scope: %s\n", result.CheckpointVerification.TrustScope)
	}
	for _, condition := range result.ProtectedConditions {
		// Protected conditions remain individual rows for auditability.
		fmt.Fprintf(stdout, "Protected condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	}
}
