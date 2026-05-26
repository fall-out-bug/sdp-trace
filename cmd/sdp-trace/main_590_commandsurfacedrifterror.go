package main

import (
	"fmt"
	"strings"
)

func commandSurfaceDriftError(missing, stale []string) error {
	// Empty diff means the registry and help text agree, so tests should return
	// nil instead of formatting an empty diagnostic.
	parts := commandSurfaceDriftParts(missing, stale)
	if len(parts) == 0 {
		return nil
	}
	return fmt.Errorf("command surface drift: %s", strings.Join(parts, " | "))
}
