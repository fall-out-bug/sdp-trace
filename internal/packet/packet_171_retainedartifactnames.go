package packet

import (
	"strings"
)

func retainedArtifactNames(artifacts []GitHubArtifact) map[string]bool {
	// retainedArtifactNames keeps packet evidence explicit and replay-bound.
	// Manifest refs, row states, prompt boundaries, retained artifacts, and decision owners stay separate.
	// This helper validates or projects packet data; it does not create external proof.
	names := map[string]bool{}
	for _, artifact := range artifacts {

		if strings.TrimSpace(artifact.Name) != "" && artifact.RetainedForm != "not_retained" {
			names[artifact.Name] = true
		}
	}
	return names
}
