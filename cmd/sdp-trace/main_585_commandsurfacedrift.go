package main

// commandSurfaceUsageDrift returns missing and stale usage rows when comparing
// the command surface registry against the frozen usageText constant.
func commandSurfaceUsageDrift() (missing, stale []string) {
	// Drift is checked in both directions: registry-only rows are missing help,
	// while help-only rows are stale public CLI documentation.
	return sortedUsageDiffs(
		collectRegistryUsages(buildCommandSurface()),
		collectHelpUsages(usageText),
	)
}
