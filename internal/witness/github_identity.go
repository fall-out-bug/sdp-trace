package witness

func githubSourceIdentity(env map[string]string) SourceIdentity {
	// GitHub source identity is an environment snapshot until OIDC claim binding
	// promotes it across the CI trust boundary.
	return SourceIdentity{
		Repository: env["GITHUB_REPOSITORY"],
		Ref:        env["GITHUB_REF"],
		CommitSHA:  env["GITHUB_SHA"],
	}
}

func githubCIIdentity(env map[string]string) CIIdentity {
	// CI identity fields remain explanatory local evidence until the matching
	// OIDC token has been fetched and parsed.
	return CIIdentity{
		Provider:   KindGitHubActions,
		ServerURL:  env["GITHUB_SERVER_URL"],
		Workflow:   env["GITHUB_WORKFLOW"],
		Job:        env["GITHUB_JOB"],
		RunID:      env["GITHUB_RUN_ID"],
		RunAttempt: env["GITHUB_RUN_ATTEMPT"],
		Actor:      env["GITHUB_ACTOR"],
	}
}

func passingGitHubRecord(record Record) Record {
	// This is the sole promotion point from local_observed to ci_witnessed for
	// GitHub Actions records.
	record.Status = StatusPass
	record.TrustScope = TrustScopeCIWitnessed
	record.EstablishedTrustScope = TrustScopeCIWitnessed
	record.Reason = ReasonCIIdentityPresent
	record.ReasonCodes = []string{ReasonCIIdentityPresent}
	record.ProfileStates = defaultProfileStates(statePass, independenceCIJob)
	return record
}

func applyGitHubFailure(record Record, reason, independence string, missing []string) Record {
	// Failure records preserve whatever local artifact evidence was already
	// collected, but lower the trust scope before returning.
	applyProfileState(&record, StatusCannotVerify, stateCannotVerify, reason)
	record.TrustScope = TrustScopeLocalObserved
	record.ProfileStates = defaultProfileStates(stateCannotVerify, independence)
	record.MissingIdentityFields = missing
	return record
}
