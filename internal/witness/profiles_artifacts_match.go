package witness

func artifactSetsMatch(expected, current []ArtifactDigest) bool {
	if len(expected) != len(current) {
		// Count differences are binding failures even if all shared paths match.
		return false
	}
	byPath := artifactDigestsByPath(current)
	for _, artifact := range expected {
		if byPath[artifact.Path] != artifact.SHA256 {
			return false
		}
	}
	return true
}

func artifactDigestsByPath(artifacts []ArtifactDigest) map[string]string {
	byPath := map[string]string{}
	for _, artifact := range artifacts {
		// Path is the binding key because witness artifacts are compared against
		// the selected run/report file set, not just a bag of digests.
		byPath[artifact.Path] = artifact.SHA256
	}
	return byPath
}
