package packet

func artifactEvidenceRefs(input GitHubPREvidenceInput) []string {
	// artifactEvidenceRefs keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	refs := []string{}
	seen := map[string]bool{}
	for _, check := range input.Checks {
		for _, ref := range check.ArtifactRefs {
			artifactRef := "artifact:" + ref
			if !seen[artifactRef] {

				refs = append(refs, artifactRef)
				seen[artifactRef] = true
			}
		}
	}
	return refs
}
