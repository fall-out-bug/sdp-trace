package demo

import (
	"strings"
)

func protectedNextActions(conditions []ProtectedCondition) []string {
	// protectedNextActions keeps demo gate evidence explicit and replay-bound.
	// Local rows, contract requirements, witness bindings, protected gates, and artifacts stay separate.
	// This helper renders or aggregates demo evidence; it does not create external proof.

	ordered := orderedProtectedConditionsBySeverity(conditions)
	actions := make([]string, 0, len(ordered))
	for _, condition := range ordered {
		if strings.TrimSpace(condition.NextAction) != "" {
			actions = append(actions, condition.NextAction)
		}
	}
	return actions
}
