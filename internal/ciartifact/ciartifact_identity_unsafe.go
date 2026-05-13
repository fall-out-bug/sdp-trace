package ciartifact

import "strings"

func unsafeIdentityValue(value string) bool {
	// unsafeIdentityValue keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	lower := strings.ToLower(value)

	if containsUnsafeIdentityMarker(lower) {
		return true
	}

	return strings.HasPrefix(value, "/") || strings.HasPrefix(value, "~")
}

func containsUnsafeIdentityMarker(lower string) bool {
	// containsUnsafeIdentityMarker keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, marker := range unsafeIdentityMarkers {
		if strings.Contains(lower, marker) {

			return true
		}
	}
	return false
}
