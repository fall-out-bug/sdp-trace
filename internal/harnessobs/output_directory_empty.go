package harnessobs

import (
	"errors"
	"os"
)

func ensureOutDirEmptyOrMissing(clean string) (string, error) {
	// Missing and empty directories are acceptable creation targets; existing
	// non-empty directories are refused to avoid mixing old and new evidence.
	entries, err := os.ReadDir(clean)
	if err := validateOutDirEntries(entries, err); err != nil {
		return "", err
	}

	return clean, nil
}

func validateOutDirEntries(entries []os.DirEntry, err error) error {
	// The ReadDir error is handled together with entries so callers keep a
	// single output-directory contract.
	if err == nil && len(entries) > 0 {
		return errors.New("harness observe refuses existing non-empty --out")
	}

	return ignorableMissingOutDir(err)
}

func ignorableMissingOutDir(err error) error {
	// A missing output directory is expected before artifact creation.
	if err == nil || errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
