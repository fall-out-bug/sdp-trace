package demo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"github.com/fall_out_bug/sdp-trace/internal/verifier"
)

const (
	GatePass         = "pass"
	GateFail         = "fail"
	GateCannotVerify = "cannot_verify"
)

type RunRow struct {
	Name          string                `json:"name"`
	RunID         string                `json:"run_id"`
	Kind          string                `json:"kind"`
	KindReason    string                `json:"kind_reason"`
	Command       string                `json:"command"`
	WrapperName   string                `json:"wrapper_name,omitempty"`
	ExitCode      *int                  `json:"exit_code"`
	ClosureState  string                `json:"closure_state"`
	Result        trace.VerifierVerdict `json:"result"`
	TrustScope    trace.TrustScope      `json:"trust_scope"`
	Completeness  trace.Completeness    `json:"completeness"`
	Replayability trace.Replayability   `json:"replayability"`
	StdoutDigest  string                `json:"stdout_digest"`
	StderrDigest  string                `json:"stderr_digest"`
	Reason        string                `json:"reason,omitempty"`
}

type Summary struct {
	GeneratedAt       string   `json:"generated_at"`
	RunCount          int      `json:"run_count"`
	ObservedCount     int      `json:"observed_count"`
	FailedCount       int      `json:"failed_count"`
	CannotVerifyCount int      `json:"cannot_verify_count"`
	NotAssessedCount  int      `json:"not_assessed_count"`
	TrustScope        string   `json:"trust_scope"`
	AuditGrade        bool     `json:"audit_grade"`
	AuditGradeReason  string   `json:"audit_grade_reason"`
	Runs              []RunRow `json:"runs"`
}

type EvidenceTable struct {
	Runs []RunRow `json:"runs"`
}

type MissingTelemetry struct {
	MissingAuditEvidence   []string `json:"missing_audit_evidence"`
	MissingHarnessEvidence []string `json:"missing_harness_evidence"`
	Notes                  []string `json:"notes"`
}

type GateResult struct {
	LocalGate            string          `json:"local_gate"`
	CIWitnessGate        string          `json:"ci_witness_gate"`
	AuditGradeGate       string          `json:"audit_grade_gate"`
	Reasons              []string        `json:"reasons"`
	RequiredEvidence     []string        `json:"required_evidence"`
	ObservedEvidence     []string        `json:"observed_evidence"`
	GateConditions       []string        `json:"gate_conditions"`
	MissingAuditEvidence []string        `json:"missing_audit_evidence"`
	Witness              *WitnessSummary `json:"witness,omitempty"`
	Runs                 []RunRow        `json:"runs"`
}

type WitnessSummary struct {
	Kind       string `json:"kind"`
	Status     string `json:"status"`
	TrustScope string `json:"trust_scope"`
	Reason     string `json:"reason"`
}

type ReportArtifacts struct {
	Summary          Summary
	EvidenceTable    EvidenceTable
	MissingTelemetry MissingTelemetry
	Timeline         string
}

func WriteReport(target, outDir, contractPath string) (ReportArtifacts, error) {
	if strings.TrimSpace(outDir) == "" {
		return ReportArtifacts{}, errors.New("report requires --out <dir>")
	}
	contract, err := trace.LoadContract(contractPath)
	if err != nil {
		return ReportArtifacts{}, err
	}
	rows, err := VerifiedRows(target, contract)
	if err != nil {
		return ReportArtifacts{}, err
	}
	artifacts := BuildReport(rows, contract)
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return ReportArtifacts{}, err
	}
	if err := writeJSON(filepath.Join(outDir, "summary.json"), artifacts.Summary); err != nil {
		return ReportArtifacts{}, err
	}
	if err := writeJSON(filepath.Join(outDir, "evidence-table.json"), artifacts.EvidenceTable); err != nil {
		return ReportArtifacts{}, err
	}
	if err := writeJSON(filepath.Join(outDir, "missing-telemetry.json"), artifacts.MissingTelemetry); err != nil {
		return ReportArtifacts{}, err
	}
	if err := os.WriteFile(filepath.Join(outDir, "timeline.md"), []byte(artifacts.Timeline), 0o644); err != nil {
		return ReportArtifacts{}, err
	}
	return artifacts, nil
}

