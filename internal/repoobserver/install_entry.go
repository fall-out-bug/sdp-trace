package repoobserver

func Install(opts Options) (Status, error) {
	// Install builds a preview first so --write callers can still receive the
	// same surface model after mutations are applied.
	opts, err := normalizeOptions(opts)
	if err != nil {
		return Status{}, err
	}
	status, err := buildStatus(opts, true)
	if err != nil {
		return status, err
	}
	if !opts.Write {
		return status, nil
	}
	return installWriteMode(opts, status)
}

func installWriteMode(opts Options, status Status) (Status, error) {
	// Preserve mutation summaries across the post-write rescan so the final
	// status shows both resulting surfaces and what changed.
	summary, err := writeInstallFiles(opts)
	status.ForceDiffSummary = summary
	if err != nil {
		return status, err
	}
	status, err = buildStatus(opts, false)
	status.ForceDiffSummary = summary
	return status, err
}
