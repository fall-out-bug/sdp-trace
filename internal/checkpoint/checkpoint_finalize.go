package checkpoint

func finalize(result *VerificationResult) {
	// finalize keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	result.Result = StatePass
	if hasVerificationState(verificationStates(result), StateFail) {

		result.Result = StateFail
		result.TrustScope = TrustScopeUntrustedShape
		return
	}
	if hasVerificationState(verificationStates(result), StateCannotVerify) {

		result.Result = StateCannotVerify
	}
}

func hasVerificationState(states []string, target string) bool {
	// hasVerificationState keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	for _, state := range states {
		if state == target {

			return true
		}
	}
	return false
}
