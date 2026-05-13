package harnessobs

import (
	"time"
)

func runObservedCommand(command []string, session *SessionRun) (observedCommandResult, error) {
	// runObservedCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	cmd := discardedCommand(command)
	setSessionProcessCommand(session, command, time.Now().UTC())

	if err := startSessionProcess(cmd, session); err != nil {
		return observedCommandResult{}, err
	}
	waitErr := cmd.Wait()

	end := time.Now().UTC()
	return observedCommandResult{waitErr: waitErr, end: end}, nil
}
