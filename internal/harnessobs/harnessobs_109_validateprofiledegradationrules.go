package harnessobs

import (
	"fmt"
)

func validateProfileDegradationRules(rules map[string]Rule) error {
	// validateProfileDegradationRules keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for key, rule := range rules {

		if !validRuleKey(key) {
			return fmt.Errorf("unsupported degradation rule: %s", key)
		}
		if !validDegradationRule(rule) {
			return fmt.Errorf("invalid degradation rule %s", key)
		}
	}
	return nil
}
