package posture

import (
	"strings"
)

func unsafeSelectionPath(value string) bool {
	clean := strings.ReplaceAll(value, "\\", "/")
	return hasUnsafeSelectionPathPrefix(clean) || strings.Contains(clean, "../") || strings.Contains(clean, "/..")
}

func hasUnsafeSelectionPathPrefix(clean string) bool {
	// hasUnsafeSelectionPathPrefix keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	return strings.Contains(clean, "://") ||
		hasWindowsVolume(clean) ||
		strings.HasPrefix(clean, "/") ||
		strings.HasPrefix(clean, "../")
}

func hasWindowsVolume(value string) bool {
	return len(value) >= 3 && isASCIIAlpha(value[0]) && value[1] == ':' && value[2] == '/'
}

func isASCIIAlpha(value byte) bool {
	return (value >= 'A' && value <= 'Z') || (value >= 'a' && value <= 'z')
}
