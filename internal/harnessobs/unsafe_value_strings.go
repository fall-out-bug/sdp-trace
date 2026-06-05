package harnessobs

import "strings"

// findUnsafeStringAt ignores blank strings, then classifies retained strings by
// path, authenticated URL, and token-like value rules.
func findUnsafeStringAt(path, value string, rawEvent bool) (string, string) {
	if strings.TrimSpace(value) == "" {
		return "", ""
	}
	if reason := unsafeStringReason(path, value, rawEvent); reason != "" {
		return path, reason
	}
	return "", ""
}

// unsafeStringReason keeps reason-code precedence stable for diagnostics and
// for tests that assert the first unsafe interpretation.
func unsafeStringReason(path, value string, rawEvent bool) string {
	if unsafeStringPath(value, rawEvent) {
		return "unsafe_path_or_private_path"
	}
	if unsafeURL(value) {
		return "authenticated_url"
	}
	if unsafeStringToken(path, value, rawEvent) {
		return "token_like_value"
	}
	return ""
}

// Path-looking strings are unsafe in retained payloads but allowed in raw-event
// fields that are later normalized under explicit raw path rules.
func unsafeStringPath(value string, rawEvent bool) bool {
	return !rawEvent && unsafePathValue(value)
}

// unsafePathValue rejects private absolute and parent-relative local paths,
// which should never become portable retained evidence values.
func unsafePathValue(value string) bool {
	return privatePathPattern.MatchString(value) || strings.HasPrefix(value, "/") || strings.Contains(value, "../")
}

// unsafeStringToken combines provider-specific token prefixes with encoded
// token detection while preserving digest/raw-path exemptions.
func unsafeStringToken(path, value string, rawEvent bool) bool {
	return providerTokenPrefix.MatchString(value) || unsafeEncodedToken(path, value, rawEvent)
}

// unsafeEncodedToken excludes known digest and raw path-like fields before
// applying the broad base64-like token detector.
func unsafeEncodedToken(path, value string, rawEvent bool) bool {
	if safeEncodedTokenExemption(path, value, rawEvent) {
		return false
	}
	return base64TokenPattern.MatchString(value)
}
