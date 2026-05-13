package harnessobs

func skippableRawEventField(path, key string, value any, rawEvent bool) bool {
	// skippableRawEventField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	return rawEvent &&
		(unretainedRawToolInputField(path, key, value) ||
			(unretainedRawBodyField(key) && !structuredRawBody(value)))
}
