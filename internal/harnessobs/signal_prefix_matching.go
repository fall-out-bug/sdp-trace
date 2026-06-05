package harnessobs

import "strings"

// Prefix signal matching is separate from exact matching because prefix
// probes are used for families of model, provider, and path signals.
func hasSignalPrefix(signals []string, prefixes ...string) bool {
	lowered := lowerStrings(prefixes)
	for _, signal := range signals {
		if signalHasPrefix(signal, lowered) {
			return true
		}
	}
	return false
}

func signalHasPrefix(signal string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(signal, prefix) {
			return true
		}
	}
	return false
}

func lowerStrings(values []string) []string {
	lowered := make([]string, 0, len(values))
	for _, value := range values {
		lowered = append(lowered, strings.ToLower(value))
	}
	return lowered
}
