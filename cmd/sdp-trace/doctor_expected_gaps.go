package main

import (
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func expectedEvidenceReferenceGaps(contract trace.Contract) []string {
	missing := make([]string, 0)
	for _, eventType := range contract.RequiredEvents {
		if !knownEventType(eventType) {
			// Required events must map to this binary's local event vocabulary.
			missing = append(missing, "required_events:"+eventType)
		}
	}
	// Evidence-specific gaps are appended after required event gaps so the
	// output separates event vocabulary drift from requirement-shape drift.
	for _, evidence := range contract.RequiredEvidence {
		missing = append(missing, expectedEvidenceGaps(evidence)...)
	}
	return missing
}

func expectedEvidenceGaps(evidence trace.EvidenceRequirement) []string {
	missing := make([]string, 0, 2)
	if strings.TrimSpace(evidence.ID) == "" {
		// Missing evidence IDs make later diagnostics ambiguous.
		missing = append(missing, "required_evidence:<missing_id>")
	}
	if strings.TrimSpace(evidence.EventType) == "" {
		// Missing event type is reported separately from unsupported types.
		return append(missing, "required_evidence:"+evidence.ID+":<missing_event_type>")
	}
	if !knownEventType(evidence.EventType) {
		// Evidence requirements must reference event types this binary can emit.
		missing = append(missing, "required_evidence:"+evidence.ID+":"+evidence.EventType)
	}
	return missing
}
