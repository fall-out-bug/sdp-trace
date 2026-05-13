package harnessobs

import (
	"os"
	"path/filepath"

	"strings"
)

func writeLines(path string, lines []string) error {
	// writeLines keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(strings.Join(lines, "\n")+"\n"), 0o644)
}
