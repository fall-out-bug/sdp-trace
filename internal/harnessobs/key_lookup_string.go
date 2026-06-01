package harnessobs

import "strings"

// String key lookup is the typed adapter for non-empty string values found by
// the generic recursive key traversal.
func findStringByKey(value any, keys ...string) string {
	wanted := map[string]bool{}
	for _, key := range keys {
		wanted[strings.ToLower(key)] = true
	}
	return findStringByKeyIn(value, wanted)
}

func findStringByKeyIn(value any, wanted map[string]bool) string {
	matchingString := func(value any) (string, bool) {
		s, ok := value.(string)
		return s, ok && strings.TrimSpace(s) != ""
	}
	s, ok := findByKeyIn(value, wanted, matchingString)
	if !ok {
		return ""
	}
	return s
}
