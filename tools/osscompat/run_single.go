package main

import (
	"fmt"
	"io"
)

// runSingleProbe runs one probe by name and prints its result.
func runSingleProbe(stdout, stderr io.Writer, reg []probe, name string, asJSON bool) int {
	name = canonicalProbeName(name)
	for _, p := range reg {
		if p.Name == name {
			return printSingleProbeResult(stdout, stderr, p, asJSON)
		}
	}
	fmt.Fprintf(stderr, "unknown probe: %s\n", name)
	return 2
}
