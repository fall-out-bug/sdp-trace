package authority

func evidenceResolutionIndex(input EvidenceResolution) map[string]string {
	// evidenceResolutionIndex keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.

	out := map[string]string{}
	addEvidenceResolution(out, input.ResolvedExternalRefs, "resolved")
	addEvidenceResolution(out, input.InaccessibleRefs, "inaccessible")
	addEvidenceResolution(out, input.MalformedRefs, "malformed")
	addEvidenceResolution(out, input.StaleRefs, "stale")
	return out
}

func addEvidenceResolution(out map[string]string, refs []string, state string) {
	// addEvidenceResolution keeps authority envelope evidence explicit and bounded to observed inputs.
	// Missing, conflicting, unsafe, and externally unresolved proof stays distinguishable.
	// This helper does not turn local policy data into external trust.
	for _, ref := range refs {

		out[ref] = state
	}
}
