package ciartifact

import "strings"

func safeIdentityToken(value, extra string) bool {
	// Identity tokens are allow-listed because artifact source labels can be echoed
	// into refs and reasons.
	if value == "" {
		return true
	}
	if !safeIdentityTokenLength(value) {

		return false
	}
	if unsafeIdentityValue(value) {

		return false
	}
	return safeIdentityTokenCharacters(value, extra)
}

func safeIdentityTokenLength(value string) bool {
	return len(value) <= 256
}

func safeIdentityTokenCharacters(value, extra string) bool {
	// safeIdentityTokenCharacters keeps CI artifact assessment evidence structured and replayable.
	// It preserves explicit states instead of turning local manifest data into proof.
	// Non-pass outcomes remain tied to a concrete source, family, or safety boundary.
	for _, r := range value {
		if !safeIdentityTokenRune(r, extra) {

			return false
		}
	}
	return true
}

func safeIdentityTokenRune(r rune, extra string) bool {
	return safeIdentityTokenAlnum(r) || strings.ContainsRune(extra, r)
}

func safeIdentityTokenAlnum(r rune) bool {
	return strings.ContainsRune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", r)
}
