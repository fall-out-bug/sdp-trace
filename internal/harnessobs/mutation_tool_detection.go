package harnessobs

import "strings"

// Native mutation detection recognizes harness tools that can mutate source
// files. It only classifies observed raw events; it does not grant proof.
func nativeMutationTool(raw map[string]any) bool {
	tool := strings.ToLower(findStringByKey(raw, "tool"))
	switch tool {
	case "edit", "write", "patch", "apply_patch", "update", "delete":
		return true
	default:
		return false
	}
}
