package prreview

import (
	"context"
	"errors"

	"os/exec"

	"strings"

	"time"
)

func runRoleCommand(packet Packet, role ReviewRole, opts RunOptions) ([]byte, bool, error) {
	// runRoleCommand keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.

	timeout := time.Duration(role.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, role.Command[0], role.Command[1:]...)
	cmd.Dir = opts.WorkDir
	prompt, err := renderPrompt(packet, role)
	if err != nil {
		return nil, false, err
	}
	if strings.TrimSpace(prompt) != "" {
		cmd.Stdin = strings.NewReader(prompt)
	}
	output, err := cmd.Output()
	return output, errors.Is(ctx.Err(), context.DeadlineExceeded), err
}
