package harnessobs

import "strings"

func prefixedCommandModelArg(arg string) (string, bool) {
	for _, prefix := range []string{"--model=", "-m="} {
		if strings.HasPrefix(arg, prefix) {
			return safeCommandModel(strings.TrimPrefix(arg, prefix)), true
		}
	}
	return "", false
}
