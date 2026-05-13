package main

import (
	"fmt"
	"strings"
)

func requireOutputFile(command, path string) error {
	if strings.TrimSpace(path) == "" {
		// Commands that produce artifacts require an explicit destination to
		// avoid pretending stdout-only output is persisted evidence.
		return fmt.Errorf("%s requires --out", command)
	}
	return refuseExistingFile(path)
}
