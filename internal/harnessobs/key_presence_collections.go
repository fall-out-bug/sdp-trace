package harnessobs

import "strings"

// Key presence collection traversal searches nested maps and slices without
// interpreting values. Value interpretation belongs to key value lookup.
func hasKeyInMap(values map[string]any, wanted map[string]bool) bool {
	for key, child := range values {
		if mapKeyMatches(key, child, wanted) {
			return true
		}
	}
	return false
}

func mapKeyMatches(key string, child any, wanted map[string]bool) bool {
	return wanted[strings.ToLower(key)] || hasKeyIn(child, wanted)
}

func hasKeyInSlice(values []any, wanted map[string]bool) bool {
	for _, child := range values {
		if hasKeyIn(child, wanted) {
			return true
		}
	}
	return false
}
