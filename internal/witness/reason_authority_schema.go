package witness

const (
	// ReasonMissingSigner records missing signer authority evidence.
	ReasonMissingSigner = "witness_signer_authority_missing"
	// ReasonMissingFreshness records missing signed freshness evidence.
	ReasonMissingFreshness = "witness_freshness_missing"
	// ReasonStaleFreshness records freshness evidence that has expired.
	ReasonStaleFreshness = "witness_freshness_stale"
	// ReasonSignerMismatch records signer identity or key mismatch evidence.
	ReasonSignerMismatch = "witness_signer_mismatch"
	// ReasonRevocationNA records that revocation evidence was not assessed.
	ReasonRevocationNA = "witness_revocation_not_assessed"
	// ReasonCertRevoked records a revoked certificate authority state.
	ReasonCertRevoked = "witness_certificate_revoked"
	// ReasonKeyCustodyNA records that key custody evidence was not assessed.
	ReasonKeyCustodyNA = "witness_key_custody_not_assessed"
)
