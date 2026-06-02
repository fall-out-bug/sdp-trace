package harnessobs

import "strings"

// findUnsafeMapAt walks maps until the first unsafe retained field is found;
// the caller only needs one field/reason pair to reject the payload.
func findUnsafeMapAt(path string, values map[string]any, rawEvent bool) (string, string) {
	for key, child := range values {
		if field, reason := findUnsafeMapChild(path, key, child, rawEvent); field != "" {
			return field, reason
		}
	}
	return "", ""
}

// findUnsafeMapChild applies field-level raw-event rules before recursing into
// nested values, so explicitly forbidden retained fields win over child checks.
func findUnsafeMapChild(path, key string, child any, rawEvent bool) (string, string) {
	childPath := childPath(path, key)
	reason, skip := unsafeMapFieldReason(childPath, strings.ToLower(key), child, rawEvent)
	if reason != "" {
		return childPath, reason
	}
	if skip {
		return "", ""
	}
	return findUnsafeValueAt(childPath, child, rawEvent)
}

func childPath(parent, key string) string {
	// Root-level keys stay unprefixed so error fields match raw payload names.
	if parent == "" {
		return key
	}
	return parent + "." + key
}

// unsafeMapFieldReason distinguishes forbidden retained fields from raw-event
// fields that are allowed only because they are intentionally not retained.
func unsafeMapFieldReason(path, key string, value any, rawEvent bool) (string, bool) {
	if skippableRawEventField(path, key, value, rawEvent) {
		return "", true
	}
	if rawFieldNames[key] {
		return "forbidden_raw_field", false
	}
	if sensitiveFieldNames[key] {
		return "sensitive_field", false
	}
	return "", false
}
