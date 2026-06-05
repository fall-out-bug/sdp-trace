package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func packetBuildPRGateErrors(bundle packet.Bundle) []string {
	rows := map[string]packet.Row{}
	for _, row := range bundle.Packet.Rows {
		// Rows are keyed by packet id so the live gate checks can reference the
		// same identifiers as the rendered packet.
		rows[row.ID] = row
	}
	errors := []string{}
	// Route proof and CI verification are the live readiness rows for this CLI
	// build path.
	errors = append(errors, packetBuildPRRouteErrors(rows)...)
	errors = append(errors, packetBuildPRVerificationErrors(rows)...)
	return errors
}

func packetBuildPRRouteErrors(rows map[string]packet.Row) []string {
	route := rows["PC-AGENT-ROUTE"]
	if route.State == packet.StatePass || route.State == packet.StatePartial {
		return nil
	}
	// Missing route proof blocks PR packet readiness but remains cannot_verify.
	return []string{"PC-AGENT-ROUTE cannot verify live route proof: " + route.Reason}
}

func packetBuildPRVerificationErrors(rows map[string]packet.Row) []string {
	verification := rows["PC-VERIFICATION"]
	if verification.State == packet.StatePass {
		return nil
	}
	// CI evidence must be live enough for the packet row to pass before the build
	// command publishes artifacts.
	return []string{"PC-VERIFICATION cannot verify live CI evidence: " + verification.Reason}
}
