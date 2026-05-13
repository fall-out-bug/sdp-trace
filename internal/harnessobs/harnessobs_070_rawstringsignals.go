package harnessobs

import (
	"strings"
)

func rawStringSignals(parentKey, value string) []string {
	// rawStringSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if rawSignalValueKey(parentKey) {

		return []string{strings.ToLower(value)}
	}
	return nil
}
