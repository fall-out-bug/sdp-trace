package harnessobs

// Raw signal dispatch keeps the recursive type switch in one place.
// Collection and value extraction live in sibling files so callers can follow
// the traversal without reopening the original numbered shards.
func rawSignals(value any) []string {
	return rawSignalsAt("", value)
}

func rawSignalsAt(parentKey string, value any) []string {
	if signals, ok := rawStructuredSignals(parentKey, value); ok {
		return signals
	}
	return rawLeafSignals(parentKey, value)
}

func rawStructuredSignals(parentKey string, value any) ([]string, bool) {
	if values, ok := value.(map[string]any); ok {
		return rawMapSignals(values), true
	}
	if values, ok := value.([]any); ok {
		return rawSliceSignals(parentKey, values), true
	}
	return nil, false
}

func rawLeafSignals(parentKey string, value any) []string {
	switch v := value.(type) {
	case string:
		return rawStringSignals(parentKey, v)
	default:
		return rawScalarSignals(v)
	}
}
