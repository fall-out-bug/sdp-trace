package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/interaction"
)

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
