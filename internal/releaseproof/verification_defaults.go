package releaseproof

func applyExternalTrustDefaults(result *Verification) {
	// External production controls are intentionally outside this local source
	// profile, so each related field stays explicitly not_assessed.
	result.ExternalTrustProfile = StatusNotAssessed
	result.ExternalAttestationRef = nil
	result.TransparencyEvidenceRef = nil
	result.SourceURIStatus = StatusNotAssessed
	result.ProtectedRefStatus = StatusNotAssessed
	result.WorkflowIdentityStatus = StatusNotAssessed
	result.ApprovalStatus = StatusNotAssessed
}

func applyReleaseTrustDefaults(result *Verification) {
	// Local source proof is not equivalent to release authorization.
	result.ProductionReleaseVerified = productionReleaseNotAssessed()
	result.TransparencyLogStatus = StatusNotAssessed
	result.FreshnessStatus = StatusNotAssessed
	result.TrustedContractRelease = false
	result.PrivateEquivalentProfile = "not_assessed"
	result.ExternalTrustReason = "external production trust is not assessed by the local source-bound profile"
}

func productionReleaseNotAssessed() ProofStateBoolean {
	// The boolean value is absent because this profile does not inspect
	// production attestation evidence.
	return ProofStateBoolean{
		State:  StatusNotAssessed,
		Value:  nil,
		Reason: "Production release verification requires external attestation in addition to source-bound local checks.",
	}
}
