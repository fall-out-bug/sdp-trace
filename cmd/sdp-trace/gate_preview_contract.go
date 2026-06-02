package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func previewGateMode(contract trace.Contract) string {
	mode := demo.GateModeObservation
	for _, required := range contract.RequiredRuns {
		switch required.Profile {
		case demo.GateModeProtectedFuture:
			// Protected-future requirements dominate advisory CI requirements in
			// preview because they imply a stricter future gate path.
			return demo.GateModeProtectedFuture
		case demo.GateModeAdvisoryCI:
			// Advisory CI is retained unless a protected-future requirement is
			// found later.
			mode = demo.GateModeAdvisoryCI
		}
	}
	return mode
}

func requiredRunIDs(contract trace.Contract) []string {
	ids := make([]string, 0, len(contract.RequiredRuns))
	for _, required := range contract.RequiredRuns {
		if required.ID != "" {
			// Empty IDs are omitted from CLI preview output rather than rendered
			// as ambiguous evidence handles.
			ids = append(ids, required.ID)
		}
	}
	return ids
}

func requiredEvidenceIDsForCLI(contract trace.Contract) []string {
	ids := make([]string, 0, len(contract.RequiredEvidence))
	for _, requirement := range contract.RequiredEvidence {
		if requirement.ID != "" {
			// Preview exposes stable evidence identifiers only.
			ids = append(ids, requirement.ID)
		}
	}
	return ids
}
