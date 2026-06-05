package harnessobs

// Generic key lookup routes recursive map and slice traversal for callers that
// provide their own typed value matcher.
func findByKeyIn[T any](value any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	var zero T
	switch v := value.(type) {
	case map[string]any:
		return findByKeyInMap(v, wanted, match)
	case []any:
		return findByKeyInSlice(v, wanted, match)
	}
	return zero, false
}

func findByKeyInSlice[T any](value []any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	var zero T
	for _, child := range value {
		if found, ok := findByKeyIn(child, wanted, match); ok {
			return found, true
		}
	}
	return zero, false
}
