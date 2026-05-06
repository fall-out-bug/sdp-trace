package demo

import (
	"crypto/sha256"
	"encoding/hex"
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
	GatePass             = "pass"
	GateFail             = "fail"
	GateCannotVerify     = "cannot_verify"
	GateNotAssessed      = "not_assessed"
	GateMissingTelemetry = "missing_telemetry"

	GateModeObservation     = "observation"
	GateModeAdvisoryCI      = "advisory_ci"
	GateModeProtectedFuture = "protected_future"
	GateSchemaVersion       = "block14-gate-result-v1"
)

type RunRow struct {
	Name             string                `json:"name"`
	RunID            string                `json:"run_id"`
	Kind             string                `json:"kind"`
	KindReason       string                `json:"kind_reason"`
	Command          string                `json:"command"`
	WrapperName      string                `json:"wrapper_name,omitempty"`
	ExitCode         *int                  `json:"exit_code"`
	ClosureState     string                `json:"closure_state"`
	Result           trace.VerifierVerdict `json:"result"`
	TrustScope       trace.TrustScope      `json:"trust_scope"`
	Completeness     trace.Completeness    `json:"completeness"`
	Replayability    trace.Replayability   `json:"replayability"`
	StdoutDigest     string                `json:"stdout_digest"`
	StderrDigest     string                `json:"stderr_digest"`
	Reason           string                `json:"reason,omitempty"`
	OverrideRequests []OverrideRequest     `json:"override_requests,omitempty"`
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
	SchemaVersion        string              `json:"schema_version"`
	GeneratedAt          string              `json:"generated_at"`
	LocalGate            string              `json:"local_gate"`
	CIWitnessGate        string              `json:"ci_witness_gate"`
	AuditGradeGate       string              `json:"audit_grade_gate"`
	GateMode             string              `json:"gate_mode"`
	TrustCap             string              `json:"trust_cap"`
	Reasons              []string            `json:"reasons"`
	NextActions          []string            `json:"next_actions"`
	RequiredRuns         []RequiredRunResult `json:"required_runs"`
	RequiredEvidence     []string            `json:"required_evidence"`
	ObservedEvidence     []string            `json:"observed_evidence"`
	GateConditions       []GateCondition     `json:"gate_conditions"`
	MissingAuditEvidence []string            `json:"missing_audit_evidence"`
	Witness              *WitnessSummary     `json:"witness,omitempty"`
	WitnessBindings      []WitnessBinding    `json:"witness_bindings"`
	OverrideRequests     []OverrideRequest   `json:"override_requests"`
	Runs                 []RunRow            `json:"runs"`
}

type RequiredRunResult struct {
	ID           string   `json:"id"`
	WrapperName  string   `json:"wrapper_name"`
	Profile      string   `json:"profile"`
	State        string   `json:"state"`
	MatchedRunID string   `json:"matched_run_id,omitempty"`
	Reasons      []string `json:"reasons"`
}

type GateCondition struct {
	ID     string `json:"id"`
	State  string `json:"state"`
	Reason string `json:"reason"`
}

type WitnessBinding struct {
	ID     string `json:"id"`
	State  string `json:"state"`
	Reason string `json:"reason"`
}

