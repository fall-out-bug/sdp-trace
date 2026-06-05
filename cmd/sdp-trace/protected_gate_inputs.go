package main

import (
	"fmt"
	"io"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

type protectedGateInputs struct {
	signed  checkpoint.SignedCheckpoint
	policy  checkpoint.TrustedCheckpointPolicy
	witness demo.WitnessSummary
}

func readProtectedGateInputs(opts *flagSet, stderr io.Writer) (protectedGateInputs, int, bool) {
	var inputs protectedGateInputs
	// Signed checkpoint data is the replay proof for the selected run.
	if code, ok := readRequiredProtectedInput("--checkpoint", opts.stringValue("checkpoint"), &inputs.signed, stderr); !ok {
		return protectedGateInputs{}, code, false
	}
	// Checkpoint policy pins the accepted signer authority.
	if code, ok := readRequiredProtectedInput("--checkpoint-policy", opts.stringValue("checkpoint-policy"), &inputs.policy, stderr); !ok {
		return protectedGateInputs{}, code, false
	}
	// Witness summary binds the protected run to observed CI evidence.
	if code, ok := readRequiredProtectedInput("--witness", opts.stringValue("witness"), &inputs.witness, stderr); !ok {
		return protectedGateInputs{}, code, false
	}
	return inputs, 0, true
}

func readRequiredProtectedInput(flag, path string, value any, stderr io.Writer) (int, bool) {
	if strings.TrimSpace(path) == "" {
		// Protected mode has no implicit local defaults for external trust
		// inputs.
		fmt.Fprintf(stderr, "protected gate requires %s\n", flag)
		return exitUsage, false
	}
	// All protected inputs are decoded as JSON artifacts before evaluation so
	// the gate never accepts unchecked path strings as trust evidence.
	if err := readJSONFile(path, value); err != nil {
		// Malformed trust inputs are usage/setup failures, not a green local gate
		// with omitted protected evidence.
		fmt.Fprintln(stderr, err)
		return exitUsage, false
	}
	return 0, true
}
