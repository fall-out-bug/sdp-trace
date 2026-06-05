package main

import (
	"fmt"
	"sort"
	"strings"
)

// Command-surface usage diff helpers grouped from numbered shards.
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

func sortedUsageDiffs(registry, help map[string]bool) (missing, stale []string) {
	// Stable sort keeps test failures deterministic and makes drift output
	// reviewable as a copied command list.
	missing = diffSets(registry, help)
	stale = diffSets(help, registry)
	sort.Strings(missing)
	sort.Strings(stale)
	return missing, stale
}

func commandSurfaceDriftError(missing, stale []string) error {
	// Empty diff means the registry and help text agree, so tests should return
	// nil instead of formatting an empty diagnostic.
	parts := commandSurfaceDriftParts(missing, stale)
	if len(parts) == 0 {
		return nil
	}
	return fmt.Errorf("command surface drift: %s", strings.Join(parts, " | "))
}

func commandSurfaceDriftParts(missing, stale []string) []string {
	// Keep missing and stale diagnostics separate so reviewers know whether to
	// update usageText, the registry, or both.
	var parts []string
	if len(missing) > 0 {
		parts = append(parts, fmt.Sprintf("missing from usageText: %s", strings.Join(missing, "; ")))
	}
	if len(stale) > 0 {
		parts = append(parts, fmt.Sprintf("stale in usageText: %s", strings.Join(stale, "; ")))
	}
	return parts
}
