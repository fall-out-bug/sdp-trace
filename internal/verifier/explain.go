package verifier

import (
	"context"
	"fmt"
	"os"
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
	if manifest.ClosureState != "" {
		lines = append(lines, fmt.Sprintf("closure_state: %s", manifest.ClosureState))
	}
	if audit != nil && audit.Issue != "" {
		lines = append(lines, fmt.Sprintf("integrity_issue: %s", audit.Issue))
		lines = append(lines, fmt.Sprintf("integrity_reason: %s", audit.Reason))
	}
	if len(table.Rows) > 0 {
		lines = append(lines, "missing_evidence:")
		for _, row := range table.Rows {
			lines = append(lines, fmt.Sprintf(" - %s: %s (%s)", row.ExpectedEvent, row.ObservedState, row.Reason))
		}
	}
	if strings.TrimSpace(manifest.ContractPath) != "" {
		lines = append(lines, fmt.Sprintf("contract_path: %s", manifest.ContractPath))
	}
	return strings.Join(lines, "\n"), nil
}

func runVerifierArtifact(runDir string) (string, error) {
	resultPath := filepath.Join(runDir, "verifier", "verifier-result.json")
	_, err := os.Stat(resultPath)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(resultPath)
	if err != nil {
		return "", err
	}
	var result trace.VerifierResult
	if err := trace.ReadJSON(context.Background(), resultPath, &result); err != nil {
		_ = data
		return "", err
	}
	return string(data), nil
}
