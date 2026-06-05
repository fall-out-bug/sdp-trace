package harnessobs

import (
	"errors"
)

// RunSession wraps setup, command execution, and post-command collection while
// preserving the observed command's wait error for the caller.
func RunSession(opts SessionOptions) (SessionRun, Run, error) {
	session, err := setupRunnableSession(opts)
	if err != nil {
		return SessionRun{}, Run{}, err
	}
	commandResult, err := runObservedCommand(opts.Command, &session)
	if err != nil {
		return SessionRun{}, Run{}, err
	}

	return collectFinishedSession(opts, session, commandResult.waitErr, commandResult.end)
}

// setupRunnableSession requires a command before creating a session artifact so
// empty runtime invocations do not leave partial observation state.
func setupRunnableSession(opts SessionOptions) (SessionRun, error) {
	if err := requireSessionCommand(opts.Command); err != nil {
		return SessionRun{}, err
	}

	return SetupSession(SessionSetupOptions{ProfilePath: opts.ProfilePath, OutDir: opts.OutDir, Now: opts.Now})
}

// requireSessionCommand preserves the CLI error for missing commands after the
// `observe session --` separator.
func requireSessionCommand(command []string) error {
	if len(command) == 0 {
		return errors.New("observe session requires command after --")
	}
	return nil
}
