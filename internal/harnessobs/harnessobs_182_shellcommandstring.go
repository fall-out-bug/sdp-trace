package harnessobs

import (
	"path/filepath"
)

func shellCommandString(command []string) string {
	// shellCommandString keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if !shellCommandShape(command) {
		return ""
	}
	base := filepath.Base(command[0])
	if base != "sh" && base != "bash" {
		return ""
	}

	return command[2]
}
