package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

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
