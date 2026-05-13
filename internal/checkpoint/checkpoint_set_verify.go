package checkpoint

import "fmt"

func verifySetCheckpoint(result *VerificationResult, runDir string, cp SignedCheckpoint, policy *TrustedCheckpointPolicy, expected setLinkExpectation) bool {
	// verifySetCheckpoint keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	checkpointResult := Verify(runDir, cp, policy)
	mergeSetVerification(result, checkpointResult)
	if checkpointResult.Result == StateFail || checkpointResult.Result == StateCannotVerify {
		applySetCheckpointResult(result, cp, checkpointResult.Result)
		return true
	}
	return applySetLinkChecks(result, cp, expected)
}

func applySetCheckpointResult(result *VerificationResult, cp SignedCheckpoint, state string) {
	// applySetCheckpointResult keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	result.Result = state
	result.SequenceState = worseState(result.SequenceState, state)
	if state == StateFail {
		result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s failed verification", cp.CheckpointID))
		return
	}
	result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s cannot verify", cp.CheckpointID))
}
