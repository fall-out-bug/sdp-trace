package harnessobs

import (
	"strings"
)

func nativeMutationTool(raw map[string]any) bool {
	// nativeMutationTool keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	tool := strings.ToLower(findStringByKey(raw, "tool"))
	switch tool {
	case "edit", "write", "patch", "apply_patch", "update", "delete":

		return true
	default:
		return false
	}
}
