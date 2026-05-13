package posture

import (
	"time"
)

func ingestRepository(repo RepositoryWindow, cutoff time.Time, hasCutoff bool) repositoryIngest {
	// ingestRepository keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if refusal, ok := repositoryPreflightRefusal(repo, cutoff, hasCutoff); ok {
		return refusal
	}
	return ingestRepositoryArtifacts(repo)
}

func repositoryPreflightRefusal(repo RepositoryWindow, cutoff time.Time, hasCutoff bool) (repositoryIngest, bool) {
	// repositoryPreflightRefusal keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	if invalidRepositoryLabels(repo) {
		return refusedInput("unsafe_label", "cannot_verify_input", "", false), true
	}
	if invalidRepositoryInputPaths(repo) {

		return refusedInput("malformed_input", "cannot_verify_input", "", false), true
	}
	if staleRepositoryInput(repo, cutoff, hasCutoff) {
		return refusedInput("stale_input", "stale_input", "", true), true
	}
	return repositoryIngest{}, false
}

func staleRepositoryInput(repo RepositoryWindow, cutoff time.Time, hasCutoff bool) bool {
	return hasCutoff && isStale(repo.InputObservedAt, cutoff)
}

func invalidRepositoryLabels(repo RepositoryWindow) bool {
	return validateRepoLabels(repo) != nil
}

func invalidRepositoryInputPaths(repo RepositoryWindow) bool {
	return validateInputPaths(repo) != nil
}
