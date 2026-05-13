package harnessobs

func ensureJSONReadDenyRule(path, pattern string) error {
	// ensureJSONReadDenyRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	config, err := readOptionalJSONObject(path)
	if err != nil {
		return err
	}
	setJSONReadDeny(config, pattern)
	return writeJSON(path, config)
}
