package releaseproof

func applySourceEvidence(result *Verification, state verificationState) {
	// Source evidence records only what was checked against the manifest's
	// source commit; it does not make external production claims.
	result.ArtifactDigestStatus = state.artifactStatus
	result.SourceCommit = state.sourceCommit
	result.SourceCommitStatus = state.commitStatus
	result.SourceCommitArtifactStatus = state.artifactStatus
	result.SourceCommitArtifactCounts = state.artifactCounts
	result.SourceCommitReason = state.sourceReason
	result.ArtifactIssues = state.artifactIssues
}
