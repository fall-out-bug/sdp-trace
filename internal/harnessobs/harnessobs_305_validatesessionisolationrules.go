package harnessobs

func validateSessionIsolationRules(rules []SessionIsolationRule) error {
	// validateSessionIsolationRules keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, rule := range rules {

		if err := validateSessionIsolationRule(rule); err != nil {
			return err
		}
	}
	return nil
}
