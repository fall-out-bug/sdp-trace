package harnessobs

import (
	"path/filepath"
	"sort"
)

func eventRefs(events []Event) []string {
	refs := make([]string, 0, len(events))
	for _, event := range events {
		refs = append(refs, filepath.ToSlash(filepath.Join("events", event.EventID+".json")))
	}
	sort.Strings(refs)
	return refs
}
