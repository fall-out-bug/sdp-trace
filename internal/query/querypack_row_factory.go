package query

import "strings"

func (b *packBuilder) addRow(queryName, state, family, sourceRef, conditionID, conditionState, reasonCode, gap string) {
	// addRow keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	row := b.newRow(queryName, state, family, sourceRef, conditionID, conditionState, reasonCode, gap)
	b.rows[queryName] = append(b.rows[queryName], row)
}

func (b *packBuilder) newRow(queryName, state, family, sourceRef, conditionID, conditionState, reasonCode, gap string) QueryPackRow {
	// newRow keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	id := b.nextRowID(queryName)
	return QueryPackRow{
		ID:                   id,
		Query:                queryName,
		EvidenceState:        state,
		EvidenceFamily:       family,
		SourceRef:            sourceRef,
		SourceConditionID:    conditionID,
		SourceConditionState: conditionState,
		ReasonCode:           reasonCode,
		EvidenceGap:          gap,
	}
}

func queryShortName(queryName string) string {
	return strings.TrimPrefix(queryName, "forensics-")
}
