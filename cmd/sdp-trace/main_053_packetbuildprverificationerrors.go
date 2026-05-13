package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func packetBuildPRVerificationErrors(rows map[string]packet.Row) []string {
	verification := rows["PC-VERIFICATION"]
	if verification.State == packet.StatePass {
		return nil
	}
	// CI evidence must be live enough for the packet row to pass before the build
	// command publishes artifacts.
	return []string{"PC-VERIFICATION cannot verify live CI evidence: " + verification.Reason}
}
