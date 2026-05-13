package harnessobs

func safeCleanOutDir(clean string) (string, error) {
	// safeCleanOutDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	exists, err := pathExistsForLstat(clean)
	if err != nil {
		return "", err
	}
	if exists {
		return safeExistingOutDir(clean)
	}

	if err := ensureOutParentInsideWorkingDirectory(clean); err != nil {
		return "", err
	}
	return ensureOutDirEmptyOrMissing(clean)
}
