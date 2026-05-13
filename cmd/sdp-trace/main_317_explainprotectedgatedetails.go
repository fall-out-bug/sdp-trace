package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

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
