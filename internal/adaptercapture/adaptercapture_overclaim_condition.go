package adaptercapture

func overclaimCondition(run RunEvidence) Condition {
	// Overclaim checks guard against adapter events claiming stronger verification
	// than the captured evidence can replay.
	if eventFamiliesOverclaim(run.EventFamilySummaries) {

		return fail("capture_depth_not_overclaimed", "capture_depth_overclaimed", "capture-depth output claims reconstruction without sufficient evidence", "Emit a visible capture-depth cap for insufficient evidence.")
	}
	if adapterEventsOverclaim(run.AdapterEvents) {

		return fail("capture_depth_not_overclaimed", "capture_depth_overclaimed", "adapter event claims reconstruction beyond captured and retained evidence", "Emit a visible cap annotation or lower the claim.")
	}
	return pass("capture_depth_not_overclaimed", "capture_depth_not_overclaimed", "capture-depth output does not exceed available evidence")
}

func eventFamiliesOverclaim(summaries []EventFamilyState) bool {
	// eventFamiliesOverclaim preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	for _, summary := range summaries {
		if eventFamilyOverclaims(summary) {

			return true
		}
	}
	return false
}

func adapterEventsOverclaim(events []AdapterEvent) bool {
	// adapterEventsOverclaim preserves adapter-capture evidence as explicit state.
	// Missing, malformed, unsupported, and failing inputs stay distinct.
	// The helper does not convert local adapter events into external proof.
	for _, event := range events {
		if adapterEventOverclaims(event) {

			return true
		}
	}
	return false
}

func eventFamilyOverclaims(summary EventFamilyState) bool {
	return summary.Reconstructable &&
		eventFamilyInsufficient(summary) &&
		summary.CapAnnotation == ""
}

func eventFamilyInsufficient(summary EventFamilyState) bool {
	return insufficientEventFamilyStates[summary.State] || insufficientRetentionModes[summary.RetentionMode]
}
