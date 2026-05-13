package query

import "fmt"

func (b *packBuilder) addUnverifiedClaimRows() {
	b.addUnverifiedClaimsFor(QueryForensicsRedactions)
	b.addUnverifiedClaimsFor(QueryForensicsCaptureDepth)
}

func (b *packBuilder) addUnverifiedClaimsFor(queryName string) {
	// addUnverifiedClaimsFor keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	for _, row := range append([]QueryPackRow{}, b.rows[queryName]...) {
		if row.EvidenceState != RowStatePresent {
			b.addReferencedClaim(row)
		}
	}
}

func (b *packBuilder) addReferencedClaim(source QueryPackRow) {
	// addReferencedClaim keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	row := b.newRow(QueryForensicsUnverifiedClaims, source.EvidenceState, EvidenceFamilyClaim, source.SourceRef, source.SourceConditionID, source.SourceConditionState, source.ReasonCode, source.EvidenceGap)
	row.RelatedRows = []string{source.ID}
	row.Reconstructable = source.Reconstructable
	b.rows[QueryForensicsUnverifiedClaims] = append(b.rows[QueryForensicsUnverifiedClaims], row)
}

func (b *packBuilder) nextRowID(queryName string) string {
	b.counters[queryName]++
	return fmt.Sprintf("%s.%04d", queryShortName(queryName), b.counters[queryName])
}
