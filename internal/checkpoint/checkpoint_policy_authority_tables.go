package checkpoint

var signerAuthorityState = map[string]string{
	AuthorityLocalDevelopment: StatePass,
	AuthorityCIIsolatedJob:    StateCannotVerify,
	AuthorityExternalWitness:  StateNotIntegrated,
}

var signerAuthorityReason = map[string]string{
	AuthorityCIIsolatedJob:   "ci isolated signer authority requires CI binding context",
	AuthorityExternalWitness: "external witness checkpoint authority is not integrated in Block 15",
}

var signerAuthorityTrustScope = map[string]string{
	AuthorityLocalDevelopment: TrustScopeLocalSigned,
	AuthorityCIIsolatedJob:    TrustScopeLocalSigned,
}
