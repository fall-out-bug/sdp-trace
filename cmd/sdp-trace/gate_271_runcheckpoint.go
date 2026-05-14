package main

import (
	"context"
	"io"
)

func runCheckpoint(_ context.Context, args []string, stdout, stderr io.Writer) int {
	// Checkpoint subcommands are explicit because create signs local run state
	// while verify only replays an existing signed checkpoint.
	command, rest, ok := checkpointCommand(args, stderr)
	if !ok {
		return exitUsage
	}
	return checkpointCommandHandlers[command](rest, stdout, stderr)
}
