package harnessobs

import (
	"errors"

	"strings"
)

func resolveSessionSetupOutDir(outDir string) (string, error) {
	// resolveSessionSetupOutDir keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(outDir) == "" {
		return "", errors.New("observe setup requires --out")
	}

	return safeOutDir(outDir)
}
