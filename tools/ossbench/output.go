package main

import (
	"fmt"
	"io"
)

func printAndExit(stdout, stderr io.Writer, results []benchmarkResult, asJSON bool) int {
	if err := printResults(stdout, results, asJSON); err != nil {
		fmt.Fprintf(stderr, "print results: %v\n", err)
		return 2
	}
	return exitCode(results)
}
