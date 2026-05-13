package harnessobs

import (
	"strings"
)

func setSessionCommand(run *SessionRun, rawCommand string) {
	// setSessionCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if strings.TrimSpace(rawCommand) == "" {
		return
	}

	command := []string{rawCommand}
	run.CommandDigest = digestCommand(command)
	run.CommandDigestState = StatePass
	if model := extractCommandModel(command); model != "" {
		run.CommandModel = model
		run.CommandModelState = StatePass
	}
}
