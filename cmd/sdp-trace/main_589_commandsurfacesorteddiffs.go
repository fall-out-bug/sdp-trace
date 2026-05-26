package main

import "sort"

func sortedUsageDiffs(registry, help map[string]bool) (missing, stale []string) {
	// Stable sort keeps test failures deterministic and makes drift output
	// reviewable as a copied command list.
	missing = diffSets(registry, help)
	stale = diffSets(help, registry)
	sort.Strings(missing)
	sort.Strings(stale)
	return missing, stale
}
