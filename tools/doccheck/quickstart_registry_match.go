package main

import "strings"

// prefixMatchesRegistry reports whether normalized starts with any registry
// usage prefix (the stable part before the first optional flag or placeholder).
func prefixMatchesRegistry(normalized string, registrySet map[string]bool) bool {
	for reg := range registrySet {
		// Quickstart examples may omit optional flags or concrete placeholders,
		// so compare against the stable leading command phrase.
		prefix := registryPrefix(reg)
		if prefix != "" && strings.HasPrefix(normalized, prefix) {
			return true
		}
	}
	return false
}

// registryPrefix returns the stable prefix of a registry usage string,
// stopping before the first optional flag ([) or placeholder (<).
func registryPrefix(usage string) string {
	for i, ch := range usage {
		if ch == '[' || ch == '<' {
			// Optional flags and placeholders start variable syntax, not the
			// command identity used for drift matching.
			return strings.TrimSpace(usage[:i])
		}
	}
	return usage
}
