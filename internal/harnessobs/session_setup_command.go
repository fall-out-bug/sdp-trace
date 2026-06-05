package harnessobs

import (
	"strings"
	"time"
)

// newSessionRunWithCommand keeps command-derived facts on the session run
// without treating them as external proof.
func newSessionRunWithCommand(profile SessionProfile, now time.Time, rawCommand string) SessionRun {
	run := newSessionRun(profile, observationTime(now))
	setSessionCommand(&run, rawCommand)
	return run
}

// setSessionCommand stores a digest for every non-blank setup command and keeps
// the model fact only when command parsing returns a safe value.
func setSessionCommand(run *SessionRun, rawCommand string) {
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
