package main

import (
	"strings"
)

func addAllowedRunnerItems(allowed map[string]bool, value string) {
	for _, item := range strings.Split(value, ",") {
		// Runner allow-lists accept comma-separated flags while preserving the
		// normalized set used by review validation.
		addAllowedRunnerItem(allowed, item)
	}
}
