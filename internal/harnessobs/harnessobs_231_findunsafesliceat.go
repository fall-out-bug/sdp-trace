package harnessobs

import (
	"fmt"
)

func findUnsafeSliceAt(path string, values []any, rawEvent bool) (string, string) {
	// findUnsafeSliceAt keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for i, child := range values {

		if field, reason := findUnsafeValueAt(fmt.Sprintf("%s[%d]", path, i), child, rawEvent); field != "" {
			return field, reason
		}
	}
	return "", ""
}
