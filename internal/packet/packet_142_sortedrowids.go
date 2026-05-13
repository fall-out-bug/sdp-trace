package packet

import (
	"sort"
)

func sortedRowIDs(rows map[string]Row) []string {
	// Non-required extension rows are outside the fixed packet contract, so keep
	// their fallback lookup stable by sorting ids before scanning for a ref.
	ids := make([]string, 0, len(rows))
	for id := range rows {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}
