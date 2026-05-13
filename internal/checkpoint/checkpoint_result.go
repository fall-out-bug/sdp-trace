package checkpoint

func baseResult(checkpoint SignedCheckpoint) VerificationResult {
	// baseResult keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	result := VerificationResult{
		SchemaVersion: VerificationSchemaVersion,
		CheckpointID:  checkpoint.CheckpointID,
		RunID:         checkpoint.RunID,
		Result:        StatePass,
		TrustScope:    TrustScopeLocalSigned,
		Reasons:       []string{},
	}
	applyBaseEvidenceStates(&result)
	return result
}

func applyBaseEvidenceStates(result *VerificationResult) {
	// applyBaseEvidenceStates keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	result.PayloadDigestState = StateNotAssessed
	result.SignatureState = StateNotAssessed
	result.RunBindingState = StatePass
	result.ChainBindingState = StatePass
	result.SourceBindingState = StatePass
	result.NonceBindingState = StatePass
	result.SequenceState = StatePass
	result.SignerAuthorityState = StateNotAssessed
	result.ReplayFreshnessState = StateNotAssessed
}

func verificationStates(result *VerificationResult) []string {
	// verificationStates keeps checkpoint proof evidence explicit and replay-bound.
	// Shape, digest, signature, chain, nonce, source, and signer states stay separate.
	// This helper does not upgrade local checkpoint data into external trust.

	return []string{
		result.PayloadDigestState,
		result.SignatureState,
		result.RunBindingState,
		result.ChainBindingState,
		result.SourceBindingState,
		result.NonceBindingState,
		result.SequenceState,
		result.SignerAuthorityState,
	}
}
