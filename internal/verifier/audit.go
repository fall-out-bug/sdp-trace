package verifier

import "github.com/fall_out_bug/sdp-trace/internal/trace"

// audit records structural verifier issues without treating absent details as
// hidden evidence.
func audit(runID, issue, reason, detailKey, detailValue string) *trace.IntegrityAudit {
	result := &trace.IntegrityAudit{
		RunID:  runID,
		Issue:  issue,
		Reason: reason,
	}
	if detailKey != "" {
		// Details stay optional so audit rows without a stable field do not
		// fabricate empty metadata.
		result.Details = map[string]string{detailKey: detailValue}
	}
	return result
}
