package harnessobs

import (
	"os"
)

func ensureOutDirEmptyOrMissing(clean string) (string, error) {
	// ensureOutDirEmptyOrMissing keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	entries, err := os.ReadDir(clean)
	if err := validateOutDirEntries(entries, err); err != nil {
		return "", err
	}

	return clean, nil
}
