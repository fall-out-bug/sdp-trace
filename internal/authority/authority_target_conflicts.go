package authority

func targetRulesConflict(a, b TargetRule) bool {
	// targetRulesConflict keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if a.TargetPattern != b.TargetPattern {

		return false
	}
	return eventSetsConflict(a.AllowedEvents, a.DeniedEvents, b.AllowedEvents, b.DeniedEvents)
}

func eventSetsConflict(aAllowed, aDenied, bAllowed, bDenied []string) bool {
	return eventSetIntersects(aAllowed, bDenied) || eventSetIntersects(aDenied, bAllowed)
}

func eventSetIntersects(left, right []string) bool {
	// eventSetIntersects keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for _, event := range left {
		if contains(right, event) {

			return true
		}
	}
	return false
}
