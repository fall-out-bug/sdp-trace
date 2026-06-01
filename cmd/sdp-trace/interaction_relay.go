package main

import (
	"context"
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/interaction"
)

func runInteractionRelay(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	// Relay records stdin plus command metadata before invoking the downstream
	// command, so the feedback event and command boundary stay coupled.
	opts, code, ok := parseInteractionRelayArgs(args, stderr)
	if !ok {
		return code
	}
	// Relay records the interaction before forwarding to the wrapped command, so
	// corrective feedback is not lost when the downstream command fails.
	_, exitCode, err := interaction.Relay(ctx, interactionRelayOptions(opts), cliStdin, stdout, stderr)
	if err != nil {
		// Relay package errors mean the interaction trace could not be recorded
		// with sufficient provenance.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return exitCode
}

func interactionRelayOptions(opts *flagSet) interaction.RelayOptions {
	return interaction.RelayOptions{
		// Identity fields bind the feedback event to a task and actor.
		TaskID:    opts.stringValue("task-id"),
		ActorType: opts.stringValue("actor-type"),
		ActorID:   opts.stringValue("actor-id"),
		// Target and event fields describe the interaction without trusting the
		// forwarded command to provide trace metadata.
		Target:      opts.stringValue("target"),
		EventType:   opts.stringValue("event-type"),
		OperationID: opts.stringValue("operation-id"),
		StageID:     opts.stringValue("stage-id"),
		// Out and Command define the durable trace location and replay boundary.
		Out:     opts.stringValue("out"),
		Command: opts.rest(),
	}
}
