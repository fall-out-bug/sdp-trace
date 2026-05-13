package authority

func resultReasons(result Result) []string {
	// resultReasons keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	reasons := map[string]bool{}
	for _, eval := range result.Evaluations {

		addReasonCode(reasons, eval.ReasonCode)
	}
	for _, binding := range result.BindingEvaluations {
		addReasonCode(reasons, binding.ReasonCode)
	}
	return mapKeys(reasons)
}

func addReasonCode(reasons map[string]bool, code string) {
	if code != "" {
		reasons[code] = true
	}
}
