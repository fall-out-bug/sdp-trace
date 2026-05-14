package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func runReleaseProof(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseReleaseProofArgs(args, stderr)
	if !ok {
		return code
	}
	// Release proof is source-bound: the CLI must evaluate the current repo root
	// before writing any proof JSON that downstream gates might cite.
	result, code, ok := evaluateAndWriteReleaseProof(opts, stderr)
	if !ok {
		return code
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return releaseProofExitCode(result.ReleaseVerificationState)
}
