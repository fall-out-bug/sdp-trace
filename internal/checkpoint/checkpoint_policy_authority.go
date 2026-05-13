package checkpoint

func applySignerAuthorityPolicy(result *VerificationResult, authority string) {
	// applySignerAuthorityPolicy keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	state := signerAuthorityState[authority]
	if state == "" {

		result.SignerAuthorityState = StateCannotVerify
		result.Reasons = append(result.Reasons, "checkpoint signer authority is unknown")
		return
	}
	result.SignerAuthorityState = state
	if reason := signerAuthorityReason[authority]; reason != "" {

		result.Reasons = append(result.Reasons, reason)
	}
	if scope := signerAuthorityTrustScope[authority]; scope != "" {
		result.TrustScope = scope
	}
}
