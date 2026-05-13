package harnessobs

import (
	"time"
)

func setSessionProcessCommand(session *SessionRun, command []string, start time.Time) {
	// setSessionProcessCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	session.CommandDigest = digestCommand(command)
	session.CommandDigestState = StatePass
	if model := extractCommandModel(command); model != "" {
		session.CommandModel = model
		session.CommandModelState = StatePass
	}
	session.StartTime = start.Format(time.RFC3339)
}
