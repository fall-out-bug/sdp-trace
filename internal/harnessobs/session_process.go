package harnessobs

import (
	"io"
	"os/exec"
	"time"
)

// runObservedCommand records command metadata before process start and returns
// both the wait result and the actual end time used for collection.
func runObservedCommand(command []string, session *SessionRun) (observedCommandResult, error) {
	cmd := discardedCommand(command)
	setSessionProcessCommand(session, command, time.Now().UTC())

	if err := startSessionProcess(cmd, session); err != nil {
		return observedCommandResult{}, err
	}
	waitErr := cmd.Wait()

	end := time.Now().UTC()
	return observedCommandResult{waitErr: waitErr, end: end}, nil
}

// discardedCommand intentionally suppresses child process stdio; observation
// evidence comes from configured event sources rather than terminal capture.
func discardedCommand(command []string) *exec.Cmd {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd
}
