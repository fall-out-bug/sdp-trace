package witness

func BuildGitHubActions(runsRoot, reportDir string, env map[string]string) (Record, error) {
	return BuildGitHubActionsWithFetcher(runsRoot, reportDir, env, FetchGitHubOIDCToken)
}

func BuildGitHubActionsWithFetcher(runsRoot, reportDir string, env map[string]string, fetcher TokenFetcher) (Record, error) {
	// GitHub Actions witness generation starts as a local record and upgrades
	// only after artifact hashing, required environment fields, and OIDC claims
	// all bind to the same execution.
	record := baseRecord(KindGitHubActions)
	record.Source = githubSourceIdentity(env)
	record.CI = githubCIIdentity(env)

	if err := hydrateGitHubArtifacts(&record, runsRoot, reportDir); err != nil {
		return Record{}, err
	}
	if record, blocked := handleGitHubIdentityChecks(record, env); blocked {
		return record, nil
	}
	if record, blocked := handleGitHubOIDCChecks(record, env, fetcher); blocked {
		return record, nil
	}
	return passingGitHubRecord(record), nil
}
