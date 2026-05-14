package main

import (
	"io"
)

func requireOverrideRequestFlags(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	// Required field validation happens before the run directory is opened for
	// append, preventing partial override events.
	if !requireRequiredFlags(opts, stderr, overrideRequestRequiredFlags) {
		// Required fields identify who asked, what scope is affected, and which
		// source evidence the override references.
		return nil, exitUsage, false
	}
	return opts, 0, true
}
