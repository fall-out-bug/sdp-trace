package main

import (
	"strings"
)

func repeatedFlagValues(args []string, key, parsedFallback string) []string {
	prefix := "--" + key + "="
	values := []string{}
	for i := 0; i < len(args); i++ {
		// Raw args preserve repeated flag order that the simple parser collapses.
		values, i = appendRepeatedFlagValue(values, args, i, key, prefix)
	}
	if len(values) == 0 && strings.TrimSpace(parsedFallback) != "" {
		// The parsed fallback covers the single-value case.
		values = append(values, parsedFallback)
	}
	return values
}
