package prreview

import "strings"

// Reviewer text safety is conservative by design.
//
// Any obvious secret marker, credential URL, token string, or local private path
// redacts the full field before it can enter trust summaries or evidence refs.
func safeText(text string) string {
	if text == "" {
		return ""
	}
	if containsUnsafeTextMarker(text) || containsUnsafeTextPattern(text) {
		return redactedUnsafeReviewerText
	}
	return text
}

func containsUnsafeTextMarker(text string) bool {
	unsafeMarkers := []string{"SYNTHETIC_", "Bearer ", "access_token=", "BEGIN PRIVATE KEY", "PRIVATE_KEY", "cookie=", "session=", "/Users/", "/private/"}
	for _, marker := range unsafeMarkers {
		if strings.Contains(text, marker) {
			return true
		}
	}
	return false
}

func containsUnsafeTextPattern(text string) bool {
	return (strings.Contains(text, "://") && strings.Contains(text, "@")) || strings.Contains(text, "token=")
}
