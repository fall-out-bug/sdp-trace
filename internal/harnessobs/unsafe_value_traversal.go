package harnessobs

import "fmt"

func findUnsafe(value any) (string, string) {
	// Generic retained payload scans start at the root with raw-event rules off.
	return findUnsafeAt("", value)
}

func findUnsafeRawEvent(value any) (string, string) {
	// Raw-event scans keep raw-event mode on so unretained raw fields can be
	// skipped while forbidden retained raw fields still fail.
	return findUnsafeRawEventAt("", value)
}

func findUnsafeRawEventAt(path string, value any) (string, string) {
	return findUnsafeValueAt(path, value, true)
}

func findUnsafeAt(path string, value any) (string, string) {
	return findUnsafeValueAt(path, value, false)
}

func findUnsafeValueAt(path string, value any, rawEvent bool) (string, string) {
	// Dispatch preserves raw-event mode while recursively checking only the
	// retained map, slice, and string value shapes that can carry unsafe data.
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

func findUnsafeSliceAt(path string, values []any, rawEvent bool) (string, string) {
	for i, child := range values {
		// Slice indexes become part of the diagnostic field path so callers can
		// identify the first unsafe nested value.
		if field, reason := findUnsafeValueAt(fmt.Sprintf("%s[%d]", path, i), child, rawEvent); field != "" {
			return field, reason
		}
	}
	return "", ""
}
