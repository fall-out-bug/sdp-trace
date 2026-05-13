package harnessobs

import (
	"os"
)

func readExistingJSONStrict(path string, target any) error {
	// readExistingJSONStrict keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	safePath, err := safeExistingFile(path)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(safePath)
	if err != nil {
		return err
	}
	return decodeStrictJSON(data, target)
}
