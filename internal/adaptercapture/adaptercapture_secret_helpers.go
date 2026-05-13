package adaptercapture

import (
	"strings"
)

func stringSliceContainsSecret(values []string) bool {
	// stringSliceContainsSecret preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	for _, value := range values {
		if containsSecret(value) {

			return true
		}
	}
	return false
}

func containsFold(value, needle string) bool {
	return strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}
