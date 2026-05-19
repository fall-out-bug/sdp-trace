package query

func (b *packBuilder) addRedactionRows() {
	// addRedactionRows keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if !b.inputs.forensicPresent {
		b.addRow(QueryForensicsRedactions, RowStateCannotVerify, EvidenceFamilyRedaction, "block_18.condition.missing", "missing_block_18_forensic_retention_result", "redaction")
		return
	}
	if b.inputs.forensicErr != nil {
		b.addRow(QueryForensicsRedactions, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_18.condition.malformed", "unreadable_or_malformed_input_artifact", "input_artifact")
		return
	}
	for _, condition := range b.inputs.forensic.ForensicConditions {
		family := familyForForensicCondition(condition.ID)
		row := b.rowFromCondition(QueryForensicsRedactions, family, "block_18.condition."+safeToken(condition.ID), condition)
		b.rows[QueryForensicsRedactions] = append(b.rows[QueryForensicsRedactions], row)
	}
}

func (b *packBuilder) addCaptureRows() {
	// addCaptureRows keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if !b.inputs.adapterPresent {
		b.addMissingAdapterCaptureRow()
		return
	}
	if b.inputs.adapterErr != nil {
		b.addMalformedAdapterCaptureRow()
		return
	}
	for _, condition := range b.inputs.adapter.AdapterCaptureConditions {
		b.addAdapterCaptureConditionRow(condition)
	}
}

func (b *packBuilder) addMissingAdapterCaptureRow() {
	b.addRow(QueryForensicsCaptureDepth, RowStateCannotVerify, EvidenceFamilyAdapterCapture, "block_19.condition.missing", "missing_block_19_adapter_capture_result", "adapter_capture")
}

func (b *packBuilder) addMalformedAdapterCaptureRow() {
	b.addRow(QueryForensicsCaptureDepth, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_19.condition.malformed", "unreadable_or_malformed_input_artifact", "input_artifact")
}

func (b *packBuilder) addAdapterCaptureConditionRow(condition assessmentCondition) {
	// addAdapterCaptureConditionRow keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	family := familyForAdapterCondition(condition.ID)
	row := b.rowFromCondition(QueryForensicsCaptureDepth, family, "block_19.condition."+safeToken(condition.ID), condition)
	b.rows[QueryForensicsCaptureDepth] = append(b.rows[QueryForensicsCaptureDepth], row)
}
