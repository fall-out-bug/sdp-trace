package verifier

import (
	"fmt"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func auditIssueLines(audit *trace.IntegrityAudit) []string {
	// Keep issue and reason adjacent so copied explain output preserves the
	// evidence claim with its replay reason.
	return []string{
		fmt.Sprintf("integrity_issue: %s", audit.Issue),
		fmt.Sprintf("integrity_reason: %s", audit.Reason),
	}
}
