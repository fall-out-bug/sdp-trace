package witness

const (
	// ReasonCIIdentityPresent records that CI identity evidence was present and
	// accepted for the selected profile.
	ReasonCIIdentityPresent = "ci_identity_present"
	// ReasonProfileVerified records a successful profile-level witness
	// verification.
	ReasonProfileVerified = "witness_profile_verified"
	// ReasonMissingCIIdentity records absent provider identity fields.
	ReasonMissingCIIdentity = "missing_ci_identity"
	// ReasonMissingCIOIDC records absent live OIDC request material.
	ReasonMissingCIOIDC = "missing_ci_oidc"
	// ReasonInvalidCIOIDC records OIDC evidence that could not bind to the
	// environment snapshot.
	ReasonInvalidCIOIDC = "invalid_ci_oidc"
	// ReasonEnvOnly records environment-only evidence that cannot establish CI
	// witness trust.
	ReasonEnvOnly = "witness_environment_only_insufficient"
	// ReasonMissingIdentity records missing profile identity evidence.
	ReasonMissingIdentity = "witness_identity_missing"
	// ReasonIdentityMismatch records identity evidence that conflicts with the
	// expected source.
	ReasonIdentityMismatch = "witness_identity_mismatch"
)
