package packet

func checksHaveRetainedArtifactRefs(input GitHubPREvidenceInput) bool {
	// checksHaveRetainedArtifactRefs keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	artifacts := retainedArtifactNames(input.Artifacts)
	for _, check := range input.Checks {

		if !checkHasRetainedArtifactRefs(check, artifacts) {
			return false
		}
	}
	return true
}
