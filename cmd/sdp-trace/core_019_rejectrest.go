package main

import (
	"fmt"
	"io"
)

func rejectRest(opts *flagSet, stderr io.Writer, message string) bool {
	if len(opts.rest()) == 0 {
		return false
	}
	// Positional arguments are rejected before required flags so diagnostics
	// cannot imply that an ignored payload was accepted.
	fmt.Fprintln(stderr, message)
	return true
}
