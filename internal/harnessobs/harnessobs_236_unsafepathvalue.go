package harnessobs

import (
	"strings"
)

func unsafePathValue(value string) bool {
	return privatePathPattern.MatchString(value) || strings.HasPrefix(value, "/") || strings.Contains(value, "../")
}
