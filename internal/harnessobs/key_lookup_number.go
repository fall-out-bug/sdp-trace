package harnessobs

import "strings"

// Number key lookup is the typed adapter for numeric timestamp-like values
// that arrive as decoded JSON floats or local integer fixtures.
func findNumberByKey(value any, keys ...string) (float64, bool) {
	wanted := map[string]bool{}
	for _, key := range keys {
		wanted[strings.ToLower(key)] = true
	}
	return findNumberByKeyIn(value, wanted)
}

func findNumberByKeyIn(value any, wanted map[string]bool) (float64, bool) {
	matchNumber := func(value any) (float64, bool) {
		switch n := value.(type) {
		case float64:
			return n, true
		case int:
			return float64(n), true
		default:
			return 0, false
		}
	}
	return findByKeyIn(value, wanted, matchNumber)
}
