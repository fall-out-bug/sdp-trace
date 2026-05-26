package main

import "strings"

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
