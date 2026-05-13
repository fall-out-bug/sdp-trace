package demo

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"errors"
	"fmt"
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/trace"

	"github.com/fall_out_bug/sdp-trace/internal/verifier"
	"os"
	"path/filepath"

	"sort"
	"strings"
	"time"
)

const (
	GatePass             = "pass"
	GateFail             = "fail"
	GateCannotVerify     = "cannot_verify"
	GateNotAssessed      = "not_assessed"
	GateMissingTelemetry = "missing_telemetry"

	GateModeObservation      = "observation"
	GateModeAdvisoryCI       = "advisory_ci"
	GateModeProtectedFuture  = "protected_future"
	GateProfileProtected     = "protected"
	GateSchemaVersion        = "block14-gate-result-v1"
	GateSchemaVersionBlock16 = "block16-gate-result-v1"
)

var protectedConditionIDs = []string{
	"protected_profile_explicitly_selected",
	"all_required_runs_present",
	"all_required_evidence_observed",
	"ci_witness_bound",
	"witness_freshness_valid",
	"checkpoint_signature_valid",
	"checkpoint_run_binding_valid",
	"checkpoint_signer_authorized",
	"protected_trust_scope_satisfied",
	"override_does_not_upgrade_profile",
}

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
	SchemaVersion          string                         `json:"schema_version"`
	GeneratedAt            string                         `json:"generated_at"`
	SelectedProfile        string                         `json:"selected_profile,omitempty"`
	LocalGate              string                         `json:"local_gate"`
	CIWitnessGate          string                         `json:"ci_witness_gate"`
	AuditGradeGate         string                         `json:"audit_grade_gate"`
	ProtectedGate          string                         `json:"protected_gate,omitempty"`
	GateMode               string                         `json:"gate_mode"`
	TrustCap               string                         `json:"trust_cap"`
	Reasons                []string                       `json:"reasons"`
	NextActions            []string                       `json:"next_actions"`
	CheckpointVerification *checkpoint.VerificationResult `json:"checkpoint_verification,omitempty"`
	ProtectedConditions    []ProtectedCondition           `json:"protected_conditions,omitempty"`
	RequiredRuns           []RequiredRunResult            `json:"required_runs"`
	RequiredEvidence       []string                       `json:"required_evidence"`
	ObservedEvidence       []string                       `json:"observed_evidence"`
	GateConditions         []GateCondition                `json:"gate_conditions"`
	MissingAuditEvidence   []string                       `json:"missing_audit_evidence"`
	Witness                *WitnessSummary                `json:"witness,omitempty"`
	WitnessBindings        []WitnessBinding               `json:"witness_bindings"`
	OverrideRequests       []OverrideRequest              `json:"override_requests"`
	Runs                   []RunRow                       `json:"runs"`
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

type ProtectedCondition struct {
	ID         string `json:"id"`
	State      string `json:"state"`
	ReasonCode string `json:"reason_code"`
	Reason     string `json:"reason"`
	NextAction string `json:"next_action,omitempty"`
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
	RunID        string
	RunArtifacts []WitnessArtifactDigest
}
type WitnessSummary struct {
	Kind            string                  `json:"kind"`
	Status          string                  `json:"status"`
	TrustScope      string                  `json:"trust_scope"`
	Reason          string                  `json:"reason"`
	GeneratedAt     string                  `json:"generated_at,omitempty"`
	Source          WitnessSourceIdentity   `json:"source"`
	CIIdentity      WitnessCIIdentity       `json:"ci_identity,omitempty"`
	RunArtifacts    []WitnessArtifactDigest `json:"run_artifacts,omitempty"`
	ReportArtifacts []WitnessArtifactDigest `json:"report_artifacts,omitempty"`
}

type WitnessCIIdentity struct {
	RunID string `json:"run_id,omitempty"`
}

