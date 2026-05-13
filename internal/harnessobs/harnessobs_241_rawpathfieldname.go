package harnessobs

import (
	"strings"
)

func rawPathFieldName(path string) string {
	// rawPathFieldName keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	field := path
	if idx := strings.LastIndex(field, "."); idx >= 0 {
		field = field[idx+1:]
	}
	if idx := strings.LastIndex(field, "["); idx >= 0 {

		field = field[:idx]
	}
	return strings.ToLower(field)
}
