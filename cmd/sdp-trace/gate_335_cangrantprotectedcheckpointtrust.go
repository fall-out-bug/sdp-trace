package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func canGrantProtectedCheckpointTrust(result checkpoint.VerificationResult, signed checkpoint.SignedCheckpoint, policy checkpoint.TrustedCheckpointPolicy, witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	// Do not upgrade explicit checkpoint failures; protected trust can only
	// tighten cannot-verify/local states when all external bindings match.
	return result.Result != checkpoint.StateFail &&
		signed.Signer.Authority == checkpoint.AuthorityCIIsolatedJob &&
		policyAllowsSigner(policy, signed) &&
		witnessMatchesProtectedInput(witnessSummary, expected)
}
