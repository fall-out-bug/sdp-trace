package main

import (
	"io"
)

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
