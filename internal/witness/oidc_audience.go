package witness

import "strings"

func audienceString(value any) string {
	// GitHub may encode audience as a string or array; normalize only those
	// shapes and reject everything else by returning an empty comparable value.
	switch typed := value.(type) {
	case string:
		return typed
	case []any:
		return strings.Join(stringItems(typed), ",")
	default:
		return ""
	}
}

func stringItems(values []any) []string {
	// Non-string audience entries are ignored instead of coerced, preserving a
	// strict comparison surface for the audience claim.
	parts := make([]string, 0, len(values))
	for _, item := range values {
		parts = appendStringItem(parts, item)
	}
	return parts
}

func appendStringItem(parts []string, item any) []string {
	// Only literal string audience values participate in the canonical audience
	// list.
	if text, ok := item.(string); ok {
		return append(parts, text)
	}
	return parts
}
