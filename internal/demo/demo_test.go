package demo

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/recorder"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func TestReportAndGatePassForObservedDemoRuns(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runCommand(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "hello")
	runCommand(t, filepath.Join(root, "002-verification-run"), "verification-run", echo, "tests")
	contractPath := writeDemoContract(t, t.TempDir())

	reportDir := filepath.Join(t.TempDir(), "report")
	artifacts, err := WriteReport(root, reportDir, contractPath)
	if err != nil {
		t.Fatalf("WriteReport: %v", err)
	}
	if artifacts.Summary.RunCount != 2 {
		t.Fatalf("run count = %d", artifacts.Summary.RunCount)
	}
	if artifacts.Summary.AuditGrade {
		t.Fatalf("local report must not be audit grade")
	}
	for _, name := range []string{"summary.json", "evidence-table.json", "missing-telemetry.json", "timeline.md"} {
		if _, err := os.Stat(filepath.Join(reportDir, name)); err != nil {
			t.Fatalf("missing report artifact %s: %v", name, err)
		}
	}
	timeline, err := os.ReadFile(filepath.Join(reportDir, "timeline.md"))
	if err != nil {
		t.Fatalf("read timeline: %v", err)
	}
	if !strings.Contains(string(timeline), "agent_session_observed") || !strings.Contains(string(timeline), "verification_run_observed") {
		t.Fatalf("timeline missing classified runs:\n%s", string(timeline))
	}

	gatePath := filepath.Join(t.TempDir(), "gate-result.json")
	gate, err := WriteGate(root, gatePath, contractPath)
	if err != nil {
		t.Fatalf("WriteGate: %v", err)
	}
	if gate.LocalGate != GatePass {
		t.Fatalf("local gate = %s reasons=%v", gate.LocalGate, gate.Reasons)
	}
	if gate.AuditGradeGate != GateCannotVerify {
		t.Fatalf("audit gate = %s", gate.AuditGradeGate)
	}
	if !contains(gate.MissingAuditEvidence, "ci_oidc_witness") {
		t.Fatalf("missing audit evidence not reported: %v", gate.MissingAuditEvidence)
	}
	if contains(gate.RequiredEvidence, "all_runs_observed") {
		t.Fatalf("gate conditions leaked into required evidence: %v", gate.RequiredEvidence)
	}
	if !containsGateCondition(gate.GateConditions, "all_required_runs_present") {
		t.Fatalf("gate conditions missing: %v", gate.GateConditions)
	}
}

func TestReportAcceptsSingleRunDirectory(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	runDir := filepath.Join(t.TempDir(), "single-run")
	runCommand(t, runDir, "agent-session", echo, "hello")

	artifacts, err := WriteReport(runDir, filepath.Join(t.TempDir(), "report"), writeDemoContract(t, t.TempDir()))
	if err != nil {
		t.Fatalf("WriteReport: %v", err)
	}
	if artifacts.Summary.RunCount != 1 {
		t.Fatalf("run count = %d", artifacts.Summary.RunCount)
	}
}

func TestGateFailsForTamperedRun(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runCommand(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "hello")
	tampered := filepath.Join(root, "002-verification-run")
	runCommand(t, tampered, "verification-run", echo, "tests")

	eventPath := filepath.Join(tampered, "events", "000003-command_finished.json")
	var event map[string]any
	raw, err := os.ReadFile(eventPath)
	if err != nil {
		t.Fatalf("read event: %v", err)
	}
	if err := json.Unmarshal(raw, &event); err != nil {
		t.Fatalf("decode event: %v", err)
	}
	event["event_hash"] = "deadbeef"
	mutated, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	if err := os.WriteFile(eventPath, mutated, 0o644); err != nil {
		t.Fatalf("write event: %v", err)
	}

	gate, err := WriteGate(root, filepath.Join(t.TempDir(), "gate-result.json"), writeDemoContract(t, t.TempDir()))
	if err != nil {
		t.Fatalf("WriteGate: %v", err)
	}
	if gate.LocalGate != GateFail {
		t.Fatalf("expected fail gate, got %s", gate.LocalGate)
	}
}

func TestGateFailsForNotAssessedRun(t *testing.T) {
	contract := traceContractForTest()
	gate := EvaluateGate([]RunRow{
		{
			Name:         "agent-session",
			Kind:         "agent_session_observed",
			Result:       "observed",
			ClosureState: "completed",
		},
		{
			Name:         "verification-run",
			Kind:         "verification_run_observed",
			Result:       "observed",
			ClosureState: "completed",
		},
		{
			Name:         "not-assessed-run",
			Kind:         "unmatched",
			Result:       "not_assessed",
			ClosureState: "completed",
		},
	}, contract)
	if gate.LocalGate != GateFail {
		t.Fatalf("expected not_assessed run to fail local gate, got %s", gate.LocalGate)
	}
}

