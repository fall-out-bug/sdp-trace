package main

import (
	"fmt"
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func explainGateResult(result demo.GateResult, stdout io.Writer) {
	if result.SchemaVersion == demo.GateSchemaVersion {
		// Legacy gate results have no protected-profile fields; state that
		// absence explicitly rather than implying not_assessed conditions.
		fmt.Fprintln(stdout, "Protected profile fields: absent")
	}
	explainGateSummary(result, stdout)
	explainProtectedGateDetails(result, stdout)
	explainGateCollections(result, stdout)
}