type ProtectedGateInput struct {
	Checkpoint         checkpoint.VerificationResult
	PolicyProvided     bool
	Witness            *WitnessSummary
	WitnessExpectation WitnessExpectation
	Now                time.Time
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

var runVerdictCounters = map[trace.VerifierVerdict]func(*Summary){
	trace.VerdictObserved:     func(summary *Summary) { summary.ObservedCount++ },
	trace.VerdictFail:         func(summary *Summary) { summary.FailedCount++ },
	trace.VerdictCannotVerify: func(summary *Summary) { summary.CannotVerifyCount++ },
	trace.VerdictNotAssessed:  func(summary *Summary) { summary.NotAssessedCount++ },
}

var checkpointGateStates = map[string]string{
	checkpoint.StatePass:          GatePass,
	checkpoint.StateFail:          GateFail,
	checkpoint.StateCannotVerify:  GateCannotVerify,
	checkpoint.StateNotIntegrated: GateCannotVerify,
	checkpoint.StateNotAssessed:   GateNotAssessed,
	"":                            GateNotAssessed,
}

type witnessScalarBinding struct {
	label    string
	expected string
	actual   string
}

var gateSeverityByState = map[string]int{
	GateFail:             4,
	GateMissingTelemetry: 4,
	GateCannotVerify:     3,
	GateNotAssessed:      2,
	GatePass:             1,
}

func WriteReport(target, outDir, contractPath string) (ReportArtifacts, error) {

	if strings.TrimSpace(outDir) == "" {

		return ReportArtifacts{}, errors.New("report requires --out <dir>")
	}
	rows, contract, err := verifiedRowsForContract(target, contractPath)
	if err != nil {
		return ReportArtifacts{}, err
	}
	artifacts := BuildReport(rows, contract)
	if err := persistReportArtifacts(outDir, artifacts); err != nil {
		return ReportArtifacts{}, err
	}
	return artifacts, nil
}

func persistReportArtifacts(outDir string, artifacts ReportArtifacts) error {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	return writeReportArtifacts(outDir, artifacts)
}

func verifiedRowsForContract(target, contractPath string) ([]RunRow, trace.Contract, error) {

	contract, err := trace.LoadContract(contractPath)
	if err != nil {
		return nil, trace.Contract{}, err
	}
	rows, err := VerifiedRows(target, contract)
	if err != nil {
		return nil, trace.Contract{}, err
	}
	return rows, contract, nil
}

func writeReportArtifacts(outDir string, artifacts ReportArtifacts) error {

	writes := []struct {
		name  string
		value any
	}{
		{name: "summary.json", value: artifacts.Summary},
		{name: "evidence-table.json", value: artifacts.EvidenceTable},
		{name: "missing-telemetry.json", value: artifacts.MissingTelemetry},
	}
	for _, write := range writes {
		if err := writeJSON(filepath.Join(outDir, write.name), write.value); err != nil {

			return err
		}
	}

	return os.WriteFile(filepath.Join(outDir, "timeline.md"), []byte(artifacts.Timeline), 0o644)
}
func WriteGate(target, outPath, contractPath string, witnessPaths ...string) (GateResult, error) {

	if strings.TrimSpace(outPath) == "" {
		return GateResult{}, errors.New("gate requires --out <file>")
	}
	rows, contract, err := verifiedRowsForContract(target, contractPath)
	if err != nil {
		return GateResult{}, err
	}
	result := EvaluateGate(rows, contract)

	result = applyOptionalWitness(result, target, witnessPaths)
	if err := persistGateResult(outPath, result); err != nil {
		return GateResult{}, err
	}
	return result, nil
}

func persistGateResult(outPath string, result GateResult) error {

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	return writeJSON(outPath, result)
}

func applyOptionalWitness(result GateResult, target string, witnessPaths []string) GateResult {

	witnessPath, ok := firstWitnessPath(witnessPaths)
	if !ok {
		return result
	}
	expected, err := witnessExpectationFromTarget(target)
	if err != nil {
		result.CIWitnessGate = GateCannotVerify
		result.Reasons = append(result.Reasons, fmt.Sprintf("ci witness cannot verify current run artifacts: %v", err))
		return result
	}
	return applyWitnessWithExpectation(result, witnessPath, expected)
}

func firstWitnessPath(witnessPaths []string) (string, bool) {
	if len(witnessPaths) == 0 || strings.TrimSpace(witnessPaths[0]) == "" {

		return "", false
	}
	return witnessPaths[0], true
}

func VerifiedRows(target string, contract trace.Contract) ([]RunRow, error) {

	runDirs, err := DiscoverRunDirs(target)
	if err != nil {
		return nil, err
	}
	rows := make([]RunRow, 0, len(runDirs))
	for _, runDir := range runDirs {
		rows = append(rows, verifiedRow(runDir, contract))
	}
	return rows, nil
}

func verifiedRow(runDir string, contract trace.Contract) RunRow {

	result, table, audit, verifyErr := verifier.VerifyRun(runDir)
	if verifyErr != nil && result.Reason == "" {
		result.Reason = verifyErr.Error()
	}
	if err := verifier.WriteVerifierArtifacts(runDir, result, table, audit); err != nil {
		result = verifierArtifactWriteFailure(runDir, result.RunID, err)
	}
	return rowFromRun(runDir, result, contract)
}

func verifierArtifactWriteFailure(runDir, runID string, err error) trace.VerifierResult {

	return trace.VerifierResult{
		RunID:         runID,
		RunDir:        runDir,
		Result:        trace.VerdictCannotVerify,
		TrustScope:    trace.TrustScopeLocalObserved,
		Completeness:  trace.CompletenessUnknown,
		Replayability: trace.ReplayabilityNone,
		Reason:        fmt.Sprintf("failed writing verifier artifacts: %v", err),
	}
}

func DiscoverRunDirs(root string) ([]string, error) {
	runDirs, err := discoverRunDirsUnder(root)
	if err != nil {
		return nil, err
	}

	sort.Strings(runDirs)
	if len(runDirs) == 0 {
		return nil, errors.New("no run directories found")
	}
	return runDirs, nil
}
func hashFile(path string) (string, error) {
	data, err := os.ReadFile(path)

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), err
}

func overrideRequestsFromEvents(events []trace.Event, contract trace.Contract) []OverrideRequest {

	requests := make([]OverrideRequest, 0)
	for _, event := range events {
		if event.EventType != trace.EventPolicyOverrideRequested {
			continue
		}
		requests = append(requests, overrideRequestFromEvent(event, contract))
	}
	sortOverrideRequests(requests)
	return requests
}
func overrideRequestFromEvent(event trace.Event, contract trace.Contract) OverrideRequest {

	request := OverrideRequest{
		OverrideID: payloadString(event, "override_id"),
		State:      GatePass,
		CreatedAt:  payloadString(event, "created_at"),
	}
	request.State, request.Reason = overrideRequestFieldState(event)

	request.State, request.Reason = overrideRequestReferenceState(event, contract, request.State, request.Reason)
	return request
}

func sortOverrideRequests(requests []OverrideRequest) {

	sort.SliceStable(requests, func(i, j int) bool {
		if requests[i].CreatedAt != requests[j].CreatedAt {
			return requests[i].CreatedAt < requests[j].CreatedAt
		}
		return requests[i].OverrideID < requests[j].OverrideID
	})
}

func overrideRequestFieldState(event trace.Event) (string, string) {

	for _, field := range []string{"override_id", "producer", "origin", "requested_by", "reason", "source_ref", "scope", "created_at"} {
		if strings.TrimSpace(payloadString(event, field)) == "" {
			return GateCannotVerify, fmt.Sprintf("override request missing %s", field)
		}
	}
	return GatePass, ""
}

func overrideRequestReferenceState(event trace.Event, contract trace.Contract, state string, reason string) (string, string) {

	state, reason = overrideRequestRequiredRunState(event, contract, state, reason)
	return overrideRequestEvidenceState(event, contract, state, reason)
}

func overrideRequestRequiredRunState(event trace.Event, contract trace.Contract, state string, reason string) (string, string) {

	for _, id := range payloadStringSlice(event, "affected_required_runs") {
		if !contractHasRequiredRun(contract, id) {
			state = GateCannotVerify
			reason = fmt.Sprintf("override request references unknown required run %s", id)
		}
	}
	return state, reason
}

func overrideRequestEvidenceState(event trace.Event, contract trace.Contract, state string, reason string) (string, string) {

	for _, id := range payloadStringSlice(event, "affected_evidence") {
		if !contractHasEvidence(contract, id) {
			state = GateCannotVerify
			reason = fmt.Sprintf("override request references unknown evidence %s", id)
		}
	}
	return state, reason
}

func payloadStringSlice(event trace.Event, key string) []string {

	value := payloadValue(event, key)
	if value == nil {
		return nil
	}
	return payloadAnyStringSlice(value)
}

