package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/posture"
)

func runCrossRepoPostureExplain(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseCrossRepoPostureExplainArgs(args, stderr)
	if !ok {
		return code
	}
	// Explanation renders a saved posture export; it never rebuilds selection
	// state from the workspace.
	result, code, ok := readCrossRepoPostureExplainResult(opts.stringValue("result"), stderr)
	if !ok {
		return code
	}
	rendered, err := posture.Explain(result)
	if err != nil {
		// Unsafe rendered text is a verification failure for the explanation.
		fmt.Fprintln(stderr, "output_safety_violation")
		return exitCannotVerify
	}
	// The explanation is intentionally stdout-only so it cannot be mistaken for
	// a new posture evidence artifact.
	fmt.Fprint(stdout, rendered)
	return 0
}
