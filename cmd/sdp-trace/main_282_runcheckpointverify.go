package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func runCheckpointVerify(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseCheckpointVerifyArgs(args, stderr)
	if !ok {
		return code
	}
	// Verification reads immutable inputs first, then delegates all replay and
	// signer-policy decisions to the checkpoint package.
	// The CLI only renders the resulting verification envelope and exit state.
	signed, policy, code, ok := loadCheckpointVerifyInputs(opts, stderr)
	if !ok {
		return code
	}
	result := checkpoint.Verify(opts.stringValue("run"), signed, policy)
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return checkpointVerifyExitCode(result.Result)
}
