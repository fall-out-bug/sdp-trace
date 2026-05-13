package harnessobs

import (
	"strings"
)

func hasSignalPrefix(signals []string, prefixes ...string) bool {
	// hasSignalPrefix keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, signal := range signals {
		for _, prefix := range prefixes {

			if strings.HasPrefix(signal, strings.ToLower(prefix)) {
				return true
			}
		}
	}
	return false
}
