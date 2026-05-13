package posture

import (
	"fmt"
	"strings"
)

func validateInputPaths(repo RepositoryWindow) error {
	// validateInputPaths keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	for _, path := range []string{repo.QueryPackResult, repo.ArtifactDigestManifest, repo.PostureSignalManifest} {
		if strings.TrimSpace(path) == "" {
			continue
		}

		if unsafeSelectionPath(path) {
			return fmt.Errorf("unsafe input path")
		}
	}
	return nil
}
