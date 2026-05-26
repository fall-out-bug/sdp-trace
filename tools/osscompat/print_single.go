package main

import (
	"fmt"
	"io"
)

func printSingleProbeResult(stdout, stderr io.Writer, p probe, asJSON bool) int {
	// Keep single-probe output on the same path as all-probe output so JSON
	// formatting and failure exit semantics cannot drift between modes.
	r := runProbe(p)
	if err := printResults(stdout, []probeResult{r}, asJSON); err != nil {
		fmt.Fprintf(stderr, "print results: %v\n", err)
		return 2
	}
	return exitCode([]probeResult{r})
}