func payloadAnyStringSlice(value any) []string {

	switch typed := value.(type) {
	case []string:
		return typed
	case []any:
		return stringItems(typed)
	default:
		return nil
	}
}

func stringItems(items []any) []string {

	values := make([]string, 0, len(items))
	for _, item := range items {
		values = appendStringItem(values, item)
	}
	return values
}

func appendStringItem(values []string, item any) []string {
	if text, ok := item.(string); ok {

		return append(values, text)
	}
	return values
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
	return gateSeverityByState[state]
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

	observed := observedEvidenceKinds(rows)
	missing := make([]string, 0)
	for _, requirement := range contract.RequiredEvidence {
		if requirement.ID == "" || observed[requirement.ID] {
			continue
		}
		missing = append(missing, requirement.ID)
	}
	return missing
}

func observedEvidenceKinds(rows []RunRow) map[string]bool {

	observed := map[string]bool{}
	for _, row := range rows {
		if rowHasObservedEvidenceKind(row) {
			observed[row.Kind] = true
		}
	}
	return observed
}

func rowHasObservedEvidenceKind(row RunRow) bool {

	return row.Kind != "" && row.Kind != "unmatched" && row.Result == trace.VerdictObserved
}

func buildTimeline(rows []RunRow) string {
	// Timeline output is derived from sanitized row fields and escapes table
	// delimiters so command text cannot corrupt Markdown structure.
	var builder strings.Builder
	builder.WriteString("# SDP Trace Timeline\n\n")

	builder.WriteString("| Run | Kind | Result | Trust Scope | Command | Exit |\n")
	builder.WriteString("|-----|------|--------|-------------|---------|------|\n")
	for _, row := range rows {
		builder.WriteString(timelineRow(row))
	}
	return builder.String()
}

func timelineRow(row RunRow) string {

	return fmt.Sprintf("| %s | %s | %s | %s | %s | %s |\n",
		escapeMD(row.Name),
		escapeMD(row.Kind),
		escapeMD(string(row.Result)),
		escapeMD(string(row.TrustScope)),
		escapeMD(row.Command),
		escapeMD(timelineExit(row)),
	)
}

func timelineExit(row RunRow) string {
	if row.ExitCode == nil {
		return ""
	}

	return fmt.Sprintf("%d", *row.ExitCode)
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
func discoverRunDirsUnder(root string) ([]string, error) {

	if err := ensureRunRootDir(root); err != nil {
		return nil, err
	}
	if hasRunManifest(root) {
		return []string{root}, nil
	}
	return collectRunDirs(root)
}

func collectRunDirs(root string) ([]string, error) {

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	runDirs := make([]string, 0)
	for _, entry := range entries {

		path := filepath.Join(root, entry.Name())
		if isRunDirCandidate(entry, path) {
			runDirs = append(runDirs, path)
		}
	}
	return runDirs, nil
}

func ensureRunRootDir(root string) error {

	info, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", root)
	}
	return nil
}

func isRunDirCandidate(entry os.DirEntry, path string) bool {
	return entry.IsDir() && hasRunManifest(path)
}
func hasRunManifest(path string) bool {
	_, err := os.Stat(filepath.Join(path, "run.json"))
	return err == nil
}

func BuildReport(rows []RunRow, contract trace.Contract) ReportArtifacts {

	return ReportArtifacts{
		Summary:          buildSummary(rows),
		EvidenceTable:    EvidenceTable{Runs: rows},
		MissingTelemetry: buildMissingTelemetry(rows, contract),
		Timeline:         buildTimeline(rows),
	}
}

func buildSummary(rows []RunRow) Summary {

	summary := Summary{
		GeneratedAt:      time.Now().UTC().Format(time.RFC3339Nano),
		RunCount:         len(rows),
		TrustScope:       string(trace.TrustScopeLocalObserved),
		AuditGrade:       false,
		AuditGradeReason: "local observed evidence has no CI/OIDC witness or external witness checkpoint",
		Runs:             rows,
	}

	summary.applyRunVerdictCounts(rows)
	return summary
}

func buildMissingTelemetry(rows []RunRow, contract trace.Contract) MissingTelemetry {
	return MissingTelemetry{

		MissingAuditEvidence:   []string{"ci_oidc_witness", "external_witness_checkpoint"},
		MissingHarnessEvidence: missingContractEvidence(rows, contract),
		Notes: []string{

			"raw stdout and stderr are not copied into demo report artifacts",
			"contract evidence is matched from redacted event metadata only",
		},
	}
}

func (summary *Summary) applyRunVerdictCounts(rows []RunRow) {
	for _, row := range rows {

		summary.applyRunVerdictCount(row.Result)
	}
}

func (summary *Summary) applyRunVerdictCount(verdict trace.VerifierVerdict) {
	counter := runVerdictCounters[verdict]
	if counter != nil {

		counter(summary)
	}
}

func EvaluateGate(rows []RunRow, contract trace.Contract) GateResult {

	result := newGateResult(rows, contract)
	observedEvidence := applyRunRows(&result, rows)
	applyRequiredRuns(&result, rows, contract)
	applyRequiredEvidence(&result, contract, observedEvidence)
	finalizeGateResult(&result)
	return result
}

func newGateResult(rows []RunRow, contract trace.Contract) GateResult {

	return GateResult{
		SchemaVersion:  GateSchemaVersion,
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		LocalGate:      GatePass,
		CIWitnessGate:  GateCannotVerify,
		AuditGradeGate: GateCannotVerify,
		GateMode:       gateMode(contract),
		TrustCap:       string(trace.TrustScopeLocalObserved),

		Reasons:          []string{},
		NextActions:      []string{},
		RequiredRuns:     []RequiredRunResult{},
		RequiredEvidence: requiredEvidenceIDs(contract),
		ObservedEvidence: []string{},

		GateConditions:       []GateCondition{},
		MissingAuditEvidence: []string{"ci_oidc_witness", "external_witness_checkpoint"},
		WitnessBindings:      []WitnessBinding{},
		OverrideRequests:     []OverrideRequest{},
		Runs:                 rows,
	}
}

