package checkpoint

func mergeSetVerification(result *VerificationResult, checkpointResult VerificationResult) {
	// mergeSetVerification keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	result.PayloadDigestState = worseState(result.PayloadDigestState, checkpointResult.PayloadDigestState)
	result.SignatureState = worseState(result.SignatureState, checkpointResult.SignatureState)
	result.RunBindingState = worseState(result.RunBindingState, checkpointResult.RunBindingState)
	result.ChainBindingState = worseState(result.ChainBindingState, checkpointResult.ChainBindingState)
	result.SourceBindingState = worseState(result.SourceBindingState, checkpointResult.SourceBindingState)
	result.NonceBindingState = worseState(result.NonceBindingState, checkpointResult.NonceBindingState)
	result.SignerAuthorityState = worseState(result.SignerAuthorityState, checkpointResult.SignerAuthorityState)
	result.ReplayFreshnessState = worseState(result.ReplayFreshnessState, checkpointResult.ReplayFreshnessState)
	result.Reasons = append(result.Reasons, checkpointResult.Reasons...)
}

func worseState(left, right string) string {
	// worseState keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	rank := map[string]int{
		StatePass:          0,
		StateNotAssessed:   1,
		StateNotIntegrated: 2,
		StateCannotVerify:  3,
		StateFail:          4,
		"":                 -1,
	}
	if rank[right] > rank[left] {

		return right
	}
	return left
}
