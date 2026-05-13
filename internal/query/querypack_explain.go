package query

import (
	"fmt"
	"sort"
	"strings"
)

func sortedQueryRows(rows []QueryPackRow) []QueryPackRow {
	// sortedQueryRows keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	sorted := append([]QueryPackRow(nil), rows...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].ID < sorted[j].ID })
	return sorted
}

func explainQueryRow(queryName string, row QueryPackRow) string {
	// explainQueryRow keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	parts := []string{queryName, row.ID, row.EvidenceState, row.EvidenceFamily}
	parts = append(parts, "source_ref="+row.SourceRef)
	parts = appendOptionalPart(parts, "source_condition_id", row.SourceConditionID)
	parts = appendOptionalPart(parts, "source_condition_state", row.SourceConditionState)
	if row.Reconstructable != nil {
		parts = append(parts, fmt.Sprintf("reconstructable=%t", *row.Reconstructable))
	}
	parts = appendOptionalPart(parts, "gap", row.EvidenceGap)
	return strings.Join(parts, " ")
}

func appendOptionalPart(parts []string, key, value string) []string {
	// appendOptionalPart keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if value == "" {
		return parts
	}
	return append(parts, key+"="+value)
}
