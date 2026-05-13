package harnessobs

import (
	"strings"
)

func unsafeCommandModelPath(model string) bool {
	return strings.Contains(model, "../") || strings.HasPrefix(model, "/")
}
