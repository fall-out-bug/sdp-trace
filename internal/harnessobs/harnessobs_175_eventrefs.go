package harnessobs

import (
	"path/filepath"

	"sort"
)

func eventRefs(events []Event) []string {
	// eventRefs keeps harness observation evidence explicit and replay-bound.
	// Source profiles, raw events, path safety, digests, validation, and command models stay separate.
	// This helper renders or aggregates harness evidence; it does not create external proof.
	refs := make([]string, 0, len(events))
	for _, event := range events {

		refs = append(refs, filepath.ToSlash(filepath.Join("events", event.EventID+".json")))
	}
	sort.Strings(refs)
	return refs
}
