package main

import (
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
