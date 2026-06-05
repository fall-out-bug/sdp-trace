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

func canGrantProtectedCheckpointTrust(result checkpoint.VerificationResult, signed checkpoint.SignedCheckpoint, policy checkpoint.TrustedCheckpointPolicy, witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	// Do not upgrade explicit checkpoint failures; protected trust can only
	// tighten cannot-verify/local states when all external bindings match.
	return result.Result != checkpoint.StateFail &&
		signed.Signer.Authority == checkpoint.AuthorityCIIsolatedJob &&
		policyAllowsSigner(policy, signed) &&
		witnessMatchesProtectedInput(witnessSummary, expected)
}
