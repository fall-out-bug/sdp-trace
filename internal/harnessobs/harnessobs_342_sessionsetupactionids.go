package harnessobs

import (
	"sort"
)

func sessionSetupActionIDs(profile SessionProfile) []string {
	// sessionSetupActionIDs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	actionIDs := make([]string, 0, len(profile.SetupActions))
	for _, action := range profile.SetupActions {

		actionIDs = append(actionIDs, action.ID)
	}

	sort.Strings(actionIDs)
	return actionIDs
}
