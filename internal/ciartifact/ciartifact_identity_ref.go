package ciartifact

import "strings"

func safeRef(value string) bool {
	// safeRef keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if value == "" {
		return true
	}
	if !safeRefPrefix(value) {

		return false
	}
	return safeIdentityToken(value, "/._-")
}

func safeRefPrefix(value string) bool {
	// safeRefPrefix keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.

	return strings.HasPrefix(value, "refs/heads/") ||
		strings.HasPrefix(value, "refs/tags/") ||
		strings.HasPrefix(value, "refs/pull/") ||
		strings.HasPrefix(value, "refs/merge-requests/")
}
