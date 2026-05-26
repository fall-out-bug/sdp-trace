package main

import (
	"fmt"
	"io"
)

func gateTargetArg(opts *flagSet, stderr io.Writer) (string, bool) {
	targets := opts.rest()
	if len(targets) == 1 {
		return targets[0], true
	}
	// Gate evaluation is bound to exactly one run root or run directory.
	fmt.Fprintln(stderr, "gate requires <runs-root-or-run-dir>")
	return "", false
}
