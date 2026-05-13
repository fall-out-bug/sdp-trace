package main

import (
	"io"
)

func runGateExplain(args []string, stdout, stderr io.Writer) int {
	path, code, ok := parseGateExplainArgs(args, stderr)
	if !ok {
		return code
	}
	result, code, ok := readGateExplainResult(path, stderr)
	if !ok {
		return code
	}
	// Explanation is read-only: it restates a persisted gate result without
	// re-running gate evaluation.
	explainGateResult(result, stdout)
	return 0
}