func TestGateReportsMissingRequiredRun(t *testing.T) {
	contract := traceContractForTest()
	contract.RequiredRuns = []trace.RequiredRun{
		{
			ID:               "agent_session",
			WrapperName:      "agent-session",
			RequiredEvidence: []string{"agent_session_observed"},
			Profile:          "observation",
		},
		{
			ID:               "verification_run",
			WrapperName:      "verification-run",
			RequiredEvidence: []string{"verification_run_observed"},
			Profile:          "observation",
		},
	}

	gate := EvaluateGate([]RunRow{
		{Name: "agent-session", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
	}, contract)

	if gate.LocalGate != GateFail {
		t.Fatalf("expected missing required run to fail local gate, got %s", gate.LocalGate)
	}
	missing := findRequiredRun(t, gate.RequiredRuns, "verification_run")
	if missing.State != "missing_telemetry" {
		t.Fatalf("required run state = %s reasons=%v", missing.State, missing.Reasons)
	}
	if !contains(gate.NextActions, "Run required wrapper verification-run through sdp-trace before evaluating advisory gate.") {
		t.Fatalf("missing next action for required run: %v", gate.NextActions)
	}
}

func TestEvaluateRequiredRunsUsesFirstMatchingWrapperRow(t *testing.T) {
	contract := traceContractForTest()
	contract.RequiredRuns = []trace.RequiredRun{
		{
			ID:               "agent_session",
			WrapperName:      "agent-session",
			RequiredEvidence: []string{"agent_session_observed"},
			Profile:          "observation",
		},
	}

	gate := EvaluateGate([]RunRow{
		{Name: "agent-session-bad", RunID: "run-1", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "fail", ClosureState: "completed"},
		{Name: "agent-session-good", RunID: "run-2", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
	}, contract)

	required := findRequiredRun(t, gate.RequiredRuns, "agent_session")
	if required.State != GateCannotVerify {
		t.Fatalf("expected first matching row to win with cannot_verify, got %s", required.State)
	}
	if required.MatchedRunID != "run-1" {
		t.Fatalf("expected first matching run to be selected, got %s", required.MatchedRunID)
	}
	if required.Reasons[0] != "required run agent_session cannot verify from run agent-session-bad" {
		t.Fatalf("unexpected reason = %v", required.Reasons)
	}
}

func TestEvaluateRequiredRunsRequiresAllEvidenceEntries(t *testing.T) {
	contract := traceContractForTest()
	contract.RequiredRuns = []trace.RequiredRun{
		{
			ID:               "agent_session",
			WrapperName:      "agent-session",
			RequiredEvidence: []string{"agent_session_observed", "verification_run_observed"},
		},
	}

	gate := EvaluateGate([]RunRow{
		{Name: "agent-session", RunID: "run-1", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
	}, contract)

	required := findRequiredRun(t, gate.RequiredRuns, "agent_session")
	if required.State != GateCannotVerify {
		t.Fatalf("expected missing second evidence to fail required run, got %s", required.State)
	}
	if required.Reasons[0] != "required run agent_session missing evidence verification_run_observed" {
		t.Fatalf("unexpected reason = %v", required.Reasons)
	}
}

func TestProtectedFutureRequiredRunCannotVerify(t *testing.T) {
	contract := traceContractForTest()
	contract.RequiredEvidence = nil
	contract.RequiredRuns = []trace.RequiredRun{
		{
			ID:          "protected_release",
			WrapperName: "verification-run",
			Profile:     "protected_future",
		},
	}

	gate := EvaluateGate([]RunRow{
		{Name: "verification-run", WrapperName: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}, contract)

	required := findRequiredRun(t, gate.RequiredRuns, "protected_release")
	if required.State != GateCannotVerify {
		t.Fatalf("protected future state = %s reasons=%v", required.State, required.Reasons)
	}
	if gate.LocalGate != GateCannotVerify {
		t.Fatalf("local gate = %s reasons=%v", gate.LocalGate, gate.Reasons)
	}
}

func TestLocalSignedCheckpointDoesNotUpgradeProtectedFutureGate(t *testing.T) {
	contract := traceContractForTest()
	contract.RequiredEvidence = nil
	contract.RequiredRuns = []trace.RequiredRun{
		{ID: "protected_release", WrapperName: "verification-run", Profile: "protected_future"},
	}
	checkpointResult := checkpoint.VerificationResult{
		Result:     checkpoint.StatePass,
		TrustScope: checkpoint.TrustScopeLocalSigned,
	}

	gate := EvaluateGate([]RunRow{
		{Name: "verification-run", WrapperName: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}, contract)

	if checkpointResult.Result != checkpoint.StatePass {
		t.Fatalf("test setup expected local signed checkpoint pass")
	}
	if gate.LocalGate != GateCannotVerify {
		t.Fatalf("local signed checkpoint must not upgrade protected future gate, got %s", gate.LocalGate)
	}
	if gate.AuditGradeGate != GateCannotVerify {
		t.Fatalf("local signed checkpoint must not upgrade audit gate, got %s", gate.AuditGradeGate)
	}
}

func TestProtectedGateRejectsLocalSignedCheckpointAndKeepsConditionRows(t *testing.T) {
	contract := traceContractForTest()
	rows := []RunRow{
		{Name: "agent-session", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
		{Name: "verification-run", WrapperName: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}

	gate := EvaluateProtectedGate(rows, contract, ProtectedGateInput{
		Checkpoint: checkpoint.VerificationResult{
			Result:               checkpoint.StatePass,
			TrustScope:           checkpoint.TrustScopeLocalSigned,
			SignatureState:       checkpoint.StatePass,
			RunBindingState:      checkpoint.StatePass,
			ChainBindingState:    checkpoint.StatePass,
			SourceBindingState:   checkpoint.StatePass,
			NonceBindingState:    checkpoint.StatePass,
			SignerAuthorityState: checkpoint.StatePass,
		},
		PolicyProvided: true,
		Witness: &WitnessSummary{
			Kind:        "github-actions",
			Status:      GatePass,
			TrustScope:  "ci_witnessed",
			Reason:      "ci_identity_present",
			GeneratedAt: "2026-05-06T00:00:00Z",
			Source:      WitnessSourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		},
		WitnessExpectation: WitnessExpectation{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		Now:                mustParseTime(t, "2026-05-06T00:10:00Z"),
	})

	if gate.SelectedProfile != GateProfileProtected {
		t.Fatalf("selected profile = %s", gate.SelectedProfile)
	}
	if gate.ProtectedGate != GateFail {
		t.Fatalf("protected gate = %s reasons=%v", gate.ProtectedGate, gate.Reasons)
	}
	for _, id := range protectedConditionIDs {
		if !containsProtectedCondition(gate.ProtectedConditions, id) {
			t.Fatalf("protected condition %s missing from %+v", id, gate.ProtectedConditions)
		}
	}
	condition := findProtectedCondition(t, gate.ProtectedConditions, "protected_trust_scope_satisfied")
	if condition.State != GateFail || condition.ReasonCode != "local_signed_not_protected" {
		t.Fatalf("trust scope condition = %+v", condition)
	}
}

func TestProtectedGateMapsAbsentAndStaleWitnessFreshness(t *testing.T) {
	contract := traceContractForTest()
	rows := []RunRow{
		{Name: "agent-session", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
		{Name: "verification-run", WrapperName: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}
	base := ProtectedGateInput{
		Checkpoint: checkpoint.VerificationResult{
			Result:               checkpoint.StatePass,
			TrustScope:           checkpoint.TrustScopeCISigned,
			SignatureState:       checkpoint.StatePass,
			RunBindingState:      checkpoint.StatePass,
			ChainBindingState:    checkpoint.StatePass,
			SourceBindingState:   checkpoint.StatePass,
			NonceBindingState:    checkpoint.StatePass,
			SignerAuthorityState: checkpoint.StatePass,
		},
		PolicyProvided:     true,
		WitnessExpectation: WitnessExpectation{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
		Now:                mustParseTime(t, "2026-05-06T12:00:00Z"),
	}

	absent := base
	absent.Witness = &WitnessSummary{
		Kind:       "github-actions",
		Status:     GatePass,
		TrustScope: "ci_witnessed",
		Reason:     "ci_identity_present",
		Source:     WitnessSourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
	}
	absentGate := EvaluateProtectedGate(rows, contract, absent)
	absentCondition := findProtectedCondition(t, absentGate.ProtectedConditions, "witness_freshness_valid")
	if absentGate.ProtectedGate != GateCannotVerify || absentCondition.State != GateCannotVerify || absentCondition.ReasonCode != "missing_witness_freshness" {
		t.Fatalf("absent freshness gate=%s condition=%+v", absentGate.ProtectedGate, absentCondition)
	}

	stale := base
	stale.Witness = &WitnessSummary{
		Kind:        "github-actions",
		Status:      GatePass,
		TrustScope:  "ci_witnessed",
		Reason:      "ci_identity_present",
		GeneratedAt: "2026-05-04T12:00:00Z",
		Source:      WitnessSourceIdentity{Repository: "org/repo", Ref: "refs/heads/main", CommitSHA: "abc123"},
	}
	staleGate := EvaluateProtectedGate(rows, contract, stale)
	staleCondition := findProtectedCondition(t, staleGate.ProtectedConditions, "witness_freshness_valid")
	if staleGate.ProtectedGate != GateFail || staleCondition.State != GateFail || staleCondition.ReasonCode != "stale_witness" {
		t.Fatalf("stale freshness gate=%s condition=%+v", staleGate.ProtectedGate, staleCondition)
	}
}

func TestProtectedGateMalformedOverrideDoesNotDowngradeFailure(t *testing.T) {
	contract := traceContractForTest()
	rows := []RunRow{
		{Name: "agent-session", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed", OverrideRequests: []OverrideRequest{{OverrideID: "override-1", State: GateCannotVerify, Reason: "malformed"}}},
		{Name: "verification-run", WrapperName: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}

	gate := EvaluateProtectedGate(rows, contract, ProtectedGateInput{
		Checkpoint: checkpoint.VerificationResult{
			Result:               checkpoint.StatePass,
			TrustScope:           checkpoint.TrustScopeLocalSigned,
			SignatureState:       checkpoint.StatePass,
			RunBindingState:      checkpoint.StatePass,
			ChainBindingState:    checkpoint.StatePass,
			SourceBindingState:   checkpoint.StatePass,
			NonceBindingState:    checkpoint.StatePass,
			SignerAuthorityState: checkpoint.StatePass,
		},
		PolicyProvided: true,
	})

	if gate.ProtectedGate != GateFail {
		t.Fatalf("malformed override downgraded protected failure: %+v", gate)
	}
	override := findProtectedCondition(t, gate.ProtectedConditions, "override_does_not_upgrade_profile")
	if override.State != GateCannotVerify || override.ReasonCode != "override_cannot_verify_non_upgrading" {
		t.Fatalf("override condition = %+v", override)
	}
}

func TestProtectedGateMapsMissingTelemetryToTopLevelCannotVerify(t *testing.T) {
	contract := traceContractForTest()
	contract.RequiredEvidence = nil
	contract.RequiredRuns = []trace.RequiredRun{
		{ID: "agent_session", WrapperName: "agent-session", Profile: GateModeObservation},
		{ID: "verification_run", WrapperName: "verification-run", Profile: GateModeObservation},
	}
	rows := []RunRow{
		{Name: "agent-session", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
	}
	gate := EvaluateProtectedGate(rows, contract, ProtectedGateInput{
		Checkpoint: checkpoint.VerificationResult{
			Result:               checkpoint.StatePass,
			TrustScope:           checkpoint.TrustScopeCISigned,
			SignatureState:       checkpoint.StatePass,
			RunBindingState:      checkpoint.StatePass,
			ChainBindingState:    checkpoint.StatePass,
			SourceBindingState:   checkpoint.StatePass,
			NonceBindingState:    checkpoint.StatePass,
			SignerAuthorityState: checkpoint.StatePass,
		},
		PolicyProvided: true,
		Witness: &WitnessSummary{
			Kind:        "github-actions",
			Status:      GatePass,
			TrustScope:  "ci_witnessed",
			Reason:      "ci_identity_present",
			GeneratedAt: "2026-05-06T00:00:00Z",
		},
		Now: mustParseTime(t, "2026-05-06T00:10:00Z"),
	})
	if gate.ProtectedGate != GateCannotVerify {
		t.Fatalf("protected gate = %s, want cannot_verify", gate.ProtectedGate)
	}
	condition := findProtectedCondition(t, gate.ProtectedConditions, "all_required_runs_present")
	if condition.State != GateMissingTelemetry {
		t.Fatalf("required-runs condition = %+v", condition)
	}
}

func TestProtectedGateMissingPolicyCannotVerify(t *testing.T) {
	contract := traceContractForTest()
	rows := []RunRow{
		{Name: "agent-session", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
		{Name: "verification-run", WrapperName: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}
	gate := EvaluateProtectedGate(rows, contract, ProtectedGateInput{
		Checkpoint: checkpoint.VerificationResult{
			Result:               checkpoint.StatePass,
			TrustScope:           checkpoint.TrustScopeCISigned,
			SignatureState:       checkpoint.StatePass,
			RunBindingState:      checkpoint.StatePass,
			ChainBindingState:    checkpoint.StatePass,
			SourceBindingState:   checkpoint.StatePass,
			NonceBindingState:    checkpoint.StatePass,
			SignerAuthorityState: checkpoint.StatePass,
		},
		Witness: &WitnessSummary{
			Kind:        "github-actions",
			Status:      GatePass,
			TrustScope:  "ci_witnessed",
			Reason:      "ci_identity_present",
			GeneratedAt: "2026-05-06T00:00:00Z",
		},
		Now: mustParseTime(t, "2026-05-06T00:10:00Z"),
	})
	condition := findProtectedCondition(t, gate.ProtectedConditions, "checkpoint_signer_authorized")
	if gate.ProtectedGate != GateCannotVerify || condition.State != GateCannotVerify || condition.ReasonCode != "missing_policy" {
		t.Fatalf("gate=%s signer condition=%+v", gate.ProtectedGate, condition)
	}
	trustScope := findProtectedCondition(t, gate.ProtectedConditions, "protected_trust_scope_satisfied")
	if trustScope.State != GateCannotVerify || trustScope.ReasonCode != "missing_policy" {
		t.Fatalf("trust scope condition=%+v", trustScope)
	}
}

func TestProtectedGateValidOverrideDoesNotUpgradeFailure(t *testing.T) {
	contract := traceContractForTest()
	rows := []RunRow{
		{Name: "agent-session", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed", OverrideRequests: []OverrideRequest{{OverrideID: "override-1", State: GatePass, Reason: "accepted for advisory only"}}},
		{Name: "verification-run", WrapperName: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}
	gate := EvaluateProtectedGate(rows, contract, ProtectedGateInput{
		Checkpoint: checkpoint.VerificationResult{
			Result:               checkpoint.StatePass,
			TrustScope:           checkpoint.TrustScopeLocalSigned,
			SignatureState:       checkpoint.StatePass,
			RunBindingState:      checkpoint.StatePass,
			ChainBindingState:    checkpoint.StatePass,
			SourceBindingState:   checkpoint.StatePass,
			NonceBindingState:    checkpoint.StatePass,
			SignerAuthorityState: checkpoint.StatePass,
		},
		PolicyProvided: true,
	})
	override := findProtectedCondition(t, gate.ProtectedConditions, "override_does_not_upgrade_profile")
	if gate.ProtectedGate != GateFail || override.State != GatePass || override.ReasonCode != "override_visible_non_upgrading" {
		t.Fatalf("gate=%s override=%+v", gate.ProtectedGate, override)
	}
}

func TestOverrideRequestsFromEvents(t *testing.T) {
	contract := trace.Contract{
		RequiredRuns: []trace.RequiredRun{
			{ID: "known-run"},
		},
		RequiredEvidence: []trace.EvidenceRequirement{
			{ID: "known-evidence"},
		},
	}
	t.Run("builds only override events and sorts by created_at and id", func(t *testing.T) {
		events := []trace.Event{
			{
				EventType: trace.EventCommandStarted,
				EventPayload: map[string]any{
					"override_id": "not-used",
				},
			},
			{
				EventType: trace.EventPolicyOverrideRequested,
				EventPayload: map[string]any{
					"override_id":  "b",
					"producer":     "policy",
					"origin":       "policy-service",
					"requested_by": "alice",
					"reason":       "manual",
					"source_ref":   "run-2",
					"scope":        "local",
					"created_at":   "2026-05-10T09:00:00Z",
				},
			},
			{
				EventType: trace.EventPolicyOverrideRequested,
				EventPayload: map[string]any{
					"override_id":  "a",
					"producer":     "policy",
					"origin":       "policy-service",
					"requested_by": "alice",
					"reason":       "manual",
					"source_ref":   "run-1",
					"scope":        "local",
					"created_at":   "2026-05-10T09:00:00Z",
				},
			},
		}

		got := overrideRequestsFromEvents(events, contract)
		if len(got) != 2 {
			t.Fatalf("got %d overrides", len(got))
		}
		if got[0].OverrideID != "a" || got[1].OverrideID != "b" {
			t.Fatalf("sorting order mismatch: %+v", got)
		}
		if got[0].State != GatePass || got[1].State != GatePass {
			t.Fatalf("expected pass state, got %+v", got)
		}
	})

	t.Run("marks missing required fields as cannot_verify", func(t *testing.T) {
		events := []trace.Event{
			{
				EventType: trace.EventPolicyOverrideRequested,
				EventPayload: map[string]any{
					"override_id":  "override-missing-origin",
					"producer":     "policy",
					"origin":       "",
					"requested_by": "alice",
					"reason":       "manual",
					"source_ref":   "run-1",
					"scope":        "local",
					"created_at":   "2026-05-10T09:00:00Z",
				},
			},
		}

		got := overrideRequestsFromEvents(events, contract)
		if len(got) != 1 {
			t.Fatalf("got %d overrides", len(got))
		}
		if got[0].State != GateCannotVerify || got[0].Reason != "override request missing origin" {
			t.Fatalf("unexpected override result: %+v", got[0])
		}
	})

	t.Run("uses unknown reference reason when references are not in contract", func(t *testing.T) {
		events := []trace.Event{
			{
				EventType: trace.EventPolicyOverrideRequested,
				EventPayload: map[string]any{
					"override_id":  "override-unknown-reference",
					"producer":     "policy",
					"origin":       "policy-service",
					"requested_by": "alice",
					"reason":       "manual",
					"source_ref":   "run-1",
					"scope":        "local",
					"created_at":   "2026-05-10T09:00:00Z",
					"affected_required_runs": []string{
						"known-run",
						"missing-run",
					},
					"affected_evidence": []string{
						"missing-evidence",
					},
				},
			},
		}

		got := overrideRequestsFromEvents(events, contract)
		if len(got) != 1 {
			t.Fatalf("got %d overrides", len(got))
		}
		if got[0].State != GateCannotVerify {
			t.Fatalf("expected cannot_verify state, got %s", got[0].State)
		}
		if got[0].Reason != "override request references unknown evidence missing-evidence" {
			t.Fatalf("unexpected reason: %s", got[0].Reason)
		}
	})
}

func TestProtectedGateRequiresWitnessRunIDBinding(t *testing.T) {
	state, reasons := witnessBindingState(WitnessSummary{
		Kind:        "github-actions",
		Status:      GatePass,
		TrustScope:  "ci_witnessed",
		Reason:      "ci_identity_present",
		GeneratedAt: "2026-05-06T00:00:00Z",
	}, WitnessExpectation{RunID: "trace-run-1"})
	if state != GateCannotVerify || !contains(reasons, "ci witness run id binding is missing") {
		t.Fatalf("state=%s reasons=%v", state, reasons)
	}

	state, reasons = witnessBindingState(WitnessSummary{
		Kind:        "github-actions",
		Status:      GatePass,
		TrustScope:  "ci_witnessed",
		Reason:      "ci_identity_present",
		GeneratedAt: "2026-05-06T00:00:00Z",
		CIIdentity:  WitnessCIIdentity{RunID: "other-run"},
	}, WitnessExpectation{RunID: "trace-run-1"})
	if state != GateFail || !contains(reasons, "ci witness run id mismatch: expected trace-run-1 got other-run") {
		t.Fatalf("state=%s reasons=%v", state, reasons)
	}
}

func TestWitnessBindingStateRejectsInvalidBindings(t *testing.T) {
	tests := []struct {
		name     string
		record   WitnessSummary
		expected WitnessExpectation
		state    string
		reason   string
	}{
		{
			name:     "missing repository",
			expected: WitnessExpectation{Repository: "org/repo"},
			state:    GateCannotVerify,
			reason:   "ci witness repository binding is missing",
		},
		{
			name:     "mismatched ref",
			record:   WitnessSummary{Source: WitnessSourceIdentity{Ref: "refs/heads/feature"}},
			expected: WitnessExpectation{Ref: "refs/heads/main"},
			state:    GateFail,
			reason:   "ci witness ref mismatch: expected refs/heads/main got refs/heads/feature",
		},
		{
			name:     "missing commit",
			expected: WitnessExpectation{CommitSHA: "abc123"},
			state:    GateCannotVerify,
			reason:   "ci witness commit binding is missing",
		},
		{
			name: "unknown artifact",
			record: WitnessSummary{RunArtifacts: []WitnessArtifactDigest{{
				Path:   "unexpected/run.json",
				SHA256: "digest",
			}}},
			expected: WitnessExpectation{RunArtifacts: []WitnessArtifactDigest{{
				Path:   "expected/run.json",
				SHA256: "digest",
			}}},
			state:  GateCannotVerify,
			reason: "ci witness artifact unexpected/run.json is not present in current gate input",
		},
		{
			name: "artifact digest mismatch",
			record: WitnessSummary{RunArtifacts: []WitnessArtifactDigest{{
				Path:   "run/run.json",
				SHA256: "actual",
			}}},
			expected: WitnessExpectation{RunArtifacts: []WitnessArtifactDigest{{
				Path:   "run/run.json",
				SHA256: "expected",
			}}},
			state:  GateFail,
			reason: "ci witness artifact digest mismatch for run/run.json",
		},
		{
			name: "missing artifact",
			expected: WitnessExpectation{RunArtifacts: []WitnessArtifactDigest{{
				Path:   "run/run.json",
				SHA256: "expected",
			}}},
			state:  GateCannotVerify,
			reason: "ci witness artifact run/run.json is missing from witness",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state, reasons := witnessBindingState(tt.record, tt.expected)
			if state != tt.state || !contains(reasons, tt.reason) {
				t.Fatalf("state=%s reasons=%v", state, reasons)
			}
		})
	}
}

func TestProtectedGateReasonsUseSeverityBeforeConditionOrder(t *testing.T) {
	contract := traceContractForTest()
	rows := []RunRow{
		{Name: "agent-session", WrapperName: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
		{Name: "verification-run", WrapperName: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}
	gate := EvaluateProtectedGate(rows, contract, ProtectedGateInput{
		Checkpoint: checkpoint.VerificationResult{
			Result:               checkpoint.StatePass,
			TrustScope:           checkpoint.TrustScopeLocalSigned,
			SignatureState:       checkpoint.StatePass,
			RunBindingState:      checkpoint.StatePass,
			ChainBindingState:    checkpoint.StatePass,
			SourceBindingState:   checkpoint.StatePass,
			NonceBindingState:    checkpoint.StatePass,
			SignerAuthorityState: checkpoint.StatePass,
		},
		PolicyProvided: true,
	})
	failIndex := indexOfReasonPrefix(gate.Reasons, "checkpoint_signer_not_protected:")
	cannotVerifyIndex := indexOfReasonPrefix(gate.Reasons, "missing_ci_witness:")
	if failIndex == -1 || cannotVerifyIndex == -1 {
		t.Fatalf("expected fail and cannot_verify reasons: %v", gate.Reasons)
	}
	if failIndex > cannotVerifyIndex {
		t.Fatalf("fail reason ordered after cannot_verify reason: %v", gate.Reasons)
	}
}

func TestReportAndGateArtifactsDoNotLeakSecretLikeCommand(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runCommand(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "SECRET_TOKEN_DEMO_TABLE")
	contractPath := writeDemoContract(t, t.TempDir())
	reportDir := filepath.Join(t.TempDir(), "report")
	if _, err := WriteReport(root, reportDir, contractPath); err != nil {
		t.Fatalf("WriteReport: %v", err)
	}
	evidence, err := os.ReadFile(filepath.Join(reportDir, "evidence-table.json"))
	if err != nil {
		t.Fatalf("read evidence table: %v", err)
	}
	if strings.Contains(string(evidence), "SECRET_TOKEN_DEMO_TABLE") {
		t.Fatalf("evidence table leaked secret-like command: %s", string(evidence))
	}
	gatePath := filepath.Join(t.TempDir(), "gate-result.json")
	if _, err := WriteGate(root, gatePath, contractPath); err != nil {
		t.Fatalf("WriteGate: %v", err)
	}
	gateRaw, err := os.ReadFile(gatePath)
	if err != nil {
		t.Fatalf("read gate result: %v", err)
	}
	if strings.Contains(string(gateRaw), "SECRET_TOKEN_DEMO_TABLE") {
		t.Fatalf("gate result leaked secret-like command: %s", string(gateRaw))
	}
}

func TestGateUsesPassingCIWitness(t *testing.T) {
	witnessPath := filepath.Join(t.TempDir(), "ci-witness.json")
	if err := os.WriteFile(witnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "pass",
	  "trust_scope": "ci_witnessed",
	  "reason": "ci_identity_present",
	  "generated_at": "2026-05-05T00:00:00Z",
	  "source": {"repository": "org/repo", "ref": "refs/heads/main", "commit_sha": "abc123"},
	  "ci": {"provider": "github-actions", "server_url": "https://github.com", "workflow": "sdp-trace", "job": "test", "run_id": "42", "run_attempt": "1", "actor": "octocat"},
	  "oidc": {"issuer": "https://token.actions.githubusercontent.com", "subject": "repo:org/repo:ref:refs/heads/main", "audience": "sdp-trace", "repository": "org/repo", "ref": "refs/heads/main", "sha": "abc123"},
	  "run_artifacts": [],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness: %v", err)
	}

	gate := EvaluateGateWithWitness([]RunRow{
		{Name: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
		{Name: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}, traceContractForTest(), witnessPath)
	if gate.CIWitnessGate != GatePass {
		t.Fatalf("ci witness gate = %s reasons=%v", gate.CIWitnessGate, gate.Reasons)
	}
	if contains(gate.MissingAuditEvidence, "ci_oidc_witness") {
		t.Fatalf("ci witness still listed as missing: %v", gate.MissingAuditEvidence)
	}
	if !contains(gate.MissingAuditEvidence, "external_witness_checkpoint") {
		t.Fatalf("external witness missing state lost: %v", gate.MissingAuditEvidence)
	}
	if gate.AuditGradeGate != GateCannotVerify {
		t.Fatalf("audit grade gate = %s", gate.AuditGradeGate)
	}
}

func TestGateFailsForMismatchedCIWitnessSource(t *testing.T) {
	witnessPath := filepath.Join(t.TempDir(), "ci-witness.json")
	if err := os.WriteFile(witnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "pass",
	  "trust_scope": "ci_witnessed",
	  "reason": "ci_identity_present",
	  "generated_at": "2026-05-05T00:00:00Z",
	  "source": {"repository": "org/other", "ref": "refs/heads/main", "commit_sha": "abc123"},
	  "ci": {"provider": "github-actions", "server_url": "https://github.com", "workflow": "sdp-trace", "job": "test", "run_id": "42", "run_attempt": "1", "actor": "octocat"},
	  "oidc": {"issuer": "https://token.actions.githubusercontent.com", "subject": "repo:org/other:ref:refs/heads/main", "audience": "sdp-trace", "repository": "org/other", "ref": "refs/heads/main", "sha": "abc123"},
	  "run_artifacts": [],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness: %v", err)
	}

	gate := EvaluateGateWithWitnessContext([]RunRow{
		{Name: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
		{Name: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}, traceContractForTest(), witnessPath, WitnessExpectation{
		Repository: "org/repo",
		Ref:        "refs/heads/main",
		CommitSHA:  "abc123",
	})

	if gate.CIWitnessGate != GateFail {
		t.Fatalf("ci witness gate = %s reasons=%v", gate.CIWitnessGate, gate.Reasons)
	}
	if !contains(gate.Reasons, "ci witness repository mismatch: expected org/repo got org/other") {
		t.Fatalf("missing source mismatch reason: %v", gate.Reasons)
	}
}

func TestGateUsesCannotVerifyCIWitness(t *testing.T) {
	witnessPath := filepath.Join(t.TempDir(), "ci-witness.json")
	if err := os.WriteFile(witnessPath, []byte(`{
	  "kind": "github-actions",
	  "status": "cannot_verify",
	  "trust_scope": "local_observed",
	  "reason": "missing_ci_oidc",
	  "generated_at": "2026-05-05T00:00:00Z",
	  "missing_identity_fields": ["ACTIONS_ID_TOKEN_REQUEST_URL"],
	  "source": {"repository": "", "ref": "", "commit_sha": ""},
	  "ci": {"provider": "github-actions", "server_url": "", "workflow": "", "job": "", "run_id": "", "run_attempt": "", "actor": ""},
	  "run_artifacts": [],
	  "report_artifacts": []
	}`), 0o644); err != nil {
		t.Fatalf("write witness: %v", err)
	}

	gate := EvaluateGateWithWitness([]RunRow{
		{Name: "agent-session", Kind: "agent_session_observed", Result: "observed", ClosureState: "completed"},
		{Name: "verification-run", Kind: "verification_run_observed", Result: "observed", ClosureState: "completed"},
	}, traceContractForTest(), witnessPath)
	if gate.CIWitnessGate != GateCannotVerify {
		t.Fatalf("ci witness gate = %s", gate.CIWitnessGate)
	}
	if !contains(gate.MissingAuditEvidence, "ci_oidc_witness") {
		t.Fatalf("ci witness missing state lost: %v", gate.MissingAuditEvidence)
	}
	if !contains(gate.MissingAuditEvidence, "external_witness_checkpoint") {
		t.Fatalf("external witness missing state lost: %v", gate.MissingAuditEvidence)
	}
}

func TestMalformedRunAppearsAsCannotVerify(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runCommand(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "hello")
	broken := filepath.Join(root, "002-verification-run")
	runCommand(t, broken, "verification-run", echo, "tests")
	if err := os.WriteFile(filepath.Join(broken, "run.json"), []byte("{"), 0o644); err != nil {
		t.Fatalf("break run manifest: %v", err)
	}

	gate, err := WriteGate(root, filepath.Join(t.TempDir(), "gate-result.json"), writeDemoContract(t, t.TempDir()))
	if err != nil {
		t.Fatalf("WriteGate: %v", err)
	}
	if gate.LocalGate != GateFail {
		t.Fatalf("expected fail gate, got %s", gate.LocalGate)
	}
	foundCannotVerify := false
	for _, row := range gate.Runs {
		if row.Name == "002-verification-run" && row.Result == "cannot_verify" {
			foundCannotVerify = true
		}
	}
	if !foundCannotVerify {
		t.Fatalf("malformed run was not reported as cannot_verify: %+v", gate.Runs)
	}
}

func TestVerifierArtifactWriteFailureAppearsAsCannotVerify(t *testing.T) {
	echo := mustFindCommand(t, "echo")
	root := t.TempDir()
	runCommand(t, filepath.Join(root, "001-agent-session"), "agent-session", echo, "hello")
	readOnlyRun := filepath.Join(root, "002-verification-run")
	runCommand(t, readOnlyRun, "verification-run", echo, "tests")
	verifierDir := filepath.Join(readOnlyRun, "verifier")
	if err := os.Chmod(verifierDir, 0o555); err != nil {
		t.Fatalf("chmod verifier dir: %v", err)
	}
	defer func() {
		_ = os.Chmod(verifierDir, 0o755)
	}()

	artifacts, err := WriteReport(root, filepath.Join(t.TempDir(), "report"), writeDemoContract(t, t.TempDir()))
	if err != nil {
		t.Fatalf("WriteReport should still produce report: %v", err)
	}
	if artifacts.Summary.CannotVerifyCount != 1 {
		t.Fatalf("cannot verify count = %d", artifacts.Summary.CannotVerifyCount)
	}
}

func TestReportRequiresOutputDirectory(t *testing.T) {
	if _, err := WriteReport(t.TempDir(), "", ""); err == nil {
		t.Fatalf("expected missing output directory error")
	}
}

func TestGateRequiresOutputFile(t *testing.T) {
	if _, err := WriteGate(t.TempDir(), "", ""); err == nil {
		t.Fatalf("expected missing output file error")
	}
}

func writeDemoContract(t *testing.T, dir string) string {
	t.Helper()
	path := filepath.Join(dir, "contract.json")
	contract := traceContractForTest()
	payload, err := json.Marshal(contract)
	if err != nil {
		t.Fatalf("marshal contract: %v", err)
	}
	if err := os.WriteFile(path, payload, 0o644); err != nil {
		t.Fatalf("write contract: %v", err)
	}
	return path
}

func traceContractForTest() trace.Contract {
	return trace.Contract{
		ContractID: "demo-contract-driven-test",
		Version:    "sdp-trace-event.v1",
		RequiredEvents: []string{
			string(trace.EventRecorderAttached),
			string(trace.EventRunStarted),
			string(trace.EventCommandStarted),
			string(trace.EventCommandFinished),
			string(trace.EventRunClosed),
		},
		RequiredEvidence: []trace.EvidenceRequirement{
			{
				ID:            "agent_session_observed",
				EventType:     string(trace.EventCommandStarted),
				PayloadField:  "wrapper_name",
				PayloadEquals: "agent-session",
			},
			{
				ID:            "verification_run_observed",
				EventType:     string(trace.EventCommandStarted),
				PayloadField:  "wrapper_name",
				PayloadEquals: "verification-run",
			},
		},
	}
}

func runCommand(t *testing.T, runDir, wrapperName, command string, args ...string) {
	t.Helper()
	cmd := append([]string{command}, args...)
	result, err := recorder.Run(context.Background(), recorder.RecorderOptions{
		WrapperName:        wrapperName,
		OutputDir:          runDir,
		UseDefaultContract: true,
		Command:            cmd,
	})
	if err != nil {
		t.Fatalf("recorder run: %v", err)
	}
	if result.ExitCode != 0 {
		t.Fatalf("recorder exit = %d", result.ExitCode)
	}
}

func mustFindCommand(t *testing.T, name string) string {
	t.Helper()
	path, err := exec.LookPath(name)
	if err != nil {
		t.Skipf("%s not available", name)
	}
	return path
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func containsGateCondition(values []GateCondition, target string) bool {
	for _, value := range values {
		if value.ID == target {
			return true
		}
	}
	return false
}

func containsProtectedCondition(values []ProtectedCondition, target string) bool {
	for _, value := range values {
		if value.ID == target {
			return true
		}
	}
	return false
}

func findProtectedCondition(t *testing.T, values []ProtectedCondition, target string) ProtectedCondition {
	t.Helper()
	for _, value := range values {
		if value.ID == target {
			return value
		}
	}
	t.Fatalf("protected condition %s not found in %+v", target, values)
	return ProtectedCondition{}
}

func indexOfReasonPrefix(values []string, prefix string) int {
	for i, value := range values {
		if strings.HasPrefix(value, prefix) {
			return i
		}
	}
	return -1
}

func mustParseTime(t *testing.T, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		t.Fatal(err)
	}
	return parsed
}

func findRequiredRun(t *testing.T, runs []RequiredRunResult, id string) RequiredRunResult {
	t.Helper()
	for _, run := range runs {
		if run.ID == id {
			return run
		}
	}
	t.Fatalf("required run %s not found in %+v", id, runs)
	return RequiredRunResult{}
}
