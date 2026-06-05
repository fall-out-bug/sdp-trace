package harnessobs

import "strings"

// Exact signal matching normalizes requested values once and compares them
// against raw signals that are already lower-cased by extraction.
func hasSignal(signals []string, values ...string) bool {
	wanted := lowerStringSet(values)
	for _, signal := range signals {
		if wanted[signal] {
			return true
		}
	}
	return false
}

func lowerStringSet(values []string) map[string]bool {
	wanted := map[string]bool{}
	for _, value := range values {
		wanted[strings.ToLower(value)] = true
	}
	return wanted
}
