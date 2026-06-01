package harnessobs

import "sort"

func sortedFamilies(families map[string]bool) []string {
	ordered := make([]string, 0, len(families))
	for family := range families {
		ordered = append(ordered, family)
	}
	sort.Strings(ordered)
	return ordered
}
