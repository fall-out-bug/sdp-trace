package main

import (
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

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
