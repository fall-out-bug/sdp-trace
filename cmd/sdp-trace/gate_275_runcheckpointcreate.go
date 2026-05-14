package main

import (
	"fmt"
	"io"
)

func runCheckpointCreate(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseCheckpointCreateArgs(args, stderr)
	if !ok {
		return code
	}
	// Creation persists the signed artifact before printing its id; stdout is
	// only a convenience pointer, not the checkpoint proof.
	created, code, ok := createCheckpointFromOptions(opts, stderr)
	if !ok {
		return code
	}
	fmt.Fprintf(stdout, "checkpoint: %s\n", created.CheckpointID)
	return 0
}
