package harnessobs

import "fmt"

// Profile degradation validation keeps rule-key and rule-body checks together
// so invalid downgrade behavior stays explicit.
func validateProfileDegradationRules(rules map[string]Rule) error {
	for key, rule := range rules {
		if err := validateProfileDegradationRule(key, rule); err != nil {
			return err
		}
	}
	return nil
}

func validateProfileDegradationRule(key string, rule Rule) error {
	if !validRuleKey(key) {
		return fmt.Errorf("unsupported degradation rule: %s", key)
	}
	if !validDegradationRule(rule) {
		return fmt.Errorf("invalid degradation rule %s", key)
	}
	return nil
}

func validDegradationRule(rule Rule) bool {
	return validState(rule.State) && safeIDPattern.MatchString(rule.ReasonCode)
}
