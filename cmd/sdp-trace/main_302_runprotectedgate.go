package main

import (
	"io"
)

func runProtectedGate(target, outPath string, opts *flagSet, stdout, stderr io.Writer) int {
	// Protected gate resolution is separated from writing so input failures do
	// not create a partial gate artifact.
	result, code := resolveProtectedGate(target, opts, stderr)
	if code != 0 {
		return code
	}
	return writeProtectedGateResult(outPath, result, stdout, stderr)
}
