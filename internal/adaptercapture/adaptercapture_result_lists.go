package adaptercapture

import (
	"sort"
)

func reasons(conditions []Condition) []string {
	// Reasons are derived from condition states so human text follows evidence and
	// does not become independent authority.
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
	// Next actions point to concrete missing adapter evidence needed for replayable
	// capture assessment.
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
	// addNextAction preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	if condition.State != StatePass && condition.NextAction != "" {

		set[condition.NextAction] = true
	}
}
