package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func runQueryPackExplain(args []string, stdout, stderr io.Writer) int {
	opts, err := parseQueryPackExplainArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	// Explain renders an existing query-pack artifact; it does not rebuild the
	// pack or re-query the original run evidence.
	result, err := readQueryPackResult(opts.resultPath)
	if err != nil {
		// Explain is artifact-only; it cannot reconstruct missing pack results
		// from the original run.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	if err := validateQueryPackExplainResult(result); err != nil {
		// Schema/profile mismatch means this binary cannot render the artifact
		// without risking stale explanation semantics.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	fmt.Fprint(stdout, query.ExplainForensicsPack(result))
	return 0
}
