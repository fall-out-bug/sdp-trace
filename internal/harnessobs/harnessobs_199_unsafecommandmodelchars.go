package harnessobs

import (
	"strings"
)

func unsafeCommandModelChars(model string) bool {
	return strings.Contains(model, "://") || strings.ContainsAny(model, " \t\n\r\"'`$\\")
}
