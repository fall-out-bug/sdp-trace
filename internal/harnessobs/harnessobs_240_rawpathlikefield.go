package harnessobs

func rawPathLikeField(path string) bool {
	// rawPathLikeField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch rawPathFieldName(path) {
	case "path", "file", "filepath", "file_path", "dir", "directory", "cwd":

		return true
	default:
		return false
	}
}
