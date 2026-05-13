package query

import "fmt"

func (b *packBuilder) addTimelineRows() {
	b.addRunTimelineRows()
	b.addOptionalTimelineRows()
}

func (b *packBuilder) addRunTimelineRows() {
	// addRunTimelineRows keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if len(b.inputs.run.EventRefs) == 0 {
		b.addRow(QueryForensicsTimeline, RowStatePresent, EvidenceFamilyRunChain, "block_09.run.run_id", "", "", "run_timeline_available", "")
		return
	}
	for i, event := range b.inputs.run.EventRefs {
		family := familyForEvent(event.EventType)
		sourceRef := fmt.Sprintf("block_09.event.%s.e%04d", family, i+1)
		b.addRow(QueryForensicsTimeline, RowStatePresent, family, sourceRef, "", "", "timeline_event_present", "")
	}
}

func (b *packBuilder) addOptionalTimelineRows() {
	// Optional block rows preserve absence or malformed state without hiding run timeline data.
	b.addOptionalTimelineRow(b.inputs.forensicPresent, b.inputs.forensicErr, EvidenceFamilyRetention, "block_18", "missing_optional_block_18_forensic_retention_result")
	b.addOptionalTimelineRow(b.inputs.adapterPresent, b.inputs.adapterErr, EvidenceFamilyAdapterCapture, "block_19", "missing_optional_block_19_adapter_capture_result")
}

func (b *packBuilder) addOptionalTimelineRow(present bool, inputErr error, family, block, missingReason string) {
	// addOptionalTimelineRow keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	if !present {
		b.addRow(QueryForensicsTimeline, RowStateNotAssessed, family, block+".condition.missing", "", "", missingReason, family)
		return
	}
	if inputErr != nil {
		b.addRow(QueryForensicsTimeline, RowStateCannotVerify, EvidenceFamilyInputArtifact, block+".condition.malformed", "", "", "unreadable_or_malformed_input_artifact", EvidenceFamilyInputArtifact)
	}
}

func (b *packBuilder) addMalformedRequiredInputRows() {
	// addMalformedRequiredInputRows keeps query-pack rows source-bound to replayed evidence artifacts.
	// Missing, malformed, redacted, retained, and adapter evidence stay separate.
	// This helper renders derived query rows; it does not create a new verdict.
	for _, queryName := range queryOrder {
		b.addRow(queryName, RowStateCannotVerify, EvidenceFamilyInputArtifact, "block_09.run.malformed", "", "", "unreadable_or_malformed_input_artifact", EvidenceFamilyInputArtifact)
	}
}