func applyRunRows(result *GateResult, rows []RunRow) map[string]bool {

	observedEvidence := map[string]bool{}
	for _, row := range rows {
		applyRunRow(result, observedEvidence, row)
	}
	return observedEvidence
}
func applyRunRow(result *GateResult, observedEvidence map[string]bool, row RunRow) {

	result.OverrideRequests = append(result.OverrideRequests, row.OverrideRequests...)
	markObservedEvidence(observedEvidence, row)
	applyRowResult(result, row)
	applyRowClosure(result, row)
}

func markObservedEvidence(observedEvidence map[string]bool, row RunRow) {
	if row.Kind != "" && row.Kind != "unmatched" && row.Result == trace.VerdictObserved {

		observedEvidence[row.Kind] = true
	}
}
func applyRowResult(result *GateResult, row RunRow) {
	if row.Result != trace.VerdictObserved {

		result.LocalGate = GateFail
		result.Reasons = append(result.Reasons, fmt.Sprintf("%s result is %s, expected observed", row.Name, row.Result))
	}
}

func applyRowClosure(result *GateResult, row RunRow) {
	if row.ClosureState != trace.ClosureStateCompleted {

		result.LocalGate = GateFail
		result.Reasons = append(result.Reasons, fmt.Sprintf("%s closure_state is %s", row.Name, row.ClosureState))
	}
}

func applyRequiredRuns(result *GateResult, rows []RunRow, contract trace.Contract) {

	result.RequiredRuns = evaluateRequiredRuns(rows, contract)
	for _, requiredRun := range result.RequiredRuns {
		applyRequiredRun(result, requiredRun)
	}
}

func applyRequiredRun(result *GateResult, requiredRun RequiredRunResult) {

	switch requiredRun.State {
	case GateMissingTelemetry:

		result.LocalGate = worseGateState(result.LocalGate, GateFail)
		result.Reasons = append(result.Reasons, requiredRun.Reasons...)
		result.NextActions = append(result.NextActions, fmt.Sprintf("Run required wrapper %s through sdp-trace before evaluating advisory gate.", requiredRun.WrapperName))
	case GateCannotVerify:

		result.LocalGate = worseGateState(result.LocalGate, GateCannotVerify)
		result.Reasons = append(result.Reasons, requiredRun.Reasons...)
	case GateFail:
		result.LocalGate = worseGateState(result.LocalGate, GateFail)
		result.Reasons = append(result.Reasons, requiredRun.Reasons...)
	}
}

func applyRequiredEvidence(result *GateResult, contract trace.Contract, observedEvidence map[string]bool) {

	for _, requirement := range contract.RequiredEvidence {
		if observedEvidence[requirement.ID] {
			result.ObservedEvidence = append(result.ObservedEvidence, requirement.ID)
			continue
		}
		result.LocalGate = GateFail
		result.Reasons = append(result.Reasons, fmt.Sprintf("missing locally observed contract evidence %s", requirement.ID))
	}
}

func finalizeGateResult(result *GateResult) {

	result.GateConditions = gateConditions(*result)
	if len(result.Reasons) == 0 {
		result.Reasons = append(result.Reasons, "local contract evidence is complete for the local gate")
	}
	result.Reasons = append(result.Reasons, "audit-grade release gate cannot verify without CI/OIDC witness and external witness checkpoint")
	sort.Strings(result.Reasons)
	sort.Strings(result.NextActions)
}

func EvaluateProtectedGate(rows []RunRow, contract trace.Contract, input ProtectedGateInput) GateResult {

	result := EvaluateGate(rows, contract)
	applyProtectedGateContext(&result, input)
	applyProtectedConditionResults(&result, input)

	return result
}

func applyProtectedGateContext(result *GateResult, input ProtectedGateInput) {

	result.SchemaVersion = GateSchemaVersionBlock16
	result.SelectedProfile = GateProfileProtected
	result.ProtectedGate = GatePass
	result.GateMode = GateProfileProtected
	result.Witness = input.Witness
	result.CIWitnessGate = protectedCIWitnessGate(input)
	result.TrustCap = protectedTrustCap(input, result.CIWitnessGate)

	result.GateConditions = gateConditions(*result)
	result.CheckpointVerification = &input.Checkpoint
}

func applyProtectedConditionResults(result *GateResult, input ProtectedGateInput) {
	result.ProtectedConditions = protectedConditions(*result, input)
	for _, condition := range result.ProtectedConditions {

		if condition.ID == "override_does_not_upgrade_profile" {
			continue
		}
		result.ProtectedGate = worseProtectedState(result.ProtectedGate, topLevelProtectedState(condition.State))
	}
	result.Reasons = append(result.Reasons, protectedReasons(result.ProtectedConditions)...)
	result.NextActions = append(result.NextActions, protectedNextActions(result.ProtectedConditions)...)
}

func protectedCIWitnessGate(input ProtectedGateInput) string {

	if input.Witness == nil {
		return GateCannotVerify
	}
	state, _ := witnessBindingState(*input.Witness, input.WitnessExpectation)
	return state
}

func protectedTrustCap(input ProtectedGateInput, ciWitnessGate string) string {

	if trustCap := protectedCheckpointTrustCap(input.Checkpoint.TrustScope); trustCap != "" {
		return trustCap
	}
	return nonCheckpointProtectedTrustCap(input.Checkpoint.TrustScope, ciWitnessGate)
}
func nonCheckpointProtectedTrustCap(checkpointTrustScope, ciWitnessGate string) string {

	if ciWitnessGate == GatePass {
		return "ci_witnessed"
	}
	if checkpointTrustScope != "" {
		return checkpointTrustScope
	}
	return string(trace.TrustScopeLocalObserved)
}

func protectedCheckpointTrustCap(trustScope string) string {

	for _, candidate := range []string{checkpoint.TrustScopeCISigned, checkpoint.TrustScopeLocalSigned} {
		if trustScope == candidate {
			return candidate
		}
	}
	return ""
}
func protectedConditions(result GateResult, input ProtectedGateInput) []ProtectedCondition {

	return []ProtectedCondition{
		protectedProfileSelectedCondition(),

		protectedConditionFromGateCondition(result.GateConditions, "all_required_runs_present"),
		protectedConditionFromGateCondition(result.GateConditions, "all_required_evidence_observed"),

		protectedCIWitnessCondition(input),
		protectedWitnessFreshnessCondition(input),

		protectedCheckpointSignatureCondition(input.Checkpoint),
		protectedCheckpointBindingCondition(input.Checkpoint),
		protectedSignerCondition(input),
		protectedTrustScopeCondition(input),
		protectedOverrideCondition(result.OverrideRequests),
	}
}

