package main

import (
	"strings"
)

func allowedRunnerSet(values []string) map[string]bool {
	allowed := map[string]bool{}
	for _, value := range values {
		addAllowedRunnerItems(allowed, value)
	}
	// Empty input intentionally means no local external runners are allowed.
	return allowed
}

func addAllowedRunnerItems(allowed map[string]bool, value string) {
	for _, item := range strings.Split(value, ",") {
		// Runner allow-lists accept comma-separated flags while preserving the
		// normalized set used by review validation.
		addAllowedRunnerItem(allowed, item)
	}
}

func addAllowedRunnerItem(allowed map[string]bool, item string) {
	item = strings.TrimSpace(item)
	if item != "" {
		// Empty allow-list entries are ignored so accidental commas do not create
		// wildcard-like runner names.
		allowed[item] = true
	}
}
