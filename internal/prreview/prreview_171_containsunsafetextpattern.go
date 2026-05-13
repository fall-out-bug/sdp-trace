package prreview

import (
	"strings"
)

func containsUnsafeTextPattern(text string) bool {
	return (strings.Contains(text, "://") && strings.Contains(text, "@")) || strings.Contains(text, "token=")
}
