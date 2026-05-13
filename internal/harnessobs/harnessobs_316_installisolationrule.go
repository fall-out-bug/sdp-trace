package harnessobs

func installIsolationRule(rule SessionIsolationRule) (SessionIsolationResult, error) {
	// installIsolationRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := ensureIsolationRule(rule); err != nil {
		return SessionIsolationResult{}, err
	}

	return verifyIsolationRule(rule)
}
