package harnessobs

import (
	"errors"

	"path/filepath"
)

func validateRawNormalizationInputs(format, rawPath, outPath string) error {
	// validateRawNormalizationInputs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if format != OpenCodeJSONLRawFormat {
		return errors.New("unsupported raw_event_format")
	}

	if filepath.Clean(rawPath) == filepath.Clean(outPath) {
		return errors.New("raw_event_source_path and event_source_path must be different files")
	}
	return nil
}
