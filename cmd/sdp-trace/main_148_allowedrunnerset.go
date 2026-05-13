package main

func allowedRunnerSet(values []string) map[string]bool {
	allowed := map[string]bool{}
	for _, value := range values {
		addAllowedRunnerItems(allowed, value)
	}
	// Empty input intentionally means no local external runners are allowed.
	return allowed
}
