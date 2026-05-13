package managed

import "sort"

func orderConditions(conditions []Condition) []Condition {
	// orderConditions preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.

	index := conditionOrderIndex()
	ordered := append([]Condition(nil), conditions...)
	sort.SliceStable(ordered, func(i, j int) bool {
		return conditionLess(ordered, index, i, j)
	})
	return ordered
}

func conditionOrderIndex() map[string]int {
	// conditionOrderIndex preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	index := map[string]int{}
	for i, id := range managedConditionIDs {

		index[id] = i
	}
	return index
}

func conditionLess(ordered []Condition, index map[string]int, i, j int) bool {
	// conditionLess preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if severity(ordered[i].State) != severity(ordered[j].State) {

		return severity(ordered[i].State) > severity(ordered[j].State)
	}
	return index[ordered[i].ID] < index[ordered[j].ID]
}
