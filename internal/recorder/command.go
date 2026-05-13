package recorder

import (
	"context"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// Command execution mirrors child output to the terminal while hashing the
// exact bytes that become recorder artifacts. Exit interpretation stays small
// and explicit because it feeds gate and closure evidence.
//
// The process wrapper does not reinterpret application output. It only wires
// stdio, waits for the child, and translates wrapper-observable failure modes
// into the exit code and signal fields written by the event layer.
//
// Keeping this layer policy-light matters because event construction is the
// only place where process facts become trace evidence.

func runCommand(ctx context.Context, command []string, env []string, writer *runWriter) (int, string) {
	// Process setup is isolated from waiting so start failures can be reported
	// without inventing an exit code from a process that never ran.
	cmd := recordedCommand(ctx, command, env, writer)
	if err := cmd.Start(); err != nil {
		return 1, "start_failed"
	}
	return waitCommand(ctx, cmd)
}

func recordedCommand(ctx context.Context, command []string, env []string, writer *runWriter) *exec.Cmd {
	// Stdout and stderr stay visible to the caller while the recorder hashes the
	// same bytes for later evidence comparison.
	cmd := exec.CommandContext(ctx, command[0], command[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin

	cmd.Stdout = io.MultiWriter(os.Stdout, &writer.stdoutHash)
	cmd.Stderr = io.MultiWriter(os.Stderr, &writer.stderrHash)
	return cmd
}

func waitCommand(ctx context.Context, cmd *exec.Cmd) (int, string) {
	// Wait distinguishes process exit, cancellation, and wrapper errors because
	// each has a different trust meaning in the terminal run record.
	err := cmd.Wait()
	if err == nil {
		return 0, ""
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), processSignal(exitErr.ProcessState)
	}
	if ctx.Err() != nil {
		return 1, "context_cancelled"
	}
	return 1, ""
}

func processSignal(processState *os.ProcessState) string {
	// Signal extraction is best-effort because non-Unix or missing process state
	// should not block recording the exit evidence we do have.
	if noProcessSignal(processState) {
		return ""
	}
	status, ok := processState.Sys().(syscall.WaitStatus)
	if !ok {
		return ""
	}
	return status.Signal().String()
}

func noProcessSignal(processState *os.ProcessState) bool {
	// A normal exit has no signal to report, and a missing process state cannot
	// add reliable signal evidence.
	return processState == nil || processState.Exited()
}
