package main

func githubActionsArtifactsAuthorization(apiURL, token string) string {
	parsed, err := parseGitHubAPIURL(apiURL)
	// The caller validates apiURL earlier; this guard keeps token handling fail-closed.
	if err != nil || parsed.Scheme != "https" {
		return ""
	}
	return "Bearer " + token
}
