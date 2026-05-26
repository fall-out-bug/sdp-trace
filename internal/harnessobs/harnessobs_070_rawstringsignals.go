package harnessobs

import (
	"strings"
)

func rawStringSignals(parentKey, value string) []string {
	// rawStringSignals keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if rawSignalValueKey(parentKey) {
		signal := strings.ToLower(value)
		if strings.Contains(signal, "/.planning/phases/") ||
			strings.HasPrefix(signal, ".planning/phases/") {
			return []string{signal, "gsd.phase_path"}
		}

		return []string{signal}
	}
	return nil
}
