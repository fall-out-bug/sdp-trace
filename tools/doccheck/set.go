package main

import "sort"

func missingCommands(expected, actual []string) []string {
	// Preserve expected ordering in the reported drift so failures are stable
	// and easy to copy back into the docs.
	actualSet := map[string]bool{}
	for _, command := range actual {
		actualSet[command] = true
	}
	var missing []string
	for _, command := range expected {
		if !actualSet[command] {
			missing = append(missing, command)
		}
	}
	return missing
}

func uniqueSorted(values []string) []string {
	// Sorting removes accidental markdown/help ordering differences while still
	// requiring exact command text matches.
	seen := map[string]bool{}
	unique := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		unique = append(unique, value)
	}
	sort.Strings(unique)
	return unique
}
