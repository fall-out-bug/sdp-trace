package verifier

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

// ExplainRun returns an operator-facing textual explanation for verify outcomes.
func ExplainRun(runDir string) (string, error) {
	manifestPath := filepath.Join(runDir, "run.json")
	var manifest trace.RunManifest
	if err := trace.ReadJSON(context.Background(), manifestPath, &manifest); err != nil {
		return "", fmt.Errorf("run directory missing: %w", err)
	}
	verification, table, audit, err := VerifyRun(runDir)
	if err != nil {
		return "", err
	}
	lines := []string{
		fmt.Sprintf("run_dir: %s", runDir),
		fmt.Sprintf("run_id: %s", manifest.RunID),
		fmt.Sprintf("contract_id: %s", manifest.ContractID),
		fmt.Sprintf("result: %s", verification.Result),
	}
	lines = appendClosureState(lines, manifest)
	lines = appendAuditIssue(lines, audit)
	lines = appendMissingEvidence(lines, table.Rows)
	lines = appendContractPath(lines, manifest)
	return strings.Join(lines, "\n"), nil
}

func appendClosureState(lines []string, manifest trace.RunManifest) []string {
	if manifest.ClosureState == "" {
		return lines
	}
	return append(lines, fmt.Sprintf("closure_state: %s", manifest.ClosureState))
}

func appendAuditIssue(lines []string, audit *trace.IntegrityAudit) []string {
	if audit == nil || audit.Issue == "" {
		return lines
	}
	lines = append(lines, fmt.Sprintf("integrity_issue: %s", audit.Issue))
	return append(lines, fmt.Sprintf("integrity_reason: %s", audit.Reason))
}

func appendMissingEvidence(lines []string, rows []trace.MissingEvidenceRow) []string {
	if len(rows) == 0 {
		return lines
	}
	lines = append(lines, "missing_evidence:")
	for _, row := range rows {
		lines = append(lines, fmt.Sprintf(" - %s: %s (%s)", row.ExpectedEvent, row.ObservedState, row.Reason))
	}
	return lines
}

func appendContractPath(lines []string, manifest trace.RunManifest) []string {
	if strings.TrimSpace(manifest.ContractPath) == "" {
		return lines
	}
	return append(lines, fmt.Sprintf("contract_path: %s", manifest.ContractPath))
}
