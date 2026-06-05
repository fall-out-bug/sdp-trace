package prreview

import "strings"

const safeIDAllowedChars = "abcdefghijklmnopqrstuvwxyz0123456789_.-"

// Safe IDs preserve readable ASCII while removing path and shell surprises.
//
// Empty or fully stripped values fall back to `item`; string defaults preserve
// non-empty values exactly but replace whitespace-only values with the caller's
// explicit fallback.
func safeID(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	out := strings.Trim(strings.Map(safeIDMapper, value), "-.")
	if out == "" {
		return "item"
	}
	return out
}

func safeIDMapper(r rune) rune {
	if r <= 127 && strings.IndexByte(safeIDAllowedChars, byte(r)) >= 0 {
		return r
	}
	return '-'
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
