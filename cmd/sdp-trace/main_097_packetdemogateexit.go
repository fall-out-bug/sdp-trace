package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func packetDemoGateExit(result packet.Validation) int {
	if result.State == packet.StatePass {
		return 0
	}
	// Demo gate validation is an expected fail/pass contract for fixtures.
	return exitFail
}
