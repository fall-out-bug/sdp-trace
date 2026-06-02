package harnessobs

import (
	"os/exec"
	"time"
)

// setSessionProcessCommand records deterministic command metadata before the
// process starts so the session still carries command facts on start failure.
func setSessionProcessCommand(session *SessionRun, command []string, start time.Time) {
	session.CommandDigest = digestCommand(command)
	session.CommandDigestState = StatePass
	if model := extractCommandModel(command); model != "" {
		session.CommandModel = model
		session.CommandModelState = StatePass
	}
	session.StartTime = start.Format(time.RFC3339)
}

// startSessionProcess records the operating-system process ID only after the
// child process successfully starts.
func startSessionProcess(cmd *exec.Cmd, session *SessionRun) error {
	if err := cmd.Start(); err != nil {
		return err
	}

	session.ProcessID = cmd.Process.Pid
	session.ProcessIDState = StatePass
	return nil
}
