package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
)

func loadCheckpointVerifyInputs(opts *flagSet, stderr io.Writer) (checkpoint.SignedCheckpoint, *checkpoint.TrustedCheckpointPolicy, int, bool) {
	var signed checkpoint.SignedCheckpoint
	// Decode the signed checkpoint before policy so malformed proof artifacts
	// fail before optional trust policy is considered.
	if err := readJSONFile(opts.stringValue("checkpoint"), &signed); err != nil {
		fmt.Fprintln(stderr, err)
		return checkpoint.SignedCheckpoint{}, nil, 1, false
	}
	policy, err := readCheckpointPolicy(opts.stringValue("policy"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return checkpoint.SignedCheckpoint{}, nil, 1, false
	}
	return signed, policy, 0, true
}
