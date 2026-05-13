package harnessobs

func setJSONReadDeny(config map[string]any, pattern string) {
	// setJSONReadDeny keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	permission := ensureObject(config, "permission")
	read := ensureObject(permission, "read")

	read[pattern] = "deny"
}