func protectedProfileSelectedCondition() ProtectedCondition {

	return ProtectedCondition{
		ID:         "protected_profile_explicitly_selected",
		State:      GatePass,
		ReasonCode: "protected_profile_selected",
		Reason:     "protected profile was explicitly selected",
	}
}

func protectedConditionFromGateCondition(conditions []GateCondition, id string) ProtectedCondition {

	for _, condition := range conditions {
		if condition.ID == id {
			return protectedConditionFromLocalGate(id, condition)
		}
	}
	return ProtectedCondition{
		ID:         id,
		State:      GateCannotVerify,
		ReasonCode: id + "_missing",
		Reason:     "required gate condition is missing",
		NextAction: "Regenerate the gate result with current sdp-trace.",
	}
}

func protectedConditionFromLocalGate(id string, condition GateCondition) ProtectedCondition {

	code := "condition_pass"
	next := ""
	if condition.State != GatePass {

		code = id + "_not_satisfied"
		next = "Supply the required run and evidence before evaluating protected profile."
	}

	return ProtectedCondition{ID: id, State: condition.State, ReasonCode: code, Reason: condition.Reason, NextAction: next}
}

func protectedCIWitnessCondition(input ProtectedGateInput) ProtectedCondition {

	if input.Witness == nil {
		return missingCIWitnessCondition()
	}
	state, reasons := witnessBindingState(*input.Witness, input.WitnessExpectation)
	code, reason, next := protectedCIWitnessFields(state, reasons)
	return ProtectedCondition{ID: "ci_witness_bound", State: state, ReasonCode: code, Reason: reason, NextAction: next}
}

func protectedCIWitnessFields(state string, reasons []string) (string, string, string) {
	if state == GatePass {

		return "ci_witness_bound", "CI witness source and artifact bindings match protected profile input", ""
	}
	return protectedCIWitnessNonPassFields(state, reasons)
}

func protectedCIWitnessNonPassFields(state string, reasons []string) (string, string, string) {

	reason := strings.Join(reasons, "; ")
	if state == GateFail {
		return "ci_witness_mismatch", reason, "Fix the CI witness source or artifact binding mismatch."
	}
	return "ci_witness_incomplete", reason, "Supply complete CI witness source and artifact bindings."
}

func missingCIWitnessCondition() ProtectedCondition {

	return ProtectedCondition{
		ID:         "ci_witness_bound",
		State:      GateCannotVerify,
		ReasonCode: "missing_ci_witness",
		Reason:     "CI witness evidence is required for protected profile",
		NextAction: "Supply a CI witness bound to the selected run.",
	}
}

func protectedWitnessFreshnessCondition(input ProtectedGateInput) ProtectedCondition {

	generatedAt, ok := protectedWitnessGeneratedAt(input.Witness)
	if !ok {
		return witnessFreshnessCannotVerify("missing_witness_freshness", "CI witness generated_at is required for protected freshness evaluation", "Supply CI witness evidence with generated_at freshness data.")
	}
	return protectedWitnessFreshnessAt(generatedAt, input.Now)
}
func protectedWitnessGeneratedAt(witness *WitnessSummary) (string, bool) {
	if witness == nil || strings.TrimSpace(witness.GeneratedAt) == "" {

		return "", false
	}
	return witness.GeneratedAt, true
}

func protectedWitnessFreshnessAt(generatedAtText string, now time.Time) ProtectedCondition {

	generatedAt, err := time.Parse(time.RFC3339, generatedAtText)
	if err != nil {

		return witnessFreshnessCannotVerify("invalid_witness_freshness", "CI witness generated_at cannot be parsed", "Regenerate CI witness evidence with an RFC3339 generated_at timestamp.")
	}
	if now.IsZero() {

		now = time.Now().UTC()
	}
	if condition, ok := invalidWitnessFreshnessCondition(generatedAt, now); ok {
		return condition
	}
	return ProtectedCondition{
		ID:         "witness_freshness_valid",
		State:      GatePass,
		ReasonCode: "witness_fresh",
		Reason:     "CI witness freshness is within the protected profile window",
	}
}

func invalidWitnessFreshnessCondition(generatedAt, now time.Time) (ProtectedCondition, bool) {

	if generatedAt.After(now.Add(5 * time.Minute)) {
		return witnessFreshnessFail("witness_from_future", "CI witness generated_at is after the verifier time window", "Regenerate CI witness evidence in the selected CI run."), true
	}
	if now.Sub(generatedAt) > 24*time.Hour {
		return witnessFreshnessFail("stale_witness", "CI witness generated_at is outside the protected freshness window", "Regenerate CI witness evidence for the selected run."), true
	}
	return ProtectedCondition{}, false
}

func witnessFreshnessCannotVerify(code, reason, next string) ProtectedCondition {
	return witnessFreshnessCondition(GateCannotVerify, code, reason, next)
}

func witnessFreshnessFail(code, reason, next string) ProtectedCondition {
	return witnessFreshnessCondition(GateFail, code, reason, next)
}
func witnessFreshnessCondition(state, code, reason, next string) ProtectedCondition {

	return ProtectedCondition{
		ID:         "witness_freshness_valid",
		State:      state,
		ReasonCode: code,
		Reason:     reason,
		NextAction: next,
	}
}

func protectedCheckpointSignatureCondition(result checkpoint.VerificationResult) ProtectedCondition {

	if result.SignatureState == checkpoint.StatePass && result.PayloadDigestState != checkpoint.StateFail {
		return ProtectedCondition{ID: "checkpoint_signature_valid", State: GatePass, ReasonCode: "checkpoint_signature_valid", Reason: "checkpoint signature verification passed"}
	}
	return ProtectedCondition{
		ID:         "checkpoint_signature_valid",
		State:      mapCheckpointState(result.SignatureState),
		ReasonCode: "checkpoint_signature_invalid",
		Reason:     "checkpoint signature verification did not pass",
		NextAction: "Regenerate the signed checkpoint for the selected run.",
	}
}

