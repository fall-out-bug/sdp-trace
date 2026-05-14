package main

import (
	"fmt"
	"sort"
	"strings"
)

// commandSurfaceUsageDrift returns missing and stale usage rows when comparing
// the command surface registry against the frozen usageText constant.
func commandSurfaceUsageDrift() (missing, stale []string, err error) {
	registryUsages := collectRegistryUsages(buildCommandSurface())
	helpUsages := collectHelpUsages(usageText)
	missing = diffSets(registryUsages, helpUsages)
	stale = diffSets(helpUsages, registryUsages)
	sort.Strings(missing)
	sort.Strings(stale)
	return missing, stale, nil
}

func collectRegistryUsages(surface commandSurfaceSchema) map[string]bool {
	usages := make(map[string]bool)
	for _, cmd := range surface.Commands {
		addUsage(usages, cmd.Usage)
		addUsages(usages, cmd.Variations)
	}
	return usages
}

func addUsage(usages map[string]bool, usage string) {
	if usage != "" {
		usages[usage] = true
	}
}

func addUsages(usages map[string]bool, values []string) {
	for _, v := range values {
		addUsage(usages, v)
	}
}

func collectHelpUsages(help string) map[string]bool {
	usages := make(map[string]bool)
	for _, line := range strings.Split(help, "\n") {
		if strings.HasPrefix(line, "  sdp-trace ") {
			usages[strings.TrimSpace(line)] = true
		}
	}
	return usages
}

func diffSets(a, b map[string]bool) []string {
	var diff []string
	for k := range a {
		if !b[k] {
			diff = append(diff, k)
		}
	}
	return diff
}

func commandSurfaceDriftError(missing, stale []string) error {
	var parts []string
	if len(missing) > 0 {
		parts = append(parts, fmt.Sprintf("missing from usageText: %s", strings.Join(missing, "; ")))
	}
	if len(stale) > 0 {
		parts = append(parts, fmt.Sprintf("stale in usageText: %s", strings.Join(stale, "; ")))
	}
	if len(parts) == 0 {
		return nil
	}
	return fmt.Errorf("command surface drift: %s", strings.Join(parts, " | "))
}
