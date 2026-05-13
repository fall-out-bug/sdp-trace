package demo

import (
	"fmt"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"strings"
)

func overrideRequestFieldState(event trace.Event) (string, string) {
	// overrideRequestFieldState keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	for _, field := range []string{"override_id", "producer", "origin", "requested_by", "reason", "source_ref", "scope", "created_at"} {
		if strings.TrimSpace(payloadString(event, field)) == "" {
			return GateCannotVerify, fmt.Sprintf("override request missing %s", field)
		}
	}
	return GatePass, ""
}
