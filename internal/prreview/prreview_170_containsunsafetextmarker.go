package prreview

import (
	"strings"
)

func containsUnsafeTextMarker(text string) bool {
	// containsUnsafeTextMarker keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	unsafeMarkers := []string{"SYNTHETIC_", "Bearer ", "access_token=", "BEGIN PRIVATE KEY", "PRIVATE_KEY", "cookie=", "session=", "/Users/", "/private/"}
	for _, marker := range unsafeMarkers {
		if strings.Contains(text, marker) {

			return true
		}
	}
	return false
}
