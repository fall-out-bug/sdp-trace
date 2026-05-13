package witness

import (
	"strings"
)

func runIDMatches(runsRoot, witnessRunID string) bool {
	if strings.TrimSpace(witnessRunID) == "" {
		// Empty run IDs cannot bind external or envelope evidence to a run.
		return false
	}
	runIDs, err := runIDsFromRoot(runsRoot)
	if err != nil || len(runIDs) == 0 {
		return false
	}
	return containsRunID(runIDs, witnessRunID)
}

func containsRunID(runIDs []string, witnessRunID string) bool {
	// Exact run ID matching prevents envelope or freshness evidence from binding
	// to a sibling run.
	for _, runID := range runIDs {
		if runID == witnessRunID {
			return true
		}
	}
	return false
}
