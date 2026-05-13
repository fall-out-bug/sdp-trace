package harnessobs

func safeParentDir(path string) (string, error) {
	// safeParentDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

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
