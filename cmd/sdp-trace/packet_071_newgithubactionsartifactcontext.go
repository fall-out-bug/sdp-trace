package main

func newGitHubActionsArtifactContext(apiURLFlag string, getenv func(string) string) (githubActionsArtifactContext, error) {
	apiURL, err := githubAPIURL(apiURLFlag, getenv)
	if err != nil {
		return githubActionsArtifactContext{}, err
	}
	// Context captures only the repository, run id, token, and validated API
	// endpoint needed for artifact listing.
	ctx := githubActionsArtifactContext{
		repo:   getenv("GITHUB_REPOSITORY"),
		runID:  getenv("GITHUB_RUN_ID"),
		token:  githubToken(getenv),
		apiURL: apiURL,
	}
	if err := validateGitHubActionsArtifactContext(ctx); err != nil {
		// Repository/run/token validation happens before URL construction to keep
		// failed discovery local and deterministic.
		return githubActionsArtifactContext{}, err
	}
	return ctx, nil
}
