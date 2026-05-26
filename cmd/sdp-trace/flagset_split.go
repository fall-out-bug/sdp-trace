package main

import "strings"

func splitFlag(arg string) (string, string, bool) {
	parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
	if len(parts) == 1 {
		// Bare flags may be boolean flags or string flags with next value.
		return parts[0], "", false
	}
	return parts[0], parts[1], true
}
