package main

import (
	"strings"
)

func addAllowedRunnerItem(allowed map[string]bool, item string) {
	item = strings.TrimSpace(item)
	if item != "" {
		// Empty allow-list entries are ignored so accidental commas do not create
		// wildcard-like runner names.
		allowed[item] = true
	}
}
