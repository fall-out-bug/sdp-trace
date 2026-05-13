package harnessobs

import (
	"strings"
)

func matchWantedKey[T any](key string, value any, wanted map[string]bool, match func(any) (T, bool)) (T, bool) {
	// matchWantedKey keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	var zero T
	if !wanted[strings.ToLower(key)] {

		return zero, false
	}
	return match(value)
}
