package packet

// A check has retained artifact evidence only when every referenced artifact is
// retained. Missing refs keep verification partial rather than pass.
func checksHaveRetainedArtifactRefs(input GitHubPREvidenceInput) bool {
	artifacts := retainedArtifactNames(input.Artifacts)
	for _, check := range input.Checks {
		if !checkHasRetainedArtifactRefs(check, artifacts) {
			return false
		}
	}
	return true
}

func checkHasRetainedArtifactRefs(check GitHubCheck, artifacts map[string]bool) bool {
	if len(check.ArtifactRefs) == 0 {
		return false
	}
	return allArtifactRefsRetained(check.ArtifactRefs, artifacts)
}

// Every check ref must resolve to a retained artifact name; one missing ref is
// enough to keep verification partial.
func allArtifactRefsRetained(refs []string, artifacts map[string]bool) bool {
	for _, ref := range refs {
		if artifactMissing(ref, artifacts) {
			return false
		}
	}
	return true
}

// Missing map entries include both unknown artifacts and explicitly
// not-retained artifacts filtered out by retainedArtifactNames.
func artifactMissing(ref string, artifacts map[string]bool) bool {
	return !artifacts[ref]
}
