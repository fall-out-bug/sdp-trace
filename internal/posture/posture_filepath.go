package posture

import (
	"strings"
)

// filepathBase extracts the basename for manifest-artifact matching.
// Slash normalization ensures Windows-style separators cannot escape the comparison.
func filepathBase(path string) string {
	// filepathBase keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	clean := strings.ReplaceAll(path, "\\", "/")
	parts := strings.Split(clean, "/")
	return parts[len(parts)-1]
}
