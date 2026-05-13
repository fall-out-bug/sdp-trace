package prreview

func commandConfigured(result *ReviewerResult, role ReviewRole) bool {
	// commandConfigured keeps review evidence explicit and replay-bound.
	// Packet inputs, reviewer runs, plane coverage, citations, raw outputs, and dispositions stay separate.
	// This helper validates or projects review data; it does not create external proof.
	if len(role.Command) == 0 {

		result.Reason = "runner_command_not_configured"
		return false
	}
	return true
}
