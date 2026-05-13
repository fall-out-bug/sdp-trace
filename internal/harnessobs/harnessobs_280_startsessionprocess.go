package harnessobs

import (
	"os/exec"
)

func startSessionProcess(cmd *exec.Cmd, session *SessionRun) error {
	// startSessionProcess keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	if err := cmd.Start(); err != nil {
		return err
	}

	session.ProcessID = cmd.Process.Pid
	session.ProcessIDState = StatePass
	return nil
}
