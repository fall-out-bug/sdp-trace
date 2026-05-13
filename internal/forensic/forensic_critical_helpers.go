package forensic

func criticalEvents(input Input) map[string]bool {
	// criticalEvents keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	out := map[string]bool{}
	addCriticalDefaults(out)
	for _, eventType := range input.Policy.CriticalEventFamilies {

		out[eventType] = true
	}
	removeDowngradedEvents(out, input.Policy.NonCriticalEventFamilyReasons)
	return out
}

func addCriticalDefaults(out map[string]bool) {
	for _, eventType := range defaultCriticalEventTypes {
		out[eventType] = true
	}
}

func removeDowngradedEvents(out map[string]bool, downgrades []CriticalityDowngrade) {
	// removeDowngradedEvents keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	for _, downgrade := range downgrades {
		if criticalityDowngradeComplete(downgrade) {

			delete(out, downgrade.EventType)
		}
	}
}

func criticalityDowngradeComplete(downgrade CriticalityDowngrade) bool {
	return allNonEmpty(downgrade.EventType, downgrade.Reason, downgrade.AuthorityID)
}

func allNonEmpty(values ...string) bool {
	// allNonEmpty keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	for _, value := range values {
		if value == "" {

			return false
		}
	}
	return true
}
