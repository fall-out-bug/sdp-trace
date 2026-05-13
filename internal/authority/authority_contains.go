package authority

func contains(values []string, needle string) bool {
	// contains keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for _, value := range values {
		if value == needle {

			return true
		}
	}
	return false
}
