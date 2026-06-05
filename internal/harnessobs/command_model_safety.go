package harnessobs

import "strings"

func safeCommandModel(model string) string {
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
