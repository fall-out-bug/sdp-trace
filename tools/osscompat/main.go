package main

import (
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, registry))
}

// run parses flags and executes the requested mode.
func run(args []string, stdout, stderr io.Writer, reg []probe) int {
	_, asJSON, list, probeName, code := parseFlags(args, stderr)
	if code >= 0 {
		return code
	}
	if *list {
		return listProbes(stdout, stderr, reg)
	}
	if *probeName != "" {
		return runSingleProbe(stdout, stderr, reg, *probeName, *asJSON)
	}
	return runAllAndPrint(stdout, stderr, reg, *asJSON)
}