func protectedCheckpointBindingCondition(result checkpoint.VerificationResult) ProtectedCondition {

	state := GatePass
	for _, candidate := range []string{result.RunBindingState, result.ChainBindingState, result.SourceBindingState, result.NonceBindingState} {

		state = worseProtectedState(state, mapCheckpointState(candidate))
	}
	if state == GatePass {
		return ProtectedCondition{ID: "checkpoint_run_binding_valid", State: GatePass, ReasonCode: "checkpoint_binding_valid", Reason: "checkpoint binding matches the selected run context"}
	}
	return ProtectedCondition{
		ID:         "checkpoint_run_binding_valid",
		State:      state,
		ReasonCode: "checkpoint_binding_invalid",
		Reason:     "checkpoint binding does not satisfy the selected run context",
		NextAction: "Regenerate checkpoint evidence from the selected run context.",
	}
}

func protectedSignerCondition(input ProtectedGateInput) ProtectedCondition {

	if !input.PolicyProvided {
		return missingSignerPolicyCondition()
	}
	state := mapCheckpointState(input.Checkpoint.SignerAuthorityState)

	if protectedSignerPass(state, input.Checkpoint.TrustScope) {

		return ProtectedCondition{ID: "checkpoint_signer_authorized", State: GatePass, ReasonCode: "checkpoint_signer_authorized", Reason: "checkpoint signer is authorized for CI signed protected profile"}
	}
	if protectedSignerLocalOnly(state, input.Checkpoint.TrustScope) {

		state = GateFail
	}
	return ProtectedCondition{
		ID:         "checkpoint_signer_authorized",
		State:      state,
		ReasonCode: "checkpoint_signer_not_protected",
		Reason:     "checkpoint signer authority does not satisfy protected profile",
		NextAction: "Run checkpoint signing in an authorized CI signer context.",
	}
}

func missingSignerPolicyCondition() ProtectedCondition {

	return ProtectedCondition{
		ID:         "checkpoint_signer_authorized",
		State:      GateCannotVerify,
		ReasonCode: "missing_policy",
		Reason:     "trusted-checkpoint policy is required for protected profile",
		NextAction: "Supply a trusted-checkpoint policy for the protected signer.",
	}
}
func protectedSignerPass(state, trustScope string) bool {
	return state == GatePass && trustScope == checkpoint.TrustScopeCISigned
}

func protectedSignerLocalOnly(state, trustScope string) bool {
	return state == GatePass && trustScope == checkpoint.TrustScopeLocalSigned
}

func protectedTrustScopeCondition(input ProtectedGateInput) ProtectedCondition {

	if !input.PolicyProvided {
		return ProtectedCondition{
			ID:         "protected_trust_scope_satisfied",
			State:      GateCannotVerify,
			ReasonCode: "missing_policy",
			Reason:     "protected trust scope cannot be verified without trusted-checkpoint policy",
			NextAction: "Supply a trusted-checkpoint policy for the protected signer.",
		}
	}
	if protectedCheckpointCanUseWitness(input) {
		return protectedWitnessTrustScopeCondition(input)
	}
	return protectedInsufficientTrustScopeCondition(input.Checkpoint.TrustScope)
}

func protectedCheckpointCanUseWitness(input ProtectedGateInput) bool {

	return input.Checkpoint.Result == checkpoint.StatePass &&
		input.Checkpoint.TrustScope == checkpoint.TrustScopeCISigned &&
		input.Checkpoint.SignerAuthorityState == checkpoint.StatePass &&
		input.Witness != nil
}

func protectedWitnessTrustScopeCondition(input ProtectedGateInput) ProtectedCondition {

	witnessState, _ := witnessBindingState(*input.Witness, input.WitnessExpectation)
	freshness := protectedWitnessFreshnessCondition(input)
	if witnessState == GatePass && freshness.State == GatePass {

		return ProtectedCondition{ID: "protected_trust_scope_satisfied", State: GatePass, ReasonCode: "protected_trust_scope_satisfied", Reason: "CI signed checkpoint and CI witness binding satisfy protected profile"}
	}
	state := worseProtectedState(witnessState, freshness.State)
	return ProtectedCondition{
		ID:         "protected_trust_scope_satisfied",
		State:      state,
		ReasonCode: "protected_trust_scope_not_satisfied",
		Reason:     "CI signed checkpoint does not have passing CI witness binding and freshness",
		NextAction: "Provide fresh CI witness binding for the selected run.",
	}
}

func protectedInsufficientTrustScopeCondition(trustScope string) ProtectedCondition {
	code := "protected_trust_scope_not_satisfied"
	if trustScope == checkpoint.TrustScopeLocalSigned {

		code = "local_signed_not_protected"
	}
	return ProtectedCondition{
		ID:         "protected_trust_scope_satisfied",
		State:      GateFail,
		ReasonCode: code,
		Reason:     "observed trust scope does not satisfy protected profile",
		NextAction: "Provide CI signed checkpoint evidence with matching CI witness binding.",
	}
}
func protectedOverrideCondition(overrides []OverrideRequest) ProtectedCondition {

	if len(overrides) == 0 {
		return ProtectedCondition{ID: "override_does_not_upgrade_profile", State: GatePass, ReasonCode: "no_override_present", Reason: "no override request is available to upgrade the profile"}
	}
	for _, override := range overrides {

		if override.State != GatePass {

			return ProtectedCondition{
				ID:         "override_does_not_upgrade_profile",
				State:      GateCannotVerify,
				ReasonCode: "override_cannot_verify_non_upgrading",
				Reason:     "override request cannot verify and remains non-upgrading",
				NextAction: "Inspect override request evidence outside protected gate evaluation.",
			}
		}
	}
	return ProtectedCondition{ID: "override_does_not_upgrade_profile", State: GatePass, ReasonCode: "override_visible_non_upgrading", Reason: "override request is visible and non-upgrading"}
}

func mapCheckpointState(state string) string {
	if mapped, ok := checkpointGateStates[state]; ok {
		return mapped
	}

	return GateCannotVerify
}

func worseProtectedState(current, next string) string {
	if protectedSeverity(next) > protectedSeverity(current) {

		return next
	}
	return current
}

func topLevelProtectedState(state string) string {
	switch state {
	case GateMissingTelemetry, "not_integrated":

		return GateCannotVerify
	default:
		return state
	}
}

func protectedSeverity(state string) int {

	return severityByState(map[string]int{
		GateFail:             5,
		GateCannotVerify:     4,
		"not_integrated":     4,
		GateMissingTelemetry: 3,
		GateNotAssessed:      2,
		GatePass:             1,
	}, state)
}
func severityByState(values map[string]int, state string) int {

	return values[state]
}

