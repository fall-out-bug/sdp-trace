package demo

import (
	"fmt"
)

func protectedReasons(conditions []ProtectedCondition) []string {
	// protectedReasons keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	ordered := orderedProtectedConditionsBySeverity(conditions)
	reasons := make([]string, 0, len(ordered))
	for _, condition := range ordered {
		if condition.ReasonCode == "" {
			continue
		}
		reasons = append(reasons, fmt.Sprintf("%s: %s", condition.ReasonCode, condition.Reason))
	}
	return reasons
}
