package main

func diffStringSliceAgainstSet(slice []string, set map[string]bool, skip string) []string {
	var diff []string
	for _, v := range slice {
		if v == skip {
			// A caller-provided skip row is an allowed sentinel, not drift.
			continue
		}
		if !set[v] {
			// Preserve original slice order so drift output follows source docs.
			diff = append(diff, v)
		}
	}
	return diff
}
