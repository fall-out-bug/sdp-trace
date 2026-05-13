package main

import (
	"strings"
)

func protectedInputStatus(path string) string {
	if strings.TrimSpace(path) == "" {
		// Missing protected inputs are reported as setup gaps, not as verdicts.
		return "absent"
	}
	var value any
	if err := readJSONFile(path, &value); err != nil {
		return protectedInputErrorStatus(err)
	}
	return "present_readable"
}