func protectedReasons(conditions []ProtectedCondition) []string {

	ordered := orderedProtectedConditionsBySeverity(conditions)
	reasons := make([]string, 0, len(ordered))
	for _, condition := range ordered {
		if condition.ReasonCode == "" {
			continue
		}
		reasons = append(reasons, fmt.Sprintf("%s: %s", condition.ReasonCode, condition.Reason))
	}
	return reasons
}

func protectedNextActions(conditions []ProtectedCondition) []string {

	ordered := orderedProtectedConditionsBySeverity(conditions)
	actions := make([]string, 0, len(ordered))
	for _, condition := range ordered {
		if strings.TrimSpace(condition.NextAction) != "" {
			actions = append(actions, condition.NextAction)
		}
	}
	return actions
}

func orderedProtectedConditionsBySeverity(conditions []ProtectedCondition) []ProtectedCondition {

	ordered := append([]ProtectedCondition(nil), conditions...)
	positions := map[string]int{}
	for i, id := range protectedConditionIDs {

		positions[id] = i
	}
	sort.SliceStable(ordered, func(i, j int) bool {
		left := protectedSeverity(ordered[i].State)
		right := protectedSeverity(ordered[j].State)
		if left != right {

			return left > right
		}
		return positions[ordered[i].ID] < positions[ordered[j].ID]
	})
	return ordered
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

	row := rowFromVerifierResult(runDir, result)
	artifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {

		row.Kind = "unmatched"
		row.KindReason = "run artifact could not be loaded"
		return row
	}
	applyRunArtifact(&row, artifact, contract)
	return row
}
func rowFromVerifierResult(runDir string, result trace.VerifierResult) RunRow {

	return RunRow{
		Name:          filepath.Base(runDir),
		RunID:         result.RunID,
		Result:        result.Result,
		TrustScope:    result.TrustScope,
		Completeness:  result.Completeness,
		Replayability: result.Replayability,
		Reason:        result.Reason,

		ClosureState: trace.ClosureStateUnknown,
	}
}
func applyRunArtifact(row *RunRow, artifact trace.RunArtifact, contract trace.Contract) {
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
}

func commandEvents(events []trace.Event) (trace.Event, trace.Event) {
	// Last-seen command events win, matching replay behavior for append-only
	// traces that may contain retries.
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
		if evidenceRequirementMatches(events, requirement) {
			return requirement.ID, "matched contract evidence requirement"
		}
	}
	return "unmatched", "no contract evidence requirement matched"
}

func evidenceRequirementMatches(events []trace.Event, requirement trace.EvidenceRequirement) bool {

	if requirement.ID == "" {
		return false
	}
	for _, event := range events {
		if eventMatchesRequirement(event, requirement) {
			return true
		}
	}
	return false
}

func eventMatchesRequirement(event trace.Event, requirement trace.EvidenceRequirement) bool {

	if requirement.EventType != "" && event.EventType != trace.EventType(requirement.EventType) {
		return false
	}
	return payloadString(event, requirement.PayloadField) == requirement.PayloadEquals
}

