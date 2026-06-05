package packet

import "strings"

func retainedArtifactNames(artifacts []GitHubArtifact) map[string]bool {
	names := map[string]bool{}
	for _, artifact := range artifacts {
		if retainedArtifactNamed(artifact) {
			names[artifact.Name] = true
		}
	}
	return names
}

// Empty names and not_retained artifacts are intentionally excluded from the
// retained set so they cannot support a pass verification row.
func retainedArtifactNamed(artifact GitHubArtifact) bool {
	return strings.TrimSpace(artifact.Name) != "" && artifact.RetainedForm != "not_retained"
}
