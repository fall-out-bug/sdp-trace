package harnessobs

func findUnsafe(value any) (string, string) {
	return findUnsafeAt("", value)
}

func findUnsafeRawEvent(value any) (string, string) {
	return findUnsafeRawEventAt("", value)
}

func findUnsafeRawEventAt(path string, value any) (string, string) {
	return findUnsafeValueAt(path, value, true)
}

func findUnsafeAt(path string, value any) (string, string) {
	return findUnsafeValueAt(path, value, false)
}

func findUnsafeValueAt(path string, value any, rawEvent bool) (string, string) {
	// findUnsafeValueAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	switch v := value.(type) {
	case map[string]any:
		return findUnsafeMapAt(path, v, rawEvent)
	case []any:
		return findUnsafeSliceAt(path, v, rawEvent)
	case string:
		return findUnsafeStringAt(path, v, rawEvent)
	}

	return "", ""
}
