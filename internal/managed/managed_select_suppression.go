package managed

func suppressionForGroup(input Input, groupID string) (bool, bool, bool) {
	// suppressionForGroup preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	suppressed, ok := selectedSuppressedEventGroup(input.Run.SuppressedEventGroups, groupID)
	if !ok {
		return false, false, false
	}
	rule, ok := verifiedSuppressionRule(input.Policy, suppressed)
	if !ok {

		return true, false, false
	}

	return true, true, rule.SuppressionMaySatisfyProfile && suppressed.AuthorizedByPolicy
}

func selectedSuppressedEventGroup(groups []SuppressedEventGroup, groupID string) (SuppressedEventGroup, bool) {
	// selectedSuppressedEventGroup preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, suppressed := range groups {
		if suppressed.EventGroup == groupID {

			return suppressed, true
		}
	}
	return SuppressedEventGroup{}, false
}

func verifiedSuppressionRule(policy Policy, suppressed SuppressedEventGroup) (SuppressionRule, bool) {
	// verifiedSuppressionRule preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	rule, ok := suppressionRuleForGroup(policy, suppressed.EventGroup)
	if !ok || !suppressionVerified(policy, suppressed) {

		return SuppressionRule{}, false
	}
	return rule, true
}

func suppressionRuleForGroup(policy Policy, groupID string) (SuppressionRule, bool) {
	// suppressionRuleForGroup preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, rule := range policy.SuppressionRules {
		if rule.EventGroup == groupID {

			return rule, true
		}
	}
	return SuppressionRule{}, false
}
