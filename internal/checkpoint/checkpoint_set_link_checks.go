package checkpoint

func applySetLinkChecks(result *VerificationResult, cp SignedCheckpoint, expected setLinkExpectation) bool {
	// applySetLinkChecks keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	if setRunMismatch(result, cp, expected.runID) {
		return true
	}
	if setSequenceMismatch(result, cp, expected.sequence) {
		return true
	}
	return setPreviousDigestMismatch(result, cp, expected.previousDigest)
}
