package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func packetBuildPRRouteErrors(rows map[string]packet.Row) []string {
	route := rows["PC-AGENT-ROUTE"]
	if route.State == packet.StatePass || route.State == packet.StatePartial {
		return nil
	}
	// Missing route proof blocks PR packet readiness but remains cannot_verify.
	return []string{"PC-AGENT-ROUTE cannot verify live route proof: " + route.Reason}
}
