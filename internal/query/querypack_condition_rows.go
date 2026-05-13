package query

func (b *packBuilder) rowFromCondition(queryName, family, sourceRef string, condition assessmentCondition) QueryPackRow {
	// rowFromCondition keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if criticalEvidenceRetentionCap(condition) {
		state := RowStateRetentionLimited
		row := b.newRow(queryName, state, family, sourceRef, condition.ID, condition.State, "digest_only_not_reconstructable", EvidenceFamilyRetention)
		row.Reconstructable = falsePointer()
		return row
	}
	state := mapSourceState(condition.State)
	gap := gapForConditionState(state, family)
	row := b.newRow(queryName, state, family, sourceRef, condition.ID, condition.State, condition.ReasonCode, gap)
	row.Reconstructable = reconstructableForCondition(condition)
	return row
}

func criticalEvidenceRetentionCap(condition assessmentCondition) bool {
	// criticalEvidenceRetentionCap keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	return condition.ID == "critical_evidence_reconstructable" &&
		(condition.CappedToRetentionMode != "" || condition.ReasonCode == "critical_evidence_digest_only")
}

func gapForConditionState(state, family string) string {
	// gapForConditionState keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if state == RowStatePresent || state == RowStateIssueObserved {
		return ""
	}
	return family
}

func reconstructableForCondition(condition assessmentCondition) *bool {
	// reconstructableForCondition keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if condition.State == RowStateRetentionLimited || condition.CappedToRetentionMode != "" {
		return falsePointer()
	}
	return nil
}

func falsePointer() *bool {
	falseValue := false
	return &falseValue
}
