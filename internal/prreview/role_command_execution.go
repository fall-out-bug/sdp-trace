package prreview

import (
	"context"
	"errors"
)

// runRoleCommand executes the configured role command with a rendered prompt
// and reports timeout separately from process failure.
func runRoleCommand(packet Packet, role ReviewRole, opts RunOptions) ([]byte, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), roleTimeout(role))
	defer cancel()
	cmd := roleCommand(ctx, role, opts.WorkDir)
	prompt, err := renderPrompt(packet, role, opts.PacketDir)
	if err != nil {
		return nil, false, err
	}
	attachPromptInput(cmd, prompt)
	output, err := cmd.Output()
	return output, roleCommandTimedOut(ctx), err
}

func roleCommandTimedOut(ctx context.Context) bool {
	return errors.Is(ctx.Err(), context.DeadlineExceeded)
}
