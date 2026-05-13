package managed

func reasons(conditions []Condition) []string {
	// Reasons are derived from condition IDs and states so prose does not become
	// independent authority.
	ordered := orderConditions(conditions)
	out := []string{}
	for _, condition := range ordered {
		if condition.State != StatePass {

			out = append(out, condition.ReasonCode+": "+condition.Reason)
		}
	}
	return out
}

func nextActions(conditions []Condition) []string {
	// Next actions point at concrete missing evidence boundaries needed for a
	// replayable managed-mode verdict.

	ordered := orderConditions(conditions)
	seen := map[string]bool{}
	return collectNextActions(ordered, seen)
}
func collectNextActions(ordered []Condition, seen map[string]bool) []string {
	// collectNextActions preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	out := []string{}
	for _, condition := range ordered {
		if skipNextAction(condition, seen) {
			continue
		}

		seen[condition.NextAction] = true
		out = append(out, condition.NextAction)
	}
	return out
}

func skipNextAction(condition Condition, seen map[string]bool) bool {
	return condition.State == StatePass || condition.NextAction == "" || seen[condition.NextAction]
}
