package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func protectedCheckpointVerification(result checkpoint.VerificationResult, signed checkpoint.SignedCheckpoint, policy checkpoint.TrustedCheckpointPolicy, witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) checkpoint.VerificationResult {
	if !canGrantProtectedCheckpointTrust(result, signed, policy, witnessSummary, expected) {
		return result
	}
	// A protected upgrade is allowed only after checkpoint replay, signer
	// policy, and witness binding all agree on the same run.
	result.SignerAuthorityState = checkpoint.StatePass
	result.TrustScope = checkpoint.TrustScopeCISigned
	result.Result = checkpoint.StatePass
	return result
}
