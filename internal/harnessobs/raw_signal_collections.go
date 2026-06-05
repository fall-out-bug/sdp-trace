package harnessobs

import "strings"

// Raw signal collection extraction walks maps and slices while preserving the
// parent key. That parent key decides whether string leaves become signals.
func rawMapSignals(values map[string]any) []string {
	parts := make([]string, 0, len(values)*2)
	for key, child := range values {
		parts = append(parts, strings.ToLower(key))
		parts = append(parts, rawSignalsAt(key, child)...)
	}
	return parts
}

func rawSliceSignals(parentKey string, values []any) []string {
	parts := make([]string, 0, len(values))
	for _, child := range values {
		parts = append(parts, rawSignalsAt(parentKey, child)...)
	}
	return parts
}
