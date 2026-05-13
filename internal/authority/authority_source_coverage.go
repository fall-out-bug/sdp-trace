package authority

func missingAttributes(eval AuthorityEvaluation) []string {
	// missingAttributes keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	var missing []string
	if eval.ActorAttribution == AttributionNotAssessed {

		missing = append(missing, "actor")
	}
	if eval.ToolAttribution == AttributionNotAssessed {

		missing = append(missing, "tool")
	}
	if eval.ModelAttribution == AttributionNotAssessed {

		missing = append(missing, "model")
	}
	return missing
}

func sourceCoverage(actions []ObservedAction) []string {
	// sourceCoverage keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	var sources []string
	for _, action := range actions {

		sources = append(sources, action.SourceType)
	}
	return uniqueStrings(sources)
}

func safeRefs(refs []string) []string {
	// safeRefs keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	out := make([]string, 0, len(refs))
	for _, ref := range refs {
		if evidenceRefPattern.MatchString(ref) && !unsafeRefPattern.MatchString(ref) {

			out = append(out, ref)
		}
	}
	return out
}
