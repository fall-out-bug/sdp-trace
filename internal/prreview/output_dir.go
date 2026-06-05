package prreview

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Output directories are created only when fresh.
//
// Empty paths fail with the stable missing_output_path contract. Existing
// non-empty directories are rejected before writes, while absent directories are
// created by the caller-owned packet or run preparation path.
// Read errors other than missing paths propagate unchanged so callers can
// distinguish inaccessible directories from normal first-time creation.
func ensureNewDir(path string) error {
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

func dirHasEntries(entries []os.DirEntry, err error) bool {
	return err == nil && len(entries) > 0
}

func readDirFailed(err error) bool {
	return err != nil && !errors.Is(err, os.ErrNotExist)
}
