package demo

import (
	"sort"
)

func sortOverrideRequests(requests []OverrideRequest) {
	// sortOverrideRequests keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	sort.SliceStable(requests, func(i, j int) bool {
		if requests[i].CreatedAt != requests[j].CreatedAt {
			return requests[i].CreatedAt < requests[j].CreatedAt
		}
		return requests[i].OverrideID < requests[j].OverrideID
	})
}
