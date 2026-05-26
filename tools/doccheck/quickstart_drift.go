package main

import "strings"

func quickstartDrift(missing, stale []string) string {
	var parts []string
	if len(missing) > 0 {
		parts = append(parts, "missing required commands: "+strings.Join(missing, "; "))
	}
	if len(stale) > 0 {
		parts = append(parts, "stale commands: "+strings.Join(stale, "; "))
	}
	return strings.Join(parts, " | ")
}
