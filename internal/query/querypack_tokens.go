package query

import "strings"

func safeToken(value string) string {
	// safeToken keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	sanitized := sanitizeToken(value)
	if sanitized == "" {
		return "unknown"
	}
	return sanitized
}

const safeTokenAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-"

func sanitizeToken(value string) string {
	// sanitizeToken keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	var builder strings.Builder
	for _, r := range value {
		if isSafeTokenChar(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

func isSafeTokenChar(r rune) bool {
	return strings.ContainsRune(safeTokenAlphabet, r)
}
