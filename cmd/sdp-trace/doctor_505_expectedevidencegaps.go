package main

import (
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

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
