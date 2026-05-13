package posture

import (
	"github.com/fall_out_bug/sdp-trace/internal/query"
)

func addTrustedRepositoryGroup(groups map[string]*aggregateGroup, repo RepositoryWindow, activeKeys []string, result query.QueryPackResult, signals map[string]PostureSignal, digest string) {
	// addTrustedRepositoryGroup keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	group := trustedAggregateGroup(groups, repo, activeKeys)
	group.rows = append(group.rows, flattenRows(result)...)
	group.inputRefs = append(group.inputRefs, repo.InputID)
	group.digests = append(group.digests, digest)
	group.trustStates["trusted_input"]++

	for rowRef, signal := range signals {
		group.signals[rowRef] = signal
	}
}

func trustedAggregateGroup(groups map[string]*aggregateGroup, repo RepositoryWindow, activeKeys []string) *aggregateGroup {
	// trustedAggregateGroup keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	key, dims := dimensionKey(repo, activeKeys)
	groupKey := repo.TimeWindow + "|" + key
	group := groups[groupKey]
	if group == nil {

		group = newAggregateGroup(repo, key, dims)
		groups[groupKey] = group
	}
	return group
}

func newAggregateGroup(repo RepositoryWindow, key string, dims map[string]string) *aggregateGroup {
	// newAggregateGroup keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	return &aggregateGroup{
		dimensions:   dims,
		dimensionKey: key,
		window:       repo.TimeWindow,
		trustStates:  map[string]int{},
		signals:      map[string]PostureSignal{},
	}
}
