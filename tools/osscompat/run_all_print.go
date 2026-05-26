package main

import (
	"fmt"
	"io"
)

// runAllAndPrint runs every probe and prints the results.
func runAllAndPrint(stdout, stderr io.Writer, reg []probe, asJSON bool) int {
	results := runAllProbes(reg)
	if err := printResults(stdout, results, asJSON); err != nil {
		fmt.Fprintf(stderr, "print results: %v\n", err)
		return 2
	}
	return exitCode(results)
}
