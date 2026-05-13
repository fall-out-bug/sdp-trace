package posture

import (
	"strings"
)

func dimensionKey(repo RepositoryWindow, keys []string) (string, map[string]string) {
	// dimensionKey keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	values := dimensionValues(repo)
	dims := map[string]string{}
	var parts []string
	for _, key := range keys {
		dims[key] = values[key]
		if key == "time_window" {

			continue
		}
		parts = append(parts, key+"="+values[key])
	}
	return strings.Join(parts, "|"), dims
}

func dimensionValues(repo RepositoryWindow) map[string]string {
	// dimensionValues keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	return map[string]string{
		"repo":        repo.Repo,
		"team":        repo.Team,
		"service":     repo.Service,
		"harness":     repo.Harness,
		"change_type": repo.ChangeType,
		"time_window": repo.TimeWindow,
	}
}
