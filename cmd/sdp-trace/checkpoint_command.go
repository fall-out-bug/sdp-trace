package main

import (
	"context"
	"fmt"
	"io"
)

var checkpointCommandHandlers = map[string]subcommandHandler{
	"create": runCheckpointCreate,
	"verify": runCheckpointVerify,
}

func runCheckpoint(_ context.Context, args []string, stdout, stderr io.Writer) int {
	// Checkpoint subcommands are explicit because create signs local run state
	// while verify only replays an existing signed checkpoint.
	command, rest, ok := checkpointCommand(args, stderr)
	if !ok {
		return exitUsage
	}
	return checkpointCommandHandlers[command](rest, stdout, stderr)
}

func checkpointCommand(args []string, stderr io.Writer) (string, []string, bool) {
	if len(args) == 0 {
		// A checkpoint command without a verb cannot decide whether it is
		// creating signing evidence or replaying it.
		fmt.Fprintln(stderr, "checkpoint requires create or verify")
		return "", nil, false
	}
	// Keep the remainder untouched so each subcommand owns its flag contract.
	return knownCheckpointCommand(args[0], args[1:], stderr)
}

func knownCheckpointCommand(command string, rest []string, stderr io.Writer) (string, []string, bool) {
	if command != "create" && command != "verify" {
		// Unknown checkpoint verbs are usage errors, not verifier states.
		fmt.Fprintln(stderr, "checkpoint requires create or verify")
		return "", nil, false
	}
	return command, rest, true
}
