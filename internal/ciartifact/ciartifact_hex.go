package ciartifact

import "strings"

func safeHex(value string, length int) bool {
	// safeHex keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	if len(value) != length {

		return false
	}
	for _, r := range value {
		if !isHexRune(r) {

			return false
		}
	}
	return true
}

func isHexRune(r rune) bool {
	return strings.ContainsRune("0123456789abcdefABCDEF", r)
}
