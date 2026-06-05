package main

import "strings"

// Command-surface usage collection helpers grouped from numbered shards.
func addCommandSurfaceUsages(usages map[string]bool, cmd commandSurfaceCmd) {
	// Variations carry documented command lines just like Usage, but blank Usage
	// fields are allowed for grouped command families.
	addCommandSurfaceUsage(usages, cmd.Usage)
	for _, variation := range cmd.Variations {
		addCommandSurfaceUsage(usages, variation)
	}
}

func addCommandSurfaceUsage(usages map[string]bool, usage string) {
	if usage != "" {
		usages[usage] = true
	}
}

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

func collectRegistryUsages(surface commandSurfaceSchema) map[string]bool {
	// The registry owns both canonical Usage rows and example Variations, so
	// drift detection must treat both as copy-pasteable public command lines.
	usages := make(map[string]bool)
	for _, cmd := range surface.Commands {
		addCommandSurfaceUsages(usages, cmd)
	}
	return usages
}

func collectHelpUsages(help string) map[string]bool {
	// Only indented sdp-trace usage rows participate; prose and section labels
	// remain human guidance rather than machine-readable command contracts.
	usages := make(map[string]bool)
	for _, line := range strings.Split(help, "\n") {
		if strings.HasPrefix(line, "  sdp-trace ") {
			usages[strings.TrimSpace(line)] = true
		}
	}
	return usages
}
