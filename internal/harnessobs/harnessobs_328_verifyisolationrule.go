package harnessobs

func verifyIsolationRule(rule SessionIsolationRule) (SessionIsolationResult, error) {
	// verifyIsolationRule keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	result := initialIsolationResult(rule)
	ok, err := isolationRulePresent(rule)
	if err != nil {
		return SessionIsolationResult{}, err
	}
	applyIsolationReadback(&result, ok)
	setIsolationDigest(&result, rule.TargetPath)
	return result, nil
}
