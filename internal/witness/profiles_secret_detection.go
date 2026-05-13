package witness

import (
	"strings"
)

func containsSecretLike(raw []byte) bool {
	// This detector is intentionally conservative. It protects witness outputs
	// from known secret shapes without trying to classify or print the secret.
	// Unknown content is not declared safe by this helper; it only blocks known
	// high-risk markers.
	// The final JWT-shape check covers tokens that lack a provider prefix.
	// All matching is done on lowercase text so provider marker casing cannot
	// bypass the deny-list.
	lower := strings.ToLower(string(raw))
	for _, marker := range secretSafetyMarkers {
		if strings.Contains(lower, marker) {
			return true
		}
	}
	return jwtLike(lower)
}

func jwtLike(text string) bool {
	// JWT detection runs on split fields only; callers get a boolean decision,
	// never decoded token claims or token material.
	for _, field := range jwtCandidateFields(text) {
		if jwtCandidate(field) {
			return true
		}
	}
	return false
}

func jwtCandidateFields(text string) []string {
	// Split on common JSON and prose separators so token-shaped substrings are
	// checked without decoding them.
	return strings.FieldsFunc(text, func(r rune) bool {
		switch r {
		case '"', '\'', ' ', '\n', '\t', '\r', ',', ':', '{', '}', '[', ']', '(', ')':
			return true
		default:
			return false
		}
	})
}

func jwtCandidate(field string) bool {
	parts := strings.Split(field, ".")
	// JWT-like values are refused by shape only; this catches common bearer
	// token leaks without logging or decoding token material.
	return len(parts) == 3 && strings.HasPrefix(parts[0], "eyj") && len(parts[1]) >= 8 && len(parts[2]) >= 8
}
