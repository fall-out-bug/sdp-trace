package harnessobs

import (
	"strings"
)

func unsafeIsolationRulePattern(pattern string) bool {
	return strings.TrimSpace(pattern) == "" || strings.Contains(pattern, "\n") || strings.Contains(pattern, "\r")
}
