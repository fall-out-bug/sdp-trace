package harnessobs

import (
	"io"

	"os/exec"
)

func discardedCommand(command []string) *exec.Cmd {
	// discardedCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd
}
