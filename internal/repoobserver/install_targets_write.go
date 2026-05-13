package repoobserver

func writeInstallTargets(opts Options) ([]DiffSummary, error) {
	// Generated config and generated docs share one repository ID for internal
	// consistency.
	// Summaries accumulate only safe path/action/digest facts from each target.
	repoID := opts.RepositoryID
	if repoID == "" {
		repoID = derivedRepositoryID(opts.RepoRoot)
	}
	summaries := make([]DiffSummary, 0)
	for _, target := range installTargets(opts, repoID) {
		summary, err := writeTarget(opts, target)
		if err != nil {
			return summaries, err
		}
		summaries = append(summaries, summary...)
	}
	return summaries, nil
}
