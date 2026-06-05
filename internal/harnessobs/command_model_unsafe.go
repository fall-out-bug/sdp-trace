package harnessobs

import "strings"

func unsafeCommandModelIdentity(model string) bool {
	return model == "" || unsafeCommandModelChars(model)
}

func unsafeCommandModelChars(model string) bool {
	return strings.Contains(model, "://") || strings.ContainsAny(model, " \t\n\r\"'`$\\")
}

func unsafeCommandModelPath(model string) bool {
	return strings.Contains(model, "../") || strings.HasPrefix(model, "/")
}
