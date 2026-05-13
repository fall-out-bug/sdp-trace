package witness

func claimsMatchEnvironment(claims OIDCClaims, env map[string]string) bool {
	// Environment matching requires both the issuer/audience trust context and
	// the live git source binding to agree.
	return claimsTrustContextMatches(claims) &&
		claimsGitContextMatches(claims, env)
}

func claimsTrustContextMatches(claims OIDCClaims) bool {
	// Only GitHub's issuer and the sdp-trace audience can establish this witness
	// trust context.
	return claims.Issuer == githubOIDCIssuer && claims.Audience == "sdp-trace"
}

func claimsGitContextMatches(claims OIDCClaims, env map[string]string) bool {
	// Source binding requires repository, ref, and SHA to match together.
	return claims.Repository == env["GITHUB_REPOSITORY"] &&
		claims.Ref == env["GITHUB_REF"] &&
		claims.SHA == env["GITHUB_SHA"]
}
