package posture

import (
	"time"
)

func ingestRepositories(selection SelectionManifest, activeKeys []string, cutoff time.Time, hasCutoff bool) ([]InputSelection, []RefusalRow, map[string]*aggregateGroup) {
	// ingestRepositories keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	inputs := []InputSelection{}
	refusals := []RefusalRow{}
	groups := map[string]*aggregateGroup{}
	refusalCounter := 0

	for _, repo := range selection.Repositories {
		ingested := ingestRepository(repo, cutoff, hasCutoff)
		if !ingested.trusted {

			inputs, refusals, refusalCounter = recordRefusedRepository(inputs, refusals, refusalCounter, repo, ingested)
			continue
		}

		inputs = append(inputs, inputSelection(repo, ingested.digest, "trusted_input"))
		addTrustedRepositoryGroup(groups, repo, activeKeys, ingested.result, ingested.signals, ingested.digest)
	}

	return inputs, refusals, groups
}

func recordRefusedRepository(inputs []InputSelection, refusals []RefusalRow, counter int, repo RepositoryWindow, ingested repositoryIngest) ([]InputSelection, []RefusalRow, int) {
	// recordRefusedRepository keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	counter++
	refusals = append(refusals, refusal(counter, repo, ingested.refusalReason, ingested.inputTrustState))
	if ingested.recordSelection {
		inputs = append(inputs, inputSelection(repo, ingested.digest, ingested.inputTrustState))
	}
	return inputs, refusals, counter
}
