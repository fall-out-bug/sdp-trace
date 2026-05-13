package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func packetValidationExit(result packet.Validation) int {
	if result.State == packet.StatePass {
		return 0
	}
	// Packet validation failures mean the packet cannot be trusted as evidence,
	// so the CLI reports cannot_verify rather than a generic runtime failure.
	return exitCannotVerify
}
