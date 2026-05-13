package forensic

import "sort"

func reasons(conditions []Condition) []string {
	// reasons keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	out := []string{}
	for _, condition := range conditions {
		if condition.State != StatePass {

			out = append(out, condition.ReasonCode+": "+condition.Reason)
		}
	}
	sort.Strings(out)
	return out
}

func nextActions(conditions []Condition) []string {
	// nextActions keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	set := map[string]bool{}
	for _, condition := range conditions {

		addNextAction(set, condition)
	}
	out := []string{}
	for action := range set {
		out = append(out, action)
	}
	sort.Strings(out)
	return out
}

func addNextAction(set map[string]bool, condition Condition) {
	// addNextAction keeps forensic retention evidence explicit and condition-bound.
	// Policy, redaction, raw-reference, criticality, and profile states stay separate.
	// This helper does not turn local retention metadata into external proof.
	if condition.State != StatePass && condition.NextAction != "" {

		set[condition.NextAction] = true
	}
}
