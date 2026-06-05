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

func appendRepeatedFlagValue(values []string, args []string, i int, key, prefix string) ([]string, int) {
	arg := args[i]
	if strings.HasPrefix(arg, prefix) {
		// --key=value contributes exactly one ordered value.
		return append(values, strings.TrimPrefix(arg, prefix)), i
	}
	if arg == "--"+key && i+1 < len(args) {
		// --key value consumes the following argument as an ordered value.
		return append(values, args[i+1]), i + 1
	}
	return values, i
}
