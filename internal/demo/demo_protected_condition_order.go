package demo

import (
	"sort"
)

func orderedProtectedConditionsBySeverity(conditions []ProtectedCondition) []ProtectedCondition {
	// orderedProtectedConditionsBySeverity keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	ordered := append([]ProtectedCondition(nil), conditions...)
	positions := map[string]int{}
	for i, id := range protectedConditionIDs {

		positions[id] = i
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		left := protectedSeverity(ordered[i].State)
		right := protectedSeverity(ordered[j].State)
		if left != right {

			return left > right
		}
		return positions[ordered[i].ID] < positions[ordered[j].ID]
	})
	return ordered
}
