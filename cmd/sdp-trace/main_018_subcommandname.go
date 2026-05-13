package main

import (
	"strings"
)

func subcommandName(label string) string {
	if before, _, ok := strings.Cut(label, " "); ok {
		// Help labels can include usage suffixes; dispatch diagnostics should
		// name only the stable subcommand token.
		return before
	}
	return label
}
