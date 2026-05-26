package main

import (
	"fmt"
	"strings"
)

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
