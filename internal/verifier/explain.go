package verifier

import (
	"fmt"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// ExplainRun returns an operator-facing textual explanation for verify outcomes.
func ExplainRun(runDir string) (string, error) {
	// Explain output is intentionally derived from a live verifier replay, not
	// from checked-in verifier artifacts. That keeps operator text subordinate
	// to machine replay state.
	manifest, err := readExplainManifest(runDir)
	if err != nil {
		return "", err
	}
	verification, table, audit, err := VerifyRun(runDir)
	if err != nil {
		return "", err
	}
	return strings.Join(explainLines(runDir, manifest, verification, table, audit), "\n"), nil
}

func explainLines(runDir string, manifest trace.RunManifest, verification trace.VerifierResult, table trace.MissingEvidenceTable, audit *trace.IntegrityAudit) []string {
	// Keep the leading rows stable for scripts and reviews that scan explain
	// output before optional audit or evidence sections.
	lines := explainHeaderLines(runDir, manifest, verification)
	lines = appendClosureState(lines, manifest)
	lines = appendAuditIssue(lines, audit)
	lines = appendMissingEvidence(lines, table.Rows)
	lines = appendContractPath(lines, manifest)
	return lines
}

func appendAuditIssue(lines []string, audit *trace.IntegrityAudit) []string {
	if audit == nil || audit.Issue == "" {
		return lines
	}
	// Audit rows are structural evidence from replay; omit them entirely when
	// no concrete issue was generated.
	return append(lines, auditIssueLines(audit)...)
}

func appendContractPath(lines []string, manifest trace.RunManifest) []string {
	if strings.TrimSpace(manifest.ContractPath) == "" {
		// Contract path is contextual metadata; missing context must not look
		// like a verified contract binding.
		return lines
	}
	return append(lines, fmt.Sprintf("contract_path: %s", manifest.ContractPath))
}
