package main

import (
	"io"
)

func requireRequiredFlags(opts *flagSet, stderr io.Writer, required []requiredCLIFlag) bool {
	for _, flag := range required {
		// Preserve caller-specific messages so command docs and tests can pin the
		// exact missing input.
		if !requireStringFlag(opts, stderr, flag.name, flag.message) {
			return false
		}
	}
	return true
}
