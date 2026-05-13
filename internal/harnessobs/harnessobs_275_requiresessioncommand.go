package harnessobs

import (
	"errors"
)

func requireSessionCommand(command []string) error {
	// requireSessionCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if len(command) == 0 {

		return errors.New("observe session requires command after --")
	}
	return nil
}