func WriteGate(target, outPath, contractPath string, witnessPaths ...string) (GateResult, error) {
	if strings.TrimSpace(outPath) == "" {
		return GateResult{}, errors.New("gate requires --out <file>")
	}
	contract, err := trace.LoadContract(contractPath)
	if err != nil {
		return GateResult{}, err
	}
	rows, err := VerifiedRows(target, contract)
	if err != nil {
		return GateResult{}, err
	}
	result := EvaluateGate(rows, contract)
	if len(witnessPaths) > 0 && strings.TrimSpace(witnessPaths[0]) != "" {
		result = applyWitness(result, witnessPaths[0])
	}
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return GateResult{}, err
	}
	if err := writeJSON(outPath, result); err != nil {
		return GateResult{}, err
	}
	return result, nil
}

func VerifiedRows(target string, contract trace.Contract) ([]RunRow, error) {
	runDirs, err := DiscoverRunDirs(target)
	if err != nil {
		return nil, err
	}
	rows := make([]RunRow, 0, len(runDirs))
	for _, runDir := range runDirs {
		result, table, audit, verifyErr := verifier.VerifyRun(runDir)
		if verifyErr != nil && result.Reason == "" {
			result.Reason = verifyErr.Error()
		}
		if err := verifier.WriteVerifierArtifacts(runDir, result, table, audit); err != nil {
			result = trace.VerifierResult{
				RunID:         result.RunID,
				RunDir:        runDir,
				Result:        trace.VerdictCannotVerify,
				TrustScope:    trace.TrustScopeLocalObserved,
				Completeness:  trace.CompletenessUnknown,
				Replayability: trace.ReplayabilityNone,
				Reason:        fmt.Sprintf("failed writing verifier artifacts: %v", err),
			}
		}
		row := rowFromRun(runDir, result, contract)
		rows = append(rows, row)
	}
	return rows, nil
}

func DiscoverRunDirs(root string) ([]string, error) {
	info, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", root)
	}
	if _, err := os.Stat(filepath.Join(root, "run.json")); err == nil {
		return []string{root}, nil
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	runDirs := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(root, entry.Name())
		if _, err := os.Stat(filepath.Join(path, "run.json")); err == nil {
			runDirs = append(runDirs, path)
		}
	}
	sort.Strings(runDirs)
	if len(runDirs) == 0 {
		return nil, errors.New("no run directories found")
	}
	return runDirs, nil
}

func BuildReport(rows []RunRow, contract trace.Contract) ReportArtifacts {
	summary := Summary{
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339Nano),
		RunCount:         len(rows),
		TrustScope:       string(trace.TrustScopeLocalObserved),
		AuditGrade:       false,
		AuditGradeReason: "local observed evidence has no CI/OIDC witness or external witness checkpoint",
		Runs:             rows,
	}
	for _, row := range rows {
		switch row.Result {
		case trace.VerdictObserved:
			summary.ObservedCount++
		case trace.VerdictFail:
			summary.FailedCount++
		case trace.VerdictCannotVerify:
			summary.CannotVerifyCount++
		case trace.VerdictNotAssessed:
			summary.NotAssessedCount++
		}
	}
	return ReportArtifacts{
		Summary:       summary,
		EvidenceTable: EvidenceTable{Runs: rows},
		MissingTelemetry: MissingTelemetry{
			MissingAuditEvidence:   []string{"ci_oidc_witness", "external_witness_checkpoint"},
			MissingHarnessEvidence: missingContractEvidence(rows, contract),
			Notes: []string{
				"raw stdout and stderr are not copied into demo report artifacts",
				"contract evidence is matched from redacted event metadata only",
			},
		},
		Timeline: buildTimeline(rows),
	}
}

