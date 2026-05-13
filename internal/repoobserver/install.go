package repoobserver

func writeInstallFiles(opts Options) ([]DiffSummary, error) {
	// Validate hooksPath before writing any files because changing it affects
	// all local git hook execution.
	// Repository files are written before hooksPath is changed so a failed file
	// write does not partially redirect local hooks.
	if err := ensureNoUnsafeHooksPath(opts); err != nil {
		return nil, err
	}
	summaries, err := writeInstallTargets(opts)
	if err != nil {
		return summaries, err
	}
	summary, err := updateGitignore(opts)
	if err != nil {
		return summaries, err
	}
	summaries = append(summaries, summary...)
	return appendHooksPathSummary(opts, summaries)
}
