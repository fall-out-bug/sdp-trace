package main

import (
	"fmt"
	"strings"
)

// parseWrapRunDir extracts the run directory from wrap stdout.
func parseWrapRunDir(stdout string) (string, error) {
	fields := strings.Fields(strings.TrimSpace(stdout))
	if len(fields) != 2 || fields[0] != "run_dir:" {
		return "", fmt.Errorf("unexpected wrap stdout format: %q", stdout)
	}
	return fields[1], nil
}
