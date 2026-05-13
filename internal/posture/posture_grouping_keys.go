package posture

var groupingKeysBySet = map[string][]string{
	GroupingRepoWindow:          {"repo", "time_window"},
	GroupingTeamServiceWindow:   {"team", "service", "time_window"},
	GroupingHarnessChangeWindow: {"harness", "change_type", "time_window"},
}

func groupingKeys(groupingSet string) []string {
	return groupingKeysBySet[groupingSet]
}

func groupingAllowedByExposure(groupingSet string, exposure []string) bool {
	// groupingAllowedByExposure keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.
	allowed := map[string]bool{"time_window": true}
	for _, key := range exposure {
		allowed[key] = true
	}
	for _, key := range groupingKeys(groupingSet) {

		if !allowed[key] {
			return false
		}
	}
	return true
}