func EvaluateGate(rows []RunRow, contract trace.Contract) GateResult {
	result := GateResult{
		LocalGate:            GatePass,
		CIWitnessGate:        GateCannotVerify,
		AuditGradeGate:       GateCannotVerify,
		RequiredEvidence:     requiredEvidenceIDs(contract),
		GateConditions:       []string{"all_runs_observed", "all_runs_completed"},
		MissingAuditEvidence: []string{"ci_oidc_witness", "external_witness_checkpoint"},
		Runs:                 rows,
	}
	observedEvidence := map[string]bool{}
	for _, row := range rows {
		if row.Kind != "" && row.Kind != "unmatched" && row.Result == trace.VerdictObserved {
			observedEvidence[row.Kind] = true
		}
		if row.Result != trace.VerdictObserved {
			result.LocalGate = GateFail
			result.Reasons = append(result.Reasons, fmt.Sprintf("%s result is %s, expected observed", row.Name, row.Result))
		}
		if row.ClosureState != trace.ClosureStateCompleted {
			result.LocalGate = GateFail
			result.Reasons = append(result.Reasons, fmt.Sprintf("%s closure_state is %s", row.Name, row.ClosureState))
		}
	}
	for _, requirement := range contract.RequiredEvidence {
		if observedEvidence[requirement.ID] {
			result.ObservedEvidence = append(result.ObservedEvidence, requirement.ID)
			continue
		}
		result.LocalGate = GateFail
		result.Reasons = append(result.Reasons, fmt.Sprintf("missing locally observed contract evidence %s", requirement.ID))
	}
	if len(result.Reasons) == 0 {
		result.Reasons = append(result.Reasons, "local contract evidence is complete for the local gate")
	}
	result.Reasons = append(result.Reasons, "audit-grade release gate cannot verify without CI/OIDC witness and external witness checkpoint")
	return result
}

func EvaluateGateWithWitness(rows []RunRow, contract trace.Contract, witnessPath string) GateResult {
	return applyWitness(EvaluateGate(rows, contract), witnessPath)
}

func rowFromRun(runDir string, result trace.VerifierResult, contract trace.Contract) RunRow {
	row := RunRow{
		Name:          filepath.Base(runDir),
		RunID:         result.RunID,
		Result:        result.Result,
		TrustScope:    result.TrustScope,
		Completeness:  result.Completeness,
		Replayability: result.Replayability,
		Reason:        result.Reason,
		ClosureState:  trace.ClosureStateUnknown,
	}
	artifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		row.Kind = "unmatched"
		row.KindReason = "run artifact could not be loaded"
		return row
	}
	row.RunID = artifact.Manifest.RunID
	row.ClosureState = artifact.Manifest.ClosureState
	commandStarted, commandFinished := commandEvents(artifact.Events)
	row.Command = payloadString(commandStarted, "command")
	row.WrapperName = payloadString(commandStarted, "wrapper_name")
	if exitCode, ok := payloadInt(commandFinished, "exit_code"); ok {
		row.ExitCode = &exitCode
	}
	row.StdoutDigest = payloadString(commandFinished, "stdout_digest")
	row.StderrDigest = payloadString(commandFinished, "stderr_digest")
	row.Kind, row.KindReason = classify(artifact.Events, contract.RequiredEvidence)
	return row
}

func commandEvents(events []trace.Event) (trace.Event, trace.Event) {
	var started trace.Event
	var finished trace.Event
	for _, event := range events {
		switch event.EventType {
		case trace.EventCommandStarted:
			started = event
		case trace.EventCommandFinished:
			finished = event
		}
	}
	return started, finished
}

