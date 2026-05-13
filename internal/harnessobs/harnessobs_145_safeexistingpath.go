package harnessobs

func safeExistingPath(path string, spec existingPathSpec) (string, error) {
	// safeExistingPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	cleanPath, err := sanitizeExistingPath(path, spec.traversalError)
	if err != nil {
		return "", err
	}

	abs, err := resolveExistingAbsolutePath(cleanPath)
	if err != nil {
		return "", err
	}
	rel, err := relativeWorkingDirectoryPath(abs)
	if err != nil {
		return "", err
	}
	return rel, ensureExpectedPathType(abs, spec)
}
