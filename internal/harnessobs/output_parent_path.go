package harnessobs

func safeParentDir(path string) (string, error) {
	// Parent directories may be missing; containment is checked against the
	// closest existing parent after resolving symlinks.
	clean, err := normalizePotentialParentPath(path)
	if err != nil {
		return "", err
	}

	resolved, err := resolveParentPathWithinWorkingDirectory(clean)
	if err != nil {
		return "", err
	}

	return ensurePathInsideWorkingDirectory(resolved)
}
