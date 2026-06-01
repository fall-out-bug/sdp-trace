package harnessobs

import "strings"

// Key presence lookup builds the wanted-key set once, then delegates recursive
// map and slice traversal to the collection helpers.
func hasKey(value any, keys ...string) bool {
	wanted := map[string]bool{}
	for _, key := range keys {
		wanted[strings.ToLower(key)] = true
	}
	return hasKeyIn(value, wanted)
}

func hasKeyIn(value any, wanted map[string]bool) bool {
	switch v := value.(type) {
	case map[string]any:
		return hasKeyInMap(v, wanted)
	case []any:
		return hasKeyInSlice(v, wanted)
	}
	return false
}
