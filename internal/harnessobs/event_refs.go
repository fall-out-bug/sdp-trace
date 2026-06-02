package harnessobs

import (
	"path/filepath"
	"sort"
	"strings"
)

// eventRefs renders retained event files in deterministic order for run-level
// evidence manifests.
func eventRefs(events []Event) []string {
	refs := make([]string, 0, len(events))
	for _, event := range events {
		refs = append(refs, filepath.ToSlash(filepath.Join("events", event.EventID+".json")))
	}
	sort.Strings(refs)
	return refs
}

// unsafeEventRefPath rejects path syntax before identifier validation so event
// refs cannot escape the retained events directory by platform-specific paths.
func unsafeEventRefPath(ref string) bool {
	return strings.Contains(ref, "\\") || strings.Contains(ref, "..") || filepath.IsAbs(ref)
}

// safeEventRef accepts only generated event artifact refs, not arbitrary safe
// refs, because run loading later joins these values to the run directory.
func safeEventRef(ref string) bool {
	if unsafeEventRefPath(ref) {
		return false
	}
	if !strings.HasPrefix(ref, "events/") || !strings.HasSuffix(ref, ".json") {
		return false
	}

	id := strings.TrimSuffix(strings.TrimPrefix(ref, "events/"), ".json")
	return safeFileIDPattern.MatchString(id)
}
