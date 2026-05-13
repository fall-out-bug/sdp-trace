package ciartifact

func sanitizeSource(input SourceIdentity) (SourceIdentity, bool) {
	// sanitizeSource keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	repository, repositoryOK := sanitizeSourceField(input.Repository, func(value string) bool { return safeIdentityToken(value, "/._-") })
	ref, refOK := sanitizeSourceField(input.Ref, safeRef)
	commitSHA, commitOK := sanitizeSourceField(input.CommitSHA, safeCommitSHA)
	return SourceIdentity{Repository: repository, Ref: ref, CommitSHA: commitSHA}, allTrue(repositoryOK, refOK, commitOK)
}

func sanitizeSourceField(value string, valid func(string) bool) (string, bool) {
	// sanitizeSourceField keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if value == "" || valid(value) {
		return value, true
	}

	return "", false
}

func safeCommitSHA(value string) bool {
	return safeHex(value, 40) || safeHex(value, 64)
}
