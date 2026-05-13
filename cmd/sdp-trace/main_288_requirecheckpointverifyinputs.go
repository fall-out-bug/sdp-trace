package main

import (
	"fmt"
	"strings"
)

func requireCheckpointVerifyInputs(opts *flagSet) error {
	if strings.TrimSpace(opts.stringValue("run")) == "" {
		// The run directory is the source replay target for the signed payload.
		return fmt.Errorf("checkpoint verify requires --run")
	}
	if strings.TrimSpace(opts.stringValue("checkpoint")) == "" {
		// The signed checkpoint artifact is mandatory verification input.
		return fmt.Errorf("checkpoint verify requires --checkpoint")
	}
	return nil
}
