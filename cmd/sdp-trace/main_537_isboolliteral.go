package main

import (
	"strings"
)

func isBoolLiteral(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "false" || lower == "1" || lower == "0"
}
