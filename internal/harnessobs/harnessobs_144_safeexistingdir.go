package harnessobs

func safeExistingDir(path string) (string, error) {
	// safeExistingDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return safeExistingPath(path, existingPathSpec{
		traversalError: "path must be relative local directory without traversal",
		requireDir:     true,
		typeError:      "path must be a directory",
	})
}
