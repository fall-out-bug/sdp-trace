package harnessobs

func childPath(parent, key string) string {
	// childPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if parent == "" {

		return key
	}
	return parent + "." + key
}
