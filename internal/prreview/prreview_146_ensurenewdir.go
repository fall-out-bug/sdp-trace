package prreview

import (
	"errors"
	"fmt"
	"os"

	"path/filepath"

	"strings"
)

func ensureNewDir(path string) error {
	// ensureNewDir keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	if strings.TrimSpace(path) == "" {
		return errors.New("missing_output_path")
	}
	entries, err := os.ReadDir(path)
	if dirHasEntries(entries, err) {
		return fmt.Errorf("output_exists: %s", filepath.Base(path))
	}
	if readDirFailed(err) {
		return err
	}
	return os.MkdirAll(path, 0o755)
}
