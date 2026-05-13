package harnessobs

import (
	"strings"
)

func hasSignal(signals []string, values ...string) bool {
	// hasSignal keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	wanted := map[string]bool{}
	for _, value := range values {

		wanted[strings.ToLower(value)] = true
	}

	for _, signal := range signals {
		if wanted[signal] {
			return true
		}
	}
	return false
}
