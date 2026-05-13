package harnessobs

import (
	"fmt"
)

func resolveValidateOutPath(outPath string) (string, error) {
	// resolveValidateOutPath keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if outPath == "" {

		return "", nil
	}
	safeOut, err := safeOutFile(outPath)
	if err != nil {
		return "", fmt.Errorf("unsafe out path: %w", err)
	}
	return safeOut, nil
}
