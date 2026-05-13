package posture

import (
	"strings"
)

func hasUnsafeTokenOrCredential(value string) bool {
	// hasUnsafeTokenOrCredential keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if strings.Contains(value, "credential_or_token") {
		return false
	}
	return strings.Contains(value, "token") || strings.Contains(value, "credential")
}

func hasUnsafePath(value string) bool {
	return strings.Contains(value, "/") || strings.Contains(value, "\\")
}

func unsafeLabel(value string) bool {
	// unsafeLabel keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	lower := strings.ToLower(value)

	return unsafeOutput(value) ||
		strings.Contains(lower, "token") ||
		strings.Contains(lower, "credential")
}
