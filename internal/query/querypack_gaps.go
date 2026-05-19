package query

import "sort"

func (b *packBuilder) addGapRows() {
	b.addVerifierGapRows()
	b.addForensicGapRows()
	b.addAdapterGapRows()
}

func (b *packBuilder) addVerifierGapRows() {
	// addVerifierGapRows keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	for _, key := range sortedVerifierStateKeys(b.inputs.run.VerifierStates) {
		b.addVerifierGapRow(key, b.inputs.run.VerifierStates[key])
	}
}

func sortedVerifierStateKeys(states map[string]verifierState) []string {
	// sortedVerifierStateKeys keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	keys := make([]string, 0, len(states))
	for key := range states {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (b *packBuilder) addVerifierGapRow(key string, state verifierState) {
	// addVerifierGapRow keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	rowState := mapSourceState(state.State)
	if rowState == RowStatePresent {
		return
	}
	family := familyForVerifierState(key)
	b.addRow(QueryForensicsGaps, rowState, family, "block_09.run."+safeToken(key), safeToken(key), family)
}

func (b *packBuilder) addForensicGapRows() {
	// addForensicGapRows keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if !b.inputs.forensicPresent {
		b.addRow(QueryForensicsGaps, RowStateNotAssessed, EvidenceFamilyRetention, "block_18.condition.missing", "missing_optional_block_18_forensic_retention_result", "retention")
	} else if b.inputs.forensicErr != nil {
		b.addRow(QueryForensicsGaps, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_18.condition.malformed", "unreadable_or_malformed_input_artifact", EvidenceFamilyInputArtifact)
	}
}

func (b *packBuilder) addAdapterGapRows() {
	// addAdapterGapRows keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if !b.inputs.adapterPresent {
		b.addRow(QueryForensicsGaps, RowStateNotAssessed, EvidenceFamilyAdapterCapture, "block_19.condition.missing", "missing_optional_block_19_adapter_capture_result", "adapter_capture")
	} else if b.inputs.adapterErr != nil {
		b.addRow(QueryForensicsGaps, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_19.condition.malformed", "unreadable_or_malformed_input_artifact", EvidenceFamilyInputArtifact)
	}
}
