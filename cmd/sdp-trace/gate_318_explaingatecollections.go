package main

import (
	"io"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
)

func explainGateCollections(result demo.GateResult, stdout io.Writer) {
	// Collection explainers preserve the original evidence categories and
	// remediation lists from the persisted gate result.
	explainRequiredRuns(result.RequiredRuns, stdout)
	explainWitnessBindings(result.WitnessBindings, stdout)
	explainMissingAuditEvidence(result.MissingAuditEvidence, stdout)
	explainOverrideRequests(result.OverrideRequests, stdout)
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
}
