package posture

import (
	"sort"
)

func metricEvidenceRefs(group *aggregateGroup) ([]string, string) {
	// metricEvidenceRefs keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	sourceRefs := sortedStringsCopy(group.inputRefs)
	digests := sortedStringsCopy(group.digests)
	return sourceRefs, digestSetHash(digests)
}

func sortedStringsCopy(values []string) []string {
	// sortedStringsCopy keeps posture export evidence explicit and source-bound.
	// Selection, digest, signal, metric, movement, refusal, and safety states stay separate.
	// This helper aggregates replayed query-pack data; it does not create new proof.

	out := append([]string(nil), values...)
	sort.Strings(out)
	return out
}
