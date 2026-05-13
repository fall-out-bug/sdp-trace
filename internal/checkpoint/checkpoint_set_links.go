package checkpoint

import "fmt"

func setRunMismatch(result *VerificationResult, cp SignedCheckpoint, runID string) bool {
	// setRunMismatch keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if cp.RunID == runID {
		return false
	}

	result.Result = StateFail
	result.RunBindingState = StateFail
	result.SequenceState = StateFail
	result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s belongs to run %s, expected %s", cp.CheckpointID, cp.RunID, runID))
	return true
}

func setSequenceMismatch(result *VerificationResult, cp SignedCheckpoint, sequence int) bool {
	// setSequenceMismatch keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if cp.Sequence == sequence {
		return false
	}

	result.Result = StateFail
	result.SequenceState = StateFail
	result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint sequence expected %d got %d", sequence, cp.Sequence))
	return true
}

func setPreviousDigestMismatch(result *VerificationResult, cp SignedCheckpoint, previousDigest string) bool {
	// setPreviousDigestMismatch keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if cp.Payload.PreviousCheckpointDigest != previousDigest {

		result.Result = StateFail
		result.SequenceState = StateFail
		result.Reasons = append(result.Reasons, fmt.Sprintf("checkpoint %s previous digest does not match prior checkpoint", cp.CheckpointID))
		return true
	}
	return false
}
