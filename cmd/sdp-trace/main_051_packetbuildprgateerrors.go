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
