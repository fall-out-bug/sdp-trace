package checkpoint

func applySignerBindingPolicy(result *VerificationResult, signer TrustedSigner, checkpoint SignedCheckpoint) bool {
	// applySignerBindingPolicy keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if signer.PublicKey == "" {

		result.SignerAuthorityState = StateCannotVerify
		result.Reasons = append(result.Reasons, "checkpoint signer policy missing public key binding")
		return false
	}
	if signer.PublicKey != checkpoint.Signature.PublicKey {

		result.SignerAuthorityState = StateFail
		result.Reasons = append(result.Reasons, "checkpoint signer public key does not match policy")
		return false
	}
	if signer.Authority != checkpoint.Signer.Authority {

		result.SignerAuthorityState = StateFail
		result.Reasons = append(result.Reasons, "checkpoint signer authority does not match policy")
		return false
	}
	return true
}

func findAllowedSigner(signers []TrustedSigner, signerID string) (TrustedSigner, bool) {
	// findAllowedSigner keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	for _, signer := range signers {
		if signer.SignerID == signerID {

			return signer, true
		}
	}
	return TrustedSigner{}, false
}
