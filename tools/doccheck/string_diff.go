package main

func diffStringSliceAgainstSet(slice []string, set map[string]bool, skip string) []string {
	var diff []string
	for _, v := range slice {
		if v == skip {
			continue
		}
		if !set[v] {
			diff = append(diff, v)
		}
	}
	return diff
}
