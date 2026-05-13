package checkpoint

func VerifySet(runDir string, checkpoints []SignedCheckpoint, policy *TrustedCheckpointPolicy) VerificationResult {
	// VerifySet keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	result := baseSetResult()
	if len(checkpoints) == 0 {
		return emptySetResult(result)
	}
	runID := checkpoints[0].RunID
	previousDigest := ""
	for i, cp := range checkpoints {

		if stop := verifySetCheckpoint(&result, runDir, cp, policy, setLinkExpectation{runID: runID, sequence: i, previousDigest: previousDigest}); stop {
			return result
		}
		previousDigest = cp.PayloadDigest
	}
	return result
}

func baseSetResult() VerificationResult {
	// baseSetResult keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	return VerificationResult{
		SchemaVersion:        VerificationSchemaVersion,
		Result:               StatePass,
		TrustScope:           TrustScopeLocalSigned,
		SequenceState:        StatePass,
		ReplayFreshnessState: StateNotAssessed,
		SignerAuthorityState: StateNotAssessed,
		Reasons:              []string{},
	}
}

func emptySetResult(result VerificationResult) VerificationResult {
	// emptySetResult keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	result.Result = StateCannotVerify
	result.SequenceState = StateCannotVerify
	result.Reasons = append(result.Reasons, "no checkpoints supplied")
	return result
}

type setLinkExpectation struct {
	runID          string
	sequence       int
	previousDigest string
}
