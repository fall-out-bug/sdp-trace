package checkpoint

func Verify(runDir string, checkpoint SignedCheckpoint, policy *TrustedCheckpointPolicy) VerificationResult {
	// Verify keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	result := baseResult(checkpoint)
	if !applyCheckpointShape(&result, checkpoint) {
		return result
	}

	expected, err := BuildPayload(runDir, checkpoint.Payload.PreviousCheckpointDigest)
	if err != nil {

		cannotVerify(&result, err.Error())
		return result
	}
	applyPayloadDigestState(&result, checkpoint)
	applySignatureState(&result, checkpoint)
	compareBindings(&result, expected, checkpoint.Payload)
	applyPolicy(&result, checkpoint, policy)
	finalize(&result)
	return result
}

func applyCheckpointShape(result *VerificationResult, checkpoint SignedCheckpoint) bool {
	// applyCheckpointShape keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if err := validateEnvelope(checkpoint); err != nil {

		failShape(result, err.Error())
		return false
	}
	if err := validateSequenceLink(checkpoint.Sequence, checkpoint.Payload.PreviousCheckpointDigest); err != nil {

		result.SequenceState = StateFail
		failShape(result, err.Error())
		return false
	}
	return true
}

func failShape(result *VerificationResult, reason string) {
	// failShape keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	result.Result = StateFail
	result.TrustScope = TrustScopeUntrustedShape
	result.Reasons = append(result.Reasons, reason)
}

func cannotVerify(result *VerificationResult, reason string) {

	result.Result = StateCannotVerify
	result.Reasons = append(result.Reasons, reason)
}
