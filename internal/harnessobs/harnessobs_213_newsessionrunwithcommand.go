package harnessobs

import (
	"time"
)

func newSessionRunWithCommand(profile SessionProfile, now time.Time, rawCommand string) SessionRun {
	// newSessionRunWithCommand keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.

	run := newSessionRun(profile, sessionRunTime(now))
	setSessionCommand(&run, rawCommand)
	return run
}
