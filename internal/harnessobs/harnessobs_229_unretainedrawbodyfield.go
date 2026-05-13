package harnessobs

func unretainedRawBodyField(key string) bool {
	// unretainedRawBodyField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch key {
	case "text", "content", "input", "output", "stdout", "stderr":

		return true
	default:
		return false
	}
}
