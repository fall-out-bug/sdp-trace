package harnessobs

import "strings"

// Map key lookup checks the current key before recursing into the child value.
// This preserves first-match traversal order from the original numbered files.
func findByKeyInMap[T any](value map[string]any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	var zero T
	for key, child := range value {
		if found, ok := findByMapEntry(key, child, wanted, match); ok {
			return found, true
		}
	}
	return zero, false
}

func findByMapEntry[T any](key string, child any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	if found, ok := matchWantedKey(key, child, wanted, match); ok {
		return found, true
	}
	return findByKeyIn(child, wanted, match)
}

func matchWantedKey[T any](key string, value any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	var zero T
	if !wanted[strings.ToLower(key)] {
		return zero, false
	}
	return match(value)
}
