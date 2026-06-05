package packet

// Artifact evidence refs are attached to pass verification rows only after
// GitHub check refs are normalized into packet manifest refs.
func artifactEvidenceRefs(input GitHubPREvidenceInput) []string {
	refs := []string{}
	seen := map[string]bool{}
	for _, check := range input.Checks {
		refs = appendArtifactEvidenceRefs(refs, seen, check.ArtifactRefs)
	}
	return refs
}

// Artifact refs are de-duplicated in first-seen order so packet row evidence is
// stable while preserving the input's check ordering.
func appendArtifactEvidenceRefs(refs []string, seen map[string]bool, artifactRefs []string) []string {
	for _, ref := range artifactRefs {
		artifactRef := "artifact:" + ref
		if !seen[artifactRef] {
			refs = append(refs, artifactRef)
			seen[artifactRef] = true
		}
	}
	return refs
}
