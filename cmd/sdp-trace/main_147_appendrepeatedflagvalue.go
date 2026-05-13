package main

import (
	"strings"
)

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
