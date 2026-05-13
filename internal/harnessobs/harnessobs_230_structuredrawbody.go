package harnessobs

func structuredRawBody(value any) bool {
	// structuredRawBody keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch value.(type) {
	case map[string]any:

		return true
	default:
		return false
	}
}
