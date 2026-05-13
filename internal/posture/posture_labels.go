package posture

import (
	"fmt"
)

func validateRepoLabels(repo RepositoryWindow) error {
	// validateRepoLabels keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	if unsafeInputLabel(repo) {
		return fmt.Errorf("unsafe label")
	}

	for _, value := range repositoryOutputLabels(repo) {

		if !safeLabel(value) {
			return fmt.Errorf("unsafe label")
		}
	}
	return nil
}

func repositoryOutputLabels(repo RepositoryWindow) []string {

	return []string{repo.Repo, repo.Team, repo.Service, repo.Harness, repo.ChangeType}
}

func unsafeInputLabel(repo RepositoryWindow) bool {
	return !safeLabel(repo.InputID) || !safeLabel(repo.TimeWindow)
}

// safeLabel enforces the input trust boundary for identifier syntax.
// Unsafe output keywords or credential/token terminology crosses the injection boundary.
func safeLabel(value string) bool {
	// safeLabel keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if !safeLabelPattern.MatchString(value) {
		return false
	}
	return !unsafeLabel(value)
}
