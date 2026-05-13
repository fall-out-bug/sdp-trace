package harnessobs

import (
	"strings"
)

func digestField(path string) bool {
	// digestField keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	last := path
	if idx := strings.LastIndex(last, "."); idx >= 0 {

		last = last[idx+1:]
	}
	return digestFieldNames[last]
}
