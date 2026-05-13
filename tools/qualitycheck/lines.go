package main

import "strings"

func sourceLines(source string) int {
	// Trim only trailing newlines so blank lines inside the measured declaration
	// still contribute to the source span.
	trimmed := strings.TrimRight(source, "\n")
	if trimmed == "" {
		// Empty source slices should not inflate MI line counts.
		return 0
	}
	return strings.Count(trimmed, "\n") + 1
}
