package witness

func handleGitHubIdentityChecks(record Record, env map[string]string) (Record, bool) {
	// Missing provider identity blocks the trust upgrade before any OIDC claim
	// is considered.
	missing := missingGitHubIdentity(env)
	if len(missing) > 0 {
		return applyGitHubFailure(record, ReasonMissingCIIdentity, independenceSameJob, missing), true
	}
	oidcMissing := missingGitHubOIDC(env)
	if len(oidcMissing) > 0 {
		return applyGitHubFailure(record, ReasonMissingCIOIDC, independenceCIJob, oidcMissing), true
	}
	return record, false
}

func handleGitHubOIDCChecks(record Record, env map[string]string, fetcher TokenFetcher) (Record, bool) {
	// Fetch and claim matching are one trust boundary: failures preserve local
	// evidence but cannot establish CI-witnessed scope.
	token, err := fetcher(env)
	if err != nil {
		return applyGitHubFailure(record, ReasonInvalidCIOIDC, independenceCIJob, nil), true
	}
	claims, err := parseOIDCClaims(token)
	if err != nil || !claimsMatchEnvironment(claims, env) {
		return applyGitHubFailure(record, ReasonInvalidCIOIDC, independenceCIJob, nil), true
	}
	record.OIDC = &claims
	return record, false
}
