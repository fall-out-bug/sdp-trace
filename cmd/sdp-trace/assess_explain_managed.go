package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/managed"
)

func explainManagedAssessment(result managed.AssessmentResult, stdout io.Writer) int {
	// Managed explanations keep condition state and remediation text visible
	// without treating setup readability as an independent pass claim.
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Managed harness assessment: %s\n", result.ManagedHarnessAssessment)
	fmt.Fprintf(stdout, "Trust scope: %s\n", result.TrustScope)
	for _, condition := range result.ManagedConditions {
		fmt.Fprintf(stdout, "Managed condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	}
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
	return 0
}
