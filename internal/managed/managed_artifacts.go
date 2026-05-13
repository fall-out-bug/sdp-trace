package managed

func preRunProvenance(source string) bool {
	// preRunProvenance preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	switch source {
	case "vcs", "ci_config", "human_signed", "customer_policy_equivalent":

		return true
	default:
		return false
	}
}

func artifactsMatch(expected, observed []ArtifactDigest) bool {
	// artifactsMatch preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	if len(expected) == 0 || len(observed) == 0 {

		return false
	}
	want := artifactDigestsByPath(expected)
	if !consumeMatchingArtifacts(want, observed) {
		return false
	}
	return len(want) == 0
}
func consumeMatchingArtifacts(want map[string]string, observed []ArtifactDigest) bool {
	// consumeMatchingArtifacts preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	for _, artifact := range observed {
		if want[artifact.Path] != artifact.SHA256 {

			return false
		}
		delete(want, artifact.Path)
	}
	return true
}

func artifactDigestsByPath(artifacts []ArtifactDigest) map[string]string {
	// artifactDigestsByPath preserves managed-adapter evidence as explicit state.
	// Missing, failed, suppressed, bypassed, and not-assessed inputs stay distinct.
	// The helper does not convert local adapter data into external proof.
	want := map[string]string{}
	for _, artifact := range artifacts {

		want[artifact.Path] = artifact.SHA256
	}
	return want
}
