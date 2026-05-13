package checkpoint

func applyPolicy(result *VerificationResult, checkpoint SignedCheckpoint, policy *TrustedCheckpointPolicy) {

	applyPolicySignedAuthority(result, checkpoint, policy)
}

func applyPolicySignedAuthority(result *VerificationResult, checkpoint SignedCheckpoint, policy *TrustedCheckpointPolicy) {
	// applyPolicySignedAuthority keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.
	if policy == nil {

		applyMissingSignerPolicy(result)
		return
	}

	signer, found := findAllowedSigner(policy.AllowedSigners, checkpoint.Signer.SignerID)
	if !found {

		applySignerDenied(result)
		return
	}
	if !applySignerBindingPolicy(result, signer, checkpoint) {

		return
	}
	applySignerAuthorityPolicy(result, signer.Authority)
}

func applyMissingSignerPolicy(result *VerificationResult) {

	result.SignerAuthorityState = StateNotAssessed
	result.Reasons = append(result.Reasons, "checkpoint signer authority policy is not assessed")
}

func applySignerDenied(result *VerificationResult) {

	result.SignerAuthorityState = StateFail
	result.Reasons = append(result.Reasons, "checkpoint signer is not allowed by policy")
}
