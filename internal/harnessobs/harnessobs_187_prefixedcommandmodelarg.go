package harnessobs

import (
	"strings"
)

func prefixedCommandModelArg(arg string) (string, bool) {
	// prefixedCommandModelArg keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	for _, prefix := range []string{"--model=", "-m="} {
		if strings.HasPrefix(arg, prefix) {

			return safeCommandModel(strings.TrimPrefix(arg, prefix)), true
		}
	}
	return "", false
}

// shellFields handles the shell field syntax needed to locate --model inside a
// controlled sh -c wrapper. It is not a general shell parser; model values still
// have to pass safeCommandModel before they become retained facts.
