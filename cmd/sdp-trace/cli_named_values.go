package main

import (
	"fmt"
	"io"
	"strings"
)

func requireNamedValues(values map[string]string, stderr io.Writer, messagePrefix string) bool {
	for flag, value := range values {
		if strings.TrimSpace(value) == "" {
			// The map caller owns flag order; this helper only centralizes the
			// trust-language for missing named inputs.
			fmt.Fprintf(stderr, "%s requires %s\n", messagePrefix, flag)
			return false
		}
	}
	return true
}
