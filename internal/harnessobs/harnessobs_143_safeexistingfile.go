package harnessobs

func safeExistingFile(path string) (string, error) {
	// safeExistingFile keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return safeExistingPath(path, existingPathSpec{
		traversalError: "path must be relative local file without traversal",
		requireDir:     false,
		typeError:      "path must be a file",
	})
}
