package demo

import (
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
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
