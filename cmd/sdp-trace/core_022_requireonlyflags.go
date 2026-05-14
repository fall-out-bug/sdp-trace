package main

import (
	"io"
)

func requireOnlyFlags(opts *flagSet, stderr io.Writer, restMessage string, required []requiredCLIFlag) bool {
	if rejectRest(opts, stderr, restMessage) {
		return false
	}
	// Required flag checks run only after the command proves it received no
	// positional payload.
	return requireRequiredFlags(opts, stderr, required)
}
