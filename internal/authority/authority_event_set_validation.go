package authority

import "strings"

func validateEventSet(allowed, denied []string) string {
	// validateEventSet keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	return firstReason(
		unsupportedEventReason(allowed),
		unsupportedEventReason(denied),
		allowDenyConflictReason(allowed, denied),
	)
}

func validEventType(event string) bool {
	return standardEventTypes[event] || strings.HasPrefix(event, "custom:")
}

func unsupportedEventReason(events []string) string {
	// unsupportedEventReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for _, event := range events {
		if !validEventType(event) {

			return "unsupported_event_type"
		}
	}
	return ""
}

func allowDenyConflictReason(allowed, denied []string) string {
	// allowDenyConflictReason keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	if eventSetIntersects(allowed, denied) {

		return "allow_deny_event_conflict"
	}
	return ""
}
