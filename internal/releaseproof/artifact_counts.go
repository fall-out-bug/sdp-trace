package releaseproof

func artifactCounts(repoRoot, sourceCommit string, artifacts []ManifestArtifact) (ArtifactCounts, []ArtifactIssue) {
	// Count every manifest artifact against the immutable source commit; each
	// issue keeps the manifest path so reports stay source-bound and auditable.
	// Checked counts manifest obligations, not successful reads, so missing
	// artifacts are visible in both the denominator and issue list.
	counts := ArtifactCounts{Checked: len(artifacts)}
	// Issues are sparse: matched artifacts stay represented by Checked, while
	// only missing or mismatched obligations get explicit rows.
	issues := []ArtifactIssue{}
	for _, artifact := range artifacts {
		issue, ok := artifactIssue(repoRoot, sourceCommit, artifact)
		if !ok {
			continue
		}
		countArtifactIssue(&counts, issue)
		issues = append(issues, issue)
	}
	return counts, issues
}

func countArtifactIssue(counts *ArtifactCounts, issue ArtifactIssue) {
	// Only replay failures change the aggregate counters; matched artifacts are
	// already represented by Checked.
	switch issue.Issue {
	case StatusMissing:
		counts.Missing++
	case StatusMismatch:
		counts.Mismatched++
	}
}
