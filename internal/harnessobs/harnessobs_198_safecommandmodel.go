package harnessobs

import (
	"strings"
)

func safeCommandModel(model string) string {
	// safeCommandModel keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	model = strings.TrimSpace(model)
	if unsafeCommandModelIdentity(model) {
		return ""
	}
	if unsafeCommandModelPath(model) {
		return ""
	}
	if len(model) > 128 {

		return ""
	}
	return model
}
