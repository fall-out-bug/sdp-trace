package harnessobs

import (
	"strings"
)

func rawMapSignals(values map[string]any) []string {
	// rawMapSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	parts := make([]string, 0, len(values)*2)
	for key, child := range values {

		parts = append(parts, strings.ToLower(key))
		parts = append(parts, rawSignalsAt(key, child)...)
	}
	return parts
}