func payloadString(event trace.Event, key string) string {

	value := payloadValue(event, key)
	if value == nil {
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

	value := payloadValue(event, key)
	if value == nil {
		return 0, false
	}
	return payloadAnyInt(value)
}

func payloadAnyInt(value any) (int, bool) {
	if typed, ok := value.(json.Number); ok {

		return jsonNumberInt(typed)
	}
	return primitivePayloadInt(value)
}

func primitivePayloadInt(value any) (int, bool) {

	switch typed := value.(type) {
	case int:
		return typed, true
	case float64:
		return int(typed), true
	default:
		return 0, false
	}
}

func jsonNumberInt(value json.Number) (int, bool) {
	i, err := value.Int64()
	return int(i), err == nil
}

func payloadValue(event trace.Event, key string) any {
	if event.EventPayload == nil {

		return nil
	}
	return event.EventPayload[key]
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
	rowsByWrapper := firstRowByWrapper(rows)
	for _, required := range contract.RequiredRuns {

		profile := required.Profile
		if profile == "" {
			profile = GateModeObservation
		}
		result := requiredRunResultTemplate(required, profile)
		if row, ok := rowsByWrapper[required.WrapperName]; ok {

			result = matchRequiredRun(row, required, result)
		}
		result = applyProtectedFutureConstraint(result, required.ID)
		results = append(results, result)
	}
	return results
}

func requiredRunResultTemplate(required trace.RequiredRun, profile string) RequiredRunResult {

	return RequiredRunResult{
		ID:          required.ID,
		WrapperName: required.WrapperName,
		Profile:     profile,
		State:       GateMissingTelemetry,
		Reasons: []string{
			fmt.Sprintf("required run %s with wrapper %s is missing", required.ID, required.WrapperName),
		},
	}
}

func firstRowByWrapper(rows []RunRow) map[string]RunRow {

	matches := make(map[string]RunRow, len(rows))
	for _, row := range rows {
		if _, ok := matches[row.WrapperName]; ok {
			continue
		}
		matches[row.WrapperName] = row
	}
	return matches
}

func matchRequiredRun(row RunRow, required trace.RequiredRun, result RequiredRunResult) RequiredRunResult {

	result.MatchedRunID = row.RunID
	result.State = GatePass
	result.Reasons = []string{fmt.Sprintf("required run %s matched wrapper %s", required.ID, required.WrapperName)}
	if row.Result != trace.VerdictObserved || row.ClosureState != trace.ClosureStateCompleted {
		result = cannotVerifyRequiredRun(result, required.ID, row.Name)
	}
	if evidenceID, ok := missingEvidenceID(row, required.RequiredEvidence); ok {
		return cannotVerifyRequiredRunEvidence(result, required.ID, evidenceID)
	}
	return result
}

func missingEvidenceID(row RunRow, requiredEvidence []string) (string, bool) {

	for _, evidenceID := range requiredEvidence {
		if row.Kind != evidenceID {
			return evidenceID, true
		}
	}
	return "", false
}

func cannotVerifyRequiredRunEvidence(result RequiredRunResult, requiredID, evidenceID string) RequiredRunResult {

	result.State = GateCannotVerify
	result.Reasons = []string{fmt.Sprintf("required run %s missing evidence %s", requiredID, evidenceID)}
	return result
}

func cannotVerifyRequiredRun(result RequiredRunResult, requiredID, runName string) RequiredRunResult {

	result.State = GateCannotVerify
	result.Reasons = []string{fmt.Sprintf("required run %s cannot verify from run %s", requiredID, runName)}
	return result
}

func applyProtectedFutureConstraint(result RequiredRunResult, requiredID string) RequiredRunResult {
	if result.Profile != GateModeProtectedFuture {
		return result
	}

	return cannotVerifyRequiredRunReason(result, requiredID, "requests protected_future profile, which cannot verify before signed checkpoint evidence exists")
}

func cannotVerifyRequiredRunReason(result RequiredRunResult, requiredID, reason string) RequiredRunResult {

	result.State = GateCannotVerify
	result.Reasons = []string{fmt.Sprintf("required run %s %s", requiredID, reason)}
	return result
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

	requiredRunsState := requiredRunsGateState(result.RequiredRuns)
	requiredEvidenceState := requiredEvidenceGateState(result.RequiredEvidence, result.ObservedEvidence)
	return []GateCondition{
		{ID: "all_required_runs_present", State: requiredRunsState, Reason: "required run observations are evaluated from contract declarations"},
		{ID: "all_required_evidence_observed", State: requiredEvidenceState, Reason: "contract evidence ids are matched against observed run events"},
		{ID: "ci_witness_bound_when_required", State: result.CIWitnessGate, Reason: "CI witness binding is advisory in Block 14"},
		{ID: "audit_grade_external_witness_present", State: result.AuditGradeGate, Reason: "external witness profile is not implemented in Block 14"},
	}
}

func requiredRunsGateState(requiredRuns []RequiredRunResult) string {

	state := GatePass
	for _, run := range requiredRuns {
		if run.State != GatePass {
			state = worseGateState(state, run.State)
		}
	}
	return state
}

func requiredEvidenceGateState(requiredEvidence, observedEvidence []string) string {

	for _, id := range requiredEvidence {
		if !containsString(observedEvidence, id) {
			return GateFail
		}
	}
	return GatePass
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
	if ciWitnessVerified(record) {

		return applyVerifiedCIWitness(result, record, expected)
	}
	result.CIWitnessGate = GateCannotVerify
	result.MissingAuditEvidence = []string{"ci_oidc_witness", "external_witness_checkpoint"}
	result.GateConditions = gateConditions(result)
	return result
}

func ciWitnessVerified(record WitnessSummary) bool {

	return record.Kind == "github-actions" && record.Status == GatePass && record.TrustScope == "ci_witnessed"
}

func applyVerifiedCIWitness(result GateResult, record WitnessSummary, expected WitnessExpectation) GateResult {

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

func loadWitnessSummary(path string) (WitnessSummary, error) {
	// Witness files are parsed into the portable demo summary shape before any
	// binding decision is made.
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

	for _, binding := range witnessScalarBindings(record, expected) {
		if state, reasons := validateWitnessScalarBinding(binding); state != GatePass {
			return state, reasons
		}
	}
	return witnessArtifactBindingState(record.RunArtifacts, expected.RunArtifacts)
}

func witnessScalarBindings(record WitnessSummary, expected WitnessExpectation) []witnessScalarBinding {

	return []witnessScalarBinding{
		{label: "repository", expected: expected.Repository, actual: record.Source.Repository},
		{label: "ref", expected: expected.Ref, actual: record.Source.Ref},
		{label: "commit", expected: expected.CommitSHA, actual: record.Source.CommitSHA},
		{label: "run id", expected: expected.RunID, actual: record.CIIdentity.RunID},
	}
}

func validateWitnessScalarBinding(binding witnessScalarBinding) (string, []string) {

	if binding.expected == "" {
		return GatePass, nil
	}
	if binding.actual == "" {
		return GateCannotVerify, []string{fmt.Sprintf("ci witness %s binding is missing", binding.label)}
	}
	if binding.actual != binding.expected {
		return GateFail, []string{fmt.Sprintf("ci witness %s mismatch: expected %s got %s", binding.label, binding.expected, binding.actual)}
	}
	return GatePass, nil
}
func witnessArtifactBindingState(actual, expected []WitnessArtifactDigest) (string, []string) {

	expectedArtifacts := witnessArtifactsByPath(expected)
	seenArtifacts := map[string]bool{}
	for _, artifact := range actual {
		seenArtifacts[artifact.Path] = true
		if state, reasons := validateWitnessArtifact(artifact, expectedArtifacts); state != GatePass {
			return state, reasons
		}
	}
	return missingWitnessArtifactState(expectedArtifacts, seenArtifacts)
}

func missingWitnessArtifactState(expectedArtifacts map[string]string, seenArtifacts map[string]bool) (string, []string) {

	for path := range expectedArtifacts {
		if !seenArtifacts[path] {
			return GateCannotVerify, []string{fmt.Sprintf("ci witness artifact %s is missing from witness", path)}
		}
	}
	return GatePass, nil
}

func witnessArtifactsByPath(artifacts []WitnessArtifactDigest) map[string]string {

	byPath := map[string]string{}
	for _, artifact := range artifacts {
		byPath[artifact.Path] = artifact.SHA256
	}
	return byPath
}

func validateWitnessArtifact(artifact WitnessArtifactDigest, expected map[string]string) (string, []string) {

	expectedDigest, ok := expected[artifact.Path]
	if !ok {
		return GateCannotVerify, []string{fmt.Sprintf("ci witness artifact %s is not present in current gate input", artifact.Path)}
	}
	if expectedDigest != artifact.SHA256 {
		return GateFail, []string{fmt.Sprintf("ci witness artifact digest mismatch for %s", artifact.Path)}
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

		artifact, err := witnessRunArtifactDigest(runDir)
		if err != nil {
			return WitnessExpectation{}, err
		}
		artifacts = append(artifacts, artifact)
	}
	return WitnessExpectation{RunArtifacts: artifacts}, nil
}

func witnessRunArtifactDigest(runDir string) (WitnessArtifactDigest, error) {

	digest, err := hashFile(filepath.Join(runDir, "run.json"))
	return WitnessArtifactDigest{
		Path:   filepath.ToSlash(filepath.Join(filepath.Base(runDir), "run.json")),
		SHA256: digest,
	}, err
}
