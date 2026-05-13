package posture

import (
	"strings"
)

var unsafeOutputKeywords = []string{
	"http://",
	"https://",
	"secret",
	"@",
}

func unsafeOutput(value string) bool {
	// unsafeOutput keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	lower := strings.ToLower(value)
	if hasUnsafeOutputKeyword(lower) {

		return true
	}
	if hasUnsafeTokenOrCredential(lower) {

		return true
	}
	if hasUnsafePath(value) {

		return true
	}
	return false
}

func hasUnsafeOutputKeyword(value string) bool {
	// hasUnsafeOutputKeyword keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	for _, keyword := range unsafeOutputKeywords {
		if strings.Contains(value, keyword) {
			return true
		}
	}
	return false
}
