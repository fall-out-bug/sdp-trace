package main

import (
	"fmt"
	"io"
)

func reportTargetArg(opts *flagSet, stderr io.Writer) (string, bool) {
	targets := opts.rest()
	if len(targets) != 1 {
		// A single target keeps report provenance bound to one run root.
		fmt.Fprintln(stderr, "report requires <runs-root-or-run-dir>")
		return "", false
	}
	return targets[0], true
}
