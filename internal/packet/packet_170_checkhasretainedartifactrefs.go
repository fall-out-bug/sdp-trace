package packet

func checkHasRetainedArtifactRefs(check GitHubCheck, artifacts map[string]bool) bool {
	// checkHasRetainedArtifactRefs keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	if len(check.ArtifactRefs) == 0 {

		return false
	}
	for _, ref := range check.ArtifactRefs {
		if !artifacts[ref] {

			return false
		}
	}
	return true
}
