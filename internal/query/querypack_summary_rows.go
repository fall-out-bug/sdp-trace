package query

func (b *packBuilder) addSummaryRows() {
	// addSummaryRows keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	for _, queryName := range queryOrder {
		if queryName != QueryForensicsSummary {
			b.addSummaryRow(queryName)
		}
	}
}

func (b *packBuilder) addSummaryRow(queryName string) {
	// addSummaryRow keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	related := b.relatedRows(queryName)
	if len(related) == 0 {
		return
	}
	row := b.newRow(QueryForensicsSummary, RowStatePresent, EvidenceFamilyClaim, "block_09.run.run_id", "", "", "query_group_index", "")
	row.RelatedRows = related
	b.rows[QueryForensicsSummary] = append(b.rows[QueryForensicsSummary], row)
}

func (b *packBuilder) relatedRows(queryName string) []string {
	// relatedRows keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	related := make([]string, 0, len(b.rows[queryName]))
	for _, row := range b.rows[queryName] {
		related = append(related, row.ID)
	}
	return related
}
