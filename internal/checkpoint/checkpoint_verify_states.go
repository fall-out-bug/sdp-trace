package checkpoint

func applyPayloadDigestState(result *VerificationResult, checkpoint SignedCheckpoint) {
	// applyPayloadDigestState keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	digest, ok := verifyPayloadDigest(checkpoint)
	if !ok {

		result.PayloadDigestState = StateCannotVerify
		result.Reasons = append(result.Reasons, "checkpoint payload cannot be canonicalized")
		return
	}
	if checkpoint.PayloadDigest == digest {

		result.PayloadDigestState = StatePass
		return
	}
	result.PayloadDigestState = StateFail
	result.Reasons = append(result.Reasons, "checkpoint payload_digest does not match canonical payload")
}

func applySignatureState(result *VerificationResult, checkpoint SignedCheckpoint) {
	// applySignatureState keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if verifySignature(checkpoint) {

		result.SignatureState = StatePass
		return
	}
	result.SignatureState = StateFail
	result.Reasons = append(result.Reasons, "checkpoint signature is invalid")
}
