package main

import (
	"path/filepath"
)

func writableProbeParent(path string) string {
	target := filepath.Dir(path)
	if target == "" {
		// Empty dirname resolves to the current directory for local probes.
		return "."
	}
	return target
}
