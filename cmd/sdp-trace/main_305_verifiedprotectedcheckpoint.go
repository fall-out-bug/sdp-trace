package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func verifiedProtectedCheckpoint(runDir string, inputs protectedGateInputs, expected demo.WitnessExpectation) checkpoint.VerificationResult {
	// Checkpoint replay runs before protected aggregation so signature/policy
	// facts stay separate from gate conditions.
	return protectedCheckpointVerification(
		checkpoint.Verify(runDir, inputs.signed, &inputs.policy),
		inputs.signed,
		inputs.policy,
		inputs.witness,
		expected,
	)
}