type OverrideRequest struct {
	OverrideID string `json:"override_id"`
	State      string `json:"state"`
	Reason     string `json:"reason,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
}

type WitnessExpectation struct {
	Repository   string
	Ref          string
	CommitSHA    string
	RunArtifacts []WitnessArtifactDigest
}

type WitnessSummary struct {
	Kind            string                  `json:"kind"`
	Status          string                  `json:"status"`
	TrustScope      string                  `json:"trust_scope"`
	Reason          string                  `json:"reason"`
	Source          WitnessSourceIdentity   `json:"source"`
	RunArtifacts    []WitnessArtifactDigest `json:"run_artifacts,omitempty"`
	ReportArtifacts []WitnessArtifactDigest `json:"report_artifacts,omitempty"`
}

type WitnessSourceIdentity struct {
	Repository string `json:"repository"`
	Ref        string `json:"ref"`
	CommitSHA  string `json:"commit_sha"`
}

type WitnessArtifactDigest struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
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
		expected, err := witnessExpectationFromTarget(target)
		if err != nil {
			result.CIWitnessGate = GateCannotVerify
			result.Reasons = append(result.Reasons, fmt.Sprintf("ci witness cannot verify current run artifacts: %v", err))
		} else {
			result = applyWitnessWithExpectation(result, witnessPaths[0], expected)
		}
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
		SchemaVersion:        GateSchemaVersion,
		GeneratedAt:          time.Now().UTC().Format(time.RFC3339),
		LocalGate:            GatePass,
		CIWitnessGate:        GateCannotVerify,
		AuditGradeGate:       GateCannotVerify,
		GateMode:             gateMode(contract),
		TrustCap:             string(trace.TrustScopeLocalObserved),
		Reasons:              []string{},
		NextActions:          []string{},
		RequiredRuns:         []RequiredRunResult{},
		RequiredEvidence:     requiredEvidenceIDs(contract),
		ObservedEvidence:     []string{},
		GateConditions:       []GateCondition{},
		MissingAuditEvidence: []string{"ci_oidc_witness", "external_witness_checkpoint"},
		WitnessBindings:      []WitnessBinding{},
		OverrideRequests:     []OverrideRequest{},
		Runs:                 rows,
	}
	observedEvidence := map[string]bool{}
	for _, row := range rows {
		result.OverrideRequests = append(result.OverrideRequests, row.OverrideRequests...)
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
	result.RequiredRuns = evaluateRequiredRuns(rows, contract)
	for _, requiredRun := range result.RequiredRuns {
		switch requiredRun.State {
		case GateMissingTelemetry:
			result.LocalGate = worseGateState(result.LocalGate, GateFail)
			for _, reason := range requiredRun.Reasons {
				result.Reasons = append(result.Reasons, reason)
			}
			result.NextActions = append(result.NextActions, fmt.Sprintf("Run required wrapper %s through sdp-trace before evaluating advisory gate.", requiredRun.WrapperName))
		case GateCannotVerify:
			result.LocalGate = worseGateState(result.LocalGate, GateCannotVerify)
			result.Reasons = append(result.Reasons, requiredRun.Reasons...)
		case GateFail:
			result.LocalGate = worseGateState(result.LocalGate, GateFail)
			result.Reasons = append(result.Reasons, requiredRun.Reasons...)
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
	result.GateConditions = gateConditions(result)
	if len(result.Reasons) == 0 {
		result.Reasons = append(result.Reasons, "local contract evidence is complete for the local gate")
	}
	result.Reasons = append(result.Reasons, "audit-grade release gate cannot verify without CI/OIDC witness and external witness checkpoint")
	sort.Strings(result.Reasons)
	sort.Strings(result.NextActions)
	return result
}

func EvaluateGateWithWitness(rows []RunRow, contract trace.Contract, witnessPath string) GateResult {
	return applyWitness(EvaluateGate(rows, contract), witnessPath)
}

func EvaluateGateWithWitnessContext(rows []RunRow, contract trace.Contract, witnessPath string, expected WitnessExpectation) GateResult {
	return applyWitnessWithExpectation(EvaluateGate(rows, contract), witnessPath, expected)
}

func PreviewWitnessBinding(witnessPath, target string) (bool, []string) {
	record, err := loadWitnessSummary(witnessPath)
	if err != nil {
		return false, []string{err.Error()}
	}
	expected, err := witnessExpectationFromTarget(target)
	if err != nil {
		return true, []string{err.Error()}
	}
	state, reasons := witnessBindingState(record, expected)
	if state == GatePass {
		return true, []string{}
	}
	return true, reasons
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
	row.OverrideRequests = overrideRequestsFromEvents(artifact.Events, contract)
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

func evaluateRequiredRuns(rows []RunRow, contract trace.Contract) []RequiredRunResult {
	results := make([]RequiredRunResult, 0, len(contract.RequiredRuns))
	for _, required := range contract.RequiredRuns {
		profile := required.Profile
		if profile == "" {
			profile = GateModeObservation
		}
		result := RequiredRunResult{
			ID:          required.ID,
			WrapperName: required.WrapperName,
			Profile:     profile,
			State:       GateMissingTelemetry,
			Reasons:     []string{fmt.Sprintf("required run %s with wrapper %s is missing", required.ID, required.WrapperName)},
		}
		for _, row := range rows {
			if row.WrapperName != required.WrapperName {
				continue
			}
			result.MatchedRunID = row.RunID
			result.State = GatePass
			result.Reasons = []string{fmt.Sprintf("required run %s matched wrapper %s", required.ID, required.WrapperName)}
			if row.Result != trace.VerdictObserved || row.ClosureState != trace.ClosureStateCompleted {
				result.State = GateCannotVerify
				result.Reasons = []string{fmt.Sprintf("required run %s cannot verify from run %s", required.ID, row.Name)}
			}
			for _, evidenceID := range required.RequiredEvidence {
				if row.Kind != evidenceID {
					result.State = GateCannotVerify
					result.Reasons = []string{fmt.Sprintf("required run %s missing evidence %s", required.ID, evidenceID)}
					break
				}
			}
			break
		}
		if profile == GateModeProtectedFuture {
			result.State = GateCannotVerify
			result.Reasons = []string{fmt.Sprintf("required run %s requests protected_future profile, which cannot verify before signed checkpoint evidence exists", required.ID)}
		}
		results = append(results, result)
	}
	return results
}

func gateMode(contract trace.Contract) string {
	mode := GateModeObservation
	for _, required := range contract.RequiredRuns {
		switch required.Profile {
		case GateModeProtectedFuture:
			return GateModeProtectedFuture
		case GateModeAdvisoryCI:
			mode = GateModeAdvisoryCI
		}
	}
	return mode
}

func gateConditions(result GateResult) []GateCondition {
	requiredRunsState := GatePass
	requiredEvidenceState := GatePass
	for _, run := range result.RequiredRuns {
		if run.State != GatePass {
			requiredRunsState = worseGateState(requiredRunsState, run.State)
		}
	}
	for _, id := range result.RequiredEvidence {
		if !containsString(result.ObservedEvidence, id) {
			requiredEvidenceState = GateFail
			break
		}
	}
	return []GateCondition{
		{ID: "all_required_runs_present", State: requiredRunsState, Reason: "required run observations are evaluated from contract declarations"},
		{ID: "all_required_evidence_observed", State: requiredEvidenceState, Reason: "contract evidence ids are matched against observed run events"},
		{ID: "ci_witness_bound_when_required", State: result.CIWitnessGate, Reason: "CI witness binding is advisory in Block 14"},
		{ID: "audit_grade_external_witness_present", State: result.AuditGradeGate, Reason: "external witness profile is not implemented in Block 14"},
	}
}

func applyWitness(result GateResult, witnessPath string) GateResult {
	return applyWitnessWithExpectation(result, witnessPath, WitnessExpectation{})
}

func applyWitnessWithExpectation(result GateResult, witnessPath string, expected WitnessExpectation) GateResult {
	record, err := loadWitnessSummary(witnessPath)
	if err != nil {
		result.CIWitnessGate = GateCannotVerify
		result.Reasons = append(result.Reasons, fmt.Sprintf("ci witness cannot verify: %v", err))
		return result
	}
	result.Witness = &record
	if record.Kind == "github-actions" && record.Status == GatePass && record.TrustScope == "ci_witnessed" {
		if bindingState, bindingReasons := witnessBindingState(record, expected); bindingState != GatePass {
			result.CIWitnessGate = bindingState
			result.Reasons = append(result.Reasons, bindingReasons...)
			for _, reason := range bindingReasons {
				result.WitnessBindings = append(result.WitnessBindings, WitnessBinding{ID: "source", State: bindingState, Reason: reason})
			}
			sort.Strings(result.Reasons)
			result.GateConditions = gateConditions(result)
			return result
		}
		result.CIWitnessGate = GatePass
		result.MissingAuditEvidence = []string{"external_witness_checkpoint"}
		result.GateConditions = gateConditions(result)
		return result
	}
	result.CIWitnessGate = GateCannotVerify
	result.MissingAuditEvidence = []string{"ci_oidc_witness", "external_witness_checkpoint"}
	result.GateConditions = gateConditions(result)
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

func witnessBindingState(record WitnessSummary, expected WitnessExpectation) (string, []string) {
	if expected.Repository != "" && record.Source.Repository == "" {
		return GateCannotVerify, []string{"ci witness repository binding is missing"}
	}
	if expected.Repository != "" && record.Source.Repository != expected.Repository {
		return GateFail, []string{fmt.Sprintf("ci witness repository mismatch: expected %s got %s", expected.Repository, record.Source.Repository)}
	}
	if expected.Ref != "" && record.Source.Ref == "" {
		return GateCannotVerify, []string{"ci witness ref binding is missing"}
	}
	if expected.Ref != "" && record.Source.Ref != expected.Ref {
		return GateFail, []string{fmt.Sprintf("ci witness ref mismatch: expected %s got %s", expected.Ref, record.Source.Ref)}
	}
	if expected.CommitSHA != "" && record.Source.CommitSHA == "" {
		return GateCannotVerify, []string{"ci witness commit binding is missing"}
	}
	if expected.CommitSHA != "" && record.Source.CommitSHA != expected.CommitSHA {
		return GateFail, []string{fmt.Sprintf("ci witness commit mismatch: expected %s got %s", expected.CommitSHA, record.Source.CommitSHA)}
	}
	expectedArtifacts := map[string]string{}
	for _, artifact := range expected.RunArtifacts {
		expectedArtifacts[artifact.Path] = artifact.SHA256
	}
	seenArtifacts := map[string]bool{}
	for _, artifact := range record.RunArtifacts {
		seenArtifacts[artifact.Path] = true
		expectedDigest, ok := expectedArtifacts[artifact.Path]
		if !ok {
			return GateCannotVerify, []string{fmt.Sprintf("ci witness artifact %s is not present in current gate input", artifact.Path)}
		}
		if expectedDigest != artifact.SHA256 {
			return GateFail, []string{fmt.Sprintf("ci witness artifact digest mismatch for %s", artifact.Path)}
		}
	}
	for path := range expectedArtifacts {
		if !seenArtifacts[path] {
			return GateCannotVerify, []string{fmt.Sprintf("ci witness artifact %s is missing from witness", path)}
		}
	}
	return GatePass, nil
}

func witnessExpectationFromTarget(target string) (WitnessExpectation, error) {
	runDirs, err := DiscoverRunDirs(target)
	if err != nil {
		return WitnessExpectation{}, err
	}
	artifacts := make([]WitnessArtifactDigest, 0, len(runDirs))
	for _, runDir := range runDirs {
		digest, err := hashFile(filepath.Join(runDir, "run.json"))
		if err != nil {
			return WitnessExpectation{}, err
		}
		artifacts = append(artifacts, WitnessArtifactDigest{
			Path:   filepath.ToSlash(filepath.Join(filepath.Base(runDir), "run.json")),
			SHA256: digest,
		})
	}
	return WitnessExpectation{RunArtifacts: artifacts}, nil
}

func hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func overrideRequestsFromEvents(events []trace.Event, contract trace.Contract) []OverrideRequest {
	requests := make([]OverrideRequest, 0)
	for _, event := range events {
		if event.EventType != trace.EventPolicyOverrideRequested {
			continue
		}
		request := OverrideRequest{
			OverrideID: payloadString(event, "override_id"),
			State:      GatePass,
			CreatedAt:  payloadString(event, "created_at"),
		}
		for _, field := range []string{"override_id", "producer", "origin", "requested_by", "reason", "source_ref", "scope", "created_at"} {
			if strings.TrimSpace(payloadString(event, field)) == "" {
				request.State = GateCannotVerify
				request.Reason = fmt.Sprintf("override request missing %s", field)
				break
			}
		}
		for _, id := range payloadStringSlice(event, "affected_required_runs") {
			if !contractHasRequiredRun(contract, id) {
				request.State = GateCannotVerify
				request.Reason = fmt.Sprintf("override request references unknown required run %s", id)
			}
		}
		for _, id := range payloadStringSlice(event, "affected_evidence") {
			if !contractHasEvidence(contract, id) {
				request.State = GateCannotVerify
				request.Reason = fmt.Sprintf("override request references unknown evidence %s", id)
			}
		}
		requests = append(requests, request)
	}
	sort.SliceStable(requests, func(i, j int) bool {
		if requests[i].CreatedAt != requests[j].CreatedAt {
			return requests[i].CreatedAt < requests[j].CreatedAt
		}
		return requests[i].OverrideID < requests[j].OverrideID
	})
	return requests
}

func payloadStringSlice(event trace.Event, key string) []string {
	value, ok := event.EventPayload[key]
	if !ok || value == nil {
		return nil
	}
	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		values := make([]string, 0, len(typed))
		for _, item := range typed {
			if text, ok := item.(string); ok {
				values = append(values, text)
			}
		}
		return values
	default:
		return nil
	}
}

func contractHasRequiredRun(contract trace.Contract, id string) bool {
	for _, required := range contract.RequiredRuns {
		if required.ID == id {
			return true
		}
	}
	return false
}

func contractHasEvidence(contract trace.Contract, id string) bool {
	for _, required := range contract.RequiredEvidence {
		if required.ID == id {
			return true
		}
	}
	return false
}

func worseGateState(current, next string) string {
	if gateSeverity(next) > gateSeverity(current) {
		return next
	}
	return current
}

func gateSeverity(state string) int {
	switch state {
	case GateFail, GateMissingTelemetry:
		return 4
	case GateCannotVerify:
		return 3
	case GateNotAssessed:
		return 2
	case GatePass:
		return 1
	default:
		return 0
	}
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
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
