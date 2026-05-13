package main

import (
	"io"
)

func requireOnlyFlagsCode(opts *flagSet, stderr io.Writer, restMessage string, required []requiredCLIFlag) (int, bool) {
	if !requireOnlyFlags(opts, stderr, restMessage, required) {
		// Keep parser helpers returning CLI usage codes instead of package
		// verifier states.
		return exitUsage, false
	}
	return 0, true
}