func classify(events []trace.Event, requirements []trace.EvidenceRequirement) (string, string) {
	for _, requirement := range requirements {
		if requirement.ID == "" {
			continue
		}
		for _, event := range events {
			if requirement.EventType != "" && event.EventType != trace.EventType(requirement.EventType) {
				continue
			}
			value := payloadString(event, requirement.PayloadField)
			if value == requirement.PayloadEquals {
				return requirement.ID, "matched contract evidence requirement"
			}
		}
	}
	return "unmatched", "no contract evidence requirement matched"
}

func payloadString(event trace.Event, key string) string {
	if event.EventPayload == nil {
		return ""
	}
	value, ok := event.EventPayload[key]
	if !ok || value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return typed
	default:
		return fmt.Sprint(typed)
	}
}

func payloadInt(event trace.Event, key string) (int, bool) {
	if event.EventPayload == nil {
		return 0, false
	}
	value, ok := event.EventPayload[key]
	if !ok || value == nil {
		return 0, false
	}
	switch typed := value.(type) {
	case int:
		return typed, true
	case float64:
		return int(typed), true
	case json.Number:
		i, err := typed.Int64()
		return int(i), err == nil
	default:
		return 0, false
	}
}

func requiredEvidenceIDs(contract trace.Contract) []string {
	ids := make([]string, 0, len(contract.RequiredEvidence))
	for _, requirement := range contract.RequiredEvidence {
		if requirement.ID != "" {
			ids = append(ids, requirement.ID)
		}
	}
	return ids
}

func applyWitness(result GateResult, witnessPath string) GateResult {
	record, err := loadWitnessSummary(witnessPath)
	if err != nil {
		result.CIWitnessGate = GateCannotVerify
		result.Reasons = append(result.Reasons, fmt.Sprintf("ci witness cannot verify: %v", err))
		return result
	}
	result.Witness = &record
	if record.Kind == "github-actions" && record.Status == GatePass && record.TrustScope == "ci_witnessed" {
		result.CIWitnessGate = GatePass
		result.MissingAuditEvidence = []string{"external_witness_checkpoint"}
		return result
	}
	result.CIWitnessGate = GateCannotVerify
	result.MissingAuditEvidence = []string{"ci_oidc_witness", "external_witness_checkpoint"}
	return result
}

func loadWitnessSummary(path string) (WitnessSummary, error) {
	var record WitnessSummary
	data, err := os.ReadFile(path)
	if err != nil {
		return WitnessSummary{}, err
	}
	if err := json.Unmarshal(data, &record); err != nil {
		return WitnessSummary{}, err
	}
	return record, nil
}

func missingContractEvidence(rows []RunRow, contract trace.Contract) []string {
	observed := map[string]bool{}
	for _, row := range rows {
		if row.Kind != "" && row.Kind != "unmatched" && row.Result == trace.VerdictObserved {
			observed[row.Kind] = true
		}
	}
	missing := make([]string, 0)
	for _, requirement := range contract.RequiredEvidence {
		if requirement.ID == "" || observed[requirement.ID] {
			continue
		}
		missing = append(missing, requirement.ID)
	}
	return missing
}

func buildTimeline(rows []RunRow) string {
	var builder strings.Builder
	builder.WriteString("# SDP Trace Timeline\n\n")
	builder.WriteString("| Run | Kind | Result | Trust Scope | Command | Exit |\n")
	builder.WriteString("|-----|------|--------|-------------|---------|------|\n")
	for _, row := range rows {
		exit := ""
		if row.ExitCode != nil {
			exit = fmt.Sprintf("%d", *row.ExitCode)
		}
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
			escapeMD(row.Name),
			escapeMD(row.Kind),
			escapeMD(string(row.Result)),
			escapeMD(string(row.TrustScope)),
			escapeMD(row.Command),
			escapeMD(exit),
		))
	}
	return builder.String()
}

func escapeMD(value string) string {
	return strings.ReplaceAll(value, "|", "\\|")
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
