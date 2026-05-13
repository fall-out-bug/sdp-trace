package harnessobs

func lineIsolationRulePresent(path, pattern string) (bool, error) {
	// lineIsolationRulePresent keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	lines, err := readOptionalLines(path)
	if err != nil {
		return false, err
	}

	for _, line := range lines {
		if line == pattern {
			return true, nil
		}
	}
	return false, nil
}
