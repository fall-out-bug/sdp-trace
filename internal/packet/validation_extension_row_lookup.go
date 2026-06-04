package packet

import "sort"

// Custom row fallback is sorted for reproducibility, but it remains fallback
// only because required rows carry the portable gate semantics.
func extensionRowIDForRef(rows map[string]Row, ref string) string {
	for _, id := range sortedRowIDs(rows) {
		// Extension rows are fallback-only; required rows keep precedence even
		// when custom rows cite the same retained evidence ref.
		if !requiredRow(id) && rowHasRef(rows[id], ref) {
			return id
		}
	}
	return ""
}

func sortedRowIDs(rows map[string]Row) []string {
	ids := make([]string, 0, len(rows))
	for id := range rows {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids
}

func rowHasRef(row Row, ref string) bool {
	for _, rowRef := range row.EvidenceRefs {
		// Evidence refs are identifiers, not prefixes; attribution must not infer
		// aliases when choosing a row for contradiction handling.
		if rowRef == ref {
			return true
		}
	}
	return false
}
