package main

func diffSets(a, b map[string]bool) []string {
	// Map keys are set members; false is unused so absent keys remain the only
	// stale/missing signal.
	var diff []string
	for k := range a {
		if !b[k] {
			diff = append(diff, k)
		}
	}
	return diff
}
