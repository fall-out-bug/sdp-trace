package repoobserver

func Doctor(opts Options) (Status, error) {
	// Doctor is read-only: it reports local structural setup without installing
	// generated observer files.
	opts, err := normalizeOptions(opts)
	if err != nil {
		return Status{}, err
	}
	opts, err = withConfiguredRepositoryID(opts)
	if err != nil {
		return Status{}, err
	}
	return buildStatus(opts, false)
}
