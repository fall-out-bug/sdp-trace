package harnessobs

func unsafeEncodedToken(path, value string, rawEvent bool) bool {
	// unsafeEncodedToken keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if safeEncodedTokenExemption(path, value, rawEvent) {
		return false
	}

	return base64TokenPattern.MatchString(value)
}
