package harnessobs

import (
	"errors"

	"os"
)

func validateOutDirEntries(entries []os.DirEntry, err error) error {
	// validateOutDirEntries keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err == nil && len(entries) > 0 {
		return errors.New("harness observe refuses existing non-empty --out")
	}

	return ignorableMissingOutDir(err)
}
