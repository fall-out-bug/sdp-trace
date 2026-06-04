package prreview

import (
	"context"
	"os/exec"
	"strings"
	"time"
)

// roleTimeout applies the default runner timeout when a role omits an explicit
// value.
func roleTimeout(role ReviewRole) time.Duration {
	timeout := time.Duration(role.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		return 30 * time.Second
	}
	return timeout
}

// roleCommand binds the role command to the caller-selected work directory.
func roleCommand(ctx context.Context, role ReviewRole, workDir string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, role.Command[0], role.Command[1:]...)
	cmd.Dir = workDir
	return cmd
}

// attachPromptInput sends non-empty rendered prompts to the runner stdin.
func attachPromptInput(cmd *exec.Cmd, prompt string) {
	if strings.TrimSpace(prompt) != "" {
		cmd.Stdin = strings.NewReader(prompt)
	}
}
