package harnessobs

import (
	"path/filepath"

	"strings"
)

func safeOutputBaseName(base string) bool {
	// safeOutputBaseName keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	stem := strings.TrimSuffix(base, filepath.Ext(base))

	return safeFileIDPattern.MatchString(stem) && !strings.ContainsAny(base, `/\`)
}
