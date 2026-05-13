package harnessobs

func installIsolationRules(profilePath string, rules []SessionIsolationRule) ([]SessionIsolationResult, error) {
	// installIsolationRules keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	results := make([]SessionIsolationResult, 0, len(rules))
	for _, rule := range rules {

		resolvedRule, err := resolveIsolationRuleTarget(profilePath, rule)
		if err != nil {
			return nil, err
		}

		result, err := installIsolationRule(resolvedRule)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
