package managed

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluatePassesWithManagedBoundaryAndWitness(t *testing.T) {
	result := Evaluate(validInput())
	if result.ManagedHarnessAssessment != StatePass {
		t.Fatalf("assessment = %s reasons=%v", result.ManagedHarnessAssessment, result.Reasons)
	}
	for _, id := range managedConditionIDs {
		if conditionByID(result.ManagedConditions, id).State != StatePass {
			t.Fatalf("condition %s = %+v", id, conditionByID(result.ManagedConditions, id))
		}
	}
}

func TestPolicyRequiredEventTypesFallback(t *testing.T) {
	input := validInput()
	input.Contract.RequiredEventTypes = nil
	got := requiredEventTypes(input)
	if len(got) == 0 {
		t.Fatalf("requiredEventTypes fallback is empty")
	}
}

func TestEvaluateFailsClosedForUnmanagedAndPostHocRuns(t *testing.T) {
	tests := []struct {
		name       string
		mutate     func(*Input)
		condition  string
		reasonCode string
	}{
		{
			name: "unmanaged run",
			mutate: func(input *Input) {
				input.Run.ManagedBoundaryEnrolled = nil
			},
			condition:  "managed_boundary_enrolled_before_run",
			reasonCode: "run_not_managed",
		},
		{
			name: "late enrollment",
			mutate: func(input *Input) {
				input.Run.ManagedBoundaryEnrolled.Sequence = 5
			},
			condition:  "managed_boundary_enrolled_before_run",
			reasonCode: "late_enrollment",
		},
		{
			name: "post hoc policy",
			mutate: func(input *Input) {
				input.Policy.PolicyProvenance.Source = "run_local"
			},
			condition:  "managed_policy_loaded",
			reasonCode: "post_hoc_policy",
		},
		{
			name: "post hoc registry",
			mutate: func(input *Input) {
				input.Registry.Provenance.Source = "run_local"
			},
			condition:  "adapter_registry_loaded",
			reasonCode: "post_hoc_registry",
		},
		{
			name: "unauthorized adapter",
			mutate: func(input *Input) {
				input.Registry.Adapters[0].IdentityState = IdentitySelfClaimed
			},
			condition:  "adapter_identity_authorized",
			reasonCode: "adapter_identity_unauthorized",
		},
		{
			name: "adapter disconnect",
			mutate: func(input *Input) {
				input.Run.AdapterDisconnectObserved = true
			},
			condition:  "adapter_connection_continuous",
			reasonCode: "adapter_disconnect_observed",
		},
		{
			name: "agent reported tests",
			mutate: func(input *Input) {
				input.Run.TestEvidence = []EvidenceEvent{{EventType: "test_observed", ProvenanceScope: "agent_reported"}}
			},
			condition:  "test_provenance_not_agent_reported",
			reasonCode: "agent_reported_test_not_executed",
		},
		{
			name: "run local suppression",
			mutate: func(input *Input) {
				input.Run.SuppressedEventGroups = []SuppressedEventGroup{{EventGroup: "tool", AuthorizedByPolicy: false, PolicyProvenanceSource: "run_local"}}
			},
			condition:  "suppression_policy_valid",
			reasonCode: "suppression_unverified",
		},
		{
			name: "witness mismatch",
			mutate: func(input *Input) {
				input.Witness.RunNonce = "wrong-nonce"
			},
			condition:  "managed_witness_bound",
			reasonCode: "managed_witness_mismatch",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			input := validInput()
			tc.mutate(&input)
			result := Evaluate(input)
			if result.ManagedHarnessAssessment != StateFail {
				t.Fatalf("assessment = %s, want fail reasons=%v", result.ManagedHarnessAssessment, result.Reasons)
			}
			condition := conditionByID(result.ManagedConditions, tc.condition)
			if condition.State != StateFail || condition.ReasonCode != tc.reasonCode {
				t.Fatalf("condition = %+v, want fail/%s", condition, tc.reasonCode)
			}
		})
	}
}

func TestEvaluateCannotVerifyForMissingTestTelemetry(t *testing.T) {
	input := validInput()
	input.Run.TestEvidence = nil
	result := Evaluate(input)
	if result.ManagedHarnessAssessment != StateCannotVerify {
		t.Fatalf("assessment = %s reasons=%v", result.ManagedHarnessAssessment, result.Reasons)
	}
	condition := conditionByID(result.ManagedConditions, "test_provenance_not_agent_reported")
	if condition.State != StateMissingTelemetry || condition.ReasonCode != "test_provenance_missing" {
		t.Fatalf("test condition = %+v", condition)
	}
}

func TestEvaluateHandlesPolicyAuthorizedSuppression(t *testing.T) {
	t.Run("satisfies managed profile", func(t *testing.T) {
		input := validInput()
		input.Run.ObservedEvents = removeEvents(input.Run.ObservedEvents, "tool_call_observed")
		input.Policy.SuppressionRules = []SuppressionRule{{
			EventGroup:                   "tool",
			AuthorityRef:                 "sig:policy",
			PolicyProvenanceSource:       "vcs",
			SuppressionMaySatisfyProfile: true,
		}}
		input.Run.SuppressedEventGroups = []SuppressedEventGroup{{
			EventGroup:             "tool",
			AuthorizedByPolicy:     true,
			PolicyProvenanceSource: "vcs",
		}}
		result := Evaluate(input)
		if result.ManagedHarnessAssessment != StatePass {
			t.Fatalf("assessment = %s reasons=%v", result.ManagedHarnessAssessment, result.Reasons)
		}
		condition := conditionByID(result.ManagedConditions, "required_tool_events_observed")
		if condition.State != StatePass || condition.ReasonCode != "tool_event_suppressed_by_policy" {
			t.Fatalf("tool condition = %+v", condition)
		}
	})

	t.Run("does not satisfy managed profile", func(t *testing.T) {
		input := validInput()
		input.Run.ObservedEvents = removeEvents(input.Run.ObservedEvents, "tool_call_observed")
		input.Policy.SuppressionRules = []SuppressionRule{{
			EventGroup:                   "tool",
			AuthorityRef:                 "sig:policy",
			PolicyProvenanceSource:       "vcs",
			SuppressionMaySatisfyProfile: false,
		}}
		input.Run.SuppressedEventGroups = []SuppressedEventGroup{{
			EventGroup:             "tool",
			AuthorizedByPolicy:     true,
			PolicyProvenanceSource: "vcs",
		}}
		result := Evaluate(input)
		if result.ManagedHarnessAssessment != StateCannotVerify {
			t.Fatalf("assessment = %s reasons=%v", result.ManagedHarnessAssessment, result.Reasons)
		}
		condition := conditionByID(result.ManagedConditions, "required_tool_events_observed")
		if condition.State != StateSuppressed || condition.ReasonCode != "tool_event_suppressed" {
			t.Fatalf("tool condition = %+v", condition)
		}
	})
}

func TestEvaluateCannotVerifyForMissingRequiredCapabilityAndWitness(t *testing.T) {
	input := validInput()
	input.Registry.Adapters[0].Capabilities = []Capability{{ID: "harness-events", EventTypes: []string{"harness_lifecycle_observed"}, ProvenanceScope: "harness_observed"}}
	input.Witness = Witness{}
	result := Evaluate(input)
	if result.ManagedHarnessAssessment != StateCannotVerify {
		t.Fatalf("assessment = %s reasons=%v", result.ManagedHarnessAssessment, result.Reasons)
	}
	if conditionByID(result.ManagedConditions, "adapter_capabilities_satisfy_contract").State != StateCannotVerify {
		t.Fatalf("capability condition = %+v", conditionByID(result.ManagedConditions, "adapter_capabilities_satisfy_contract"))
	}
	if conditionByID(result.ManagedConditions, "managed_witness_bound").State != StateCannotVerify {
		t.Fatalf("witness condition = %+v", conditionByID(result.ManagedConditions, "managed_witness_bound"))
	}
}

func TestEvaluateCannotVerifyForMismatchedCapabilityReferences(t *testing.T) {
	input := validInput()
	input.Registry.Adapters[0].CapabilityRefs = []string{"harness-events", "tool-events", "file-events"}
	result := Evaluate(input)
	if result.ManagedHarnessAssessment != StateCannotVerify {
		t.Fatalf("assessment = %s reasons=%v", result.ManagedHarnessAssessment, result.Reasons)
	}
	condition := conditionByID(result.ManagedConditions, "adapter_capabilities_satisfy_contract")
	if condition.State != StateCannotVerify || condition.ReasonCode != "adapter_capability_missing" {
		t.Fatalf("capability condition = %+v", condition)
	}
}

func TestCapabilityConditionPolicyAndEventCoverage(t *testing.T) {
	for name, tc := range map[string]struct {
		mutate     func(*Input)
		wantState  string
		wantReason string
	}{
		"missing-authorized-capability-id": {
			mutate: func(input *Input) {
				input.Policy.AuthorizedAdapters[0].CapabilityIDs = []string{"harness-events", "missing-capability"}
			},
			wantState:  StateCannotVerify,
			wantReason: "adapter_capability_missing",
		},
		"missing-required-event-type": {
			mutate: func(input *Input) {
				input.Registry.Adapters[0].Capabilities[3].EventTypes = nil
			},
			wantState:  StateCannotVerify,
			wantReason: "adapter_capability_missing",
		},
		"valid-capabilities": {
			mutate:     func(*Input) {},
			wantState:  StatePass,
			wantReason: "adapter_capabilities_satisfy_contract",
		},
	} {
		t.Run(name, func(t *testing.T) {
			input := validInput()
			tc.mutate(&input)
			condition := capabilityCondition(input)
			if condition.State != tc.wantState || condition.ReasonCode != tc.wantReason {
				t.Fatalf("capabilityCondition() = %+v, want state=%s reason=%s", condition, tc.wantState, tc.wantReason)
			}
			if condition.State != StatePass && condition.NextAction == "" {
				t.Fatalf("capabilityCondition() = %+v, want non-pass next action", condition)
			}
		})
	}
}

func TestEvaluateCannotVerifyWhenWitnessHasNoArtifactBinding(t *testing.T) {
	input := validInput()
	input.Run.OutputArtifacts = nil
	input.Witness.ArtifactDigests = nil
	result := Evaluate(input)
	if result.ManagedHarnessAssessment != StateCannotVerify {
		t.Fatalf("assessment = %s reasons=%v", result.ManagedHarnessAssessment, result.Reasons)
	}
	condition := conditionByID(result.ManagedConditions, "managed_witness_bound")
	if condition.State != StateCannotVerify || condition.ReasonCode != "managed_witness_missing" {
		t.Fatalf("witness condition = %+v", condition)
	}
}

func TestWitnessConditionRejectsInvalidBindings(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Input)
		state  string
		reason string
	}{
		{name: "missing witness id", mutate: func(input *Input) { input.Witness.WitnessID = "" }, state: StateCannotVerify, reason: "managed_witness_missing"},
		{name: "stale witness", mutate: func(input *Input) { input.Witness.FreshnessState = StateCannotVerify }, state: StateCannotVerify, reason: "managed_witness_missing"},
		{name: "missing run artifacts", mutate: func(input *Input) { input.Run.OutputArtifacts = nil }, state: StateCannotVerify, reason: "managed_witness_missing"},
		{name: "missing witness artifacts", mutate: func(input *Input) { input.Witness.ArtifactDigests = nil }, state: StateCannotVerify, reason: "managed_witness_missing"},
		{name: "missing boundary", mutate: func(input *Input) { input.Run.ManagedBoundaryEnrolled = nil }, state: StateFail, reason: "managed_witness_mismatch"},
		{name: "run mismatch", mutate: func(input *Input) { input.Witness.RunNonce = "wrong" }, state: StateFail, reason: "managed_witness_mismatch"},
		{name: "authority mismatch", mutate: func(input *Input) { input.Witness.ManagedPolicyDigest = "wrong" }, state: StateFail, reason: "managed_witness_mismatch"},
		{name: "event mismatch", mutate: func(input *Input) { input.Witness.LaunchEventDigest = "wrong" }, state: StateFail, reason: "managed_witness_mismatch"},
		{name: "artifact mismatch", mutate: func(input *Input) { input.Witness.ArtifactDigests[0].SHA256 = "wrong" }, state: StateFail, reason: "managed_witness_mismatch"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := validInput()
			tt.mutate(&input)
			condition := witnessCondition(input)
			if condition.State != tt.state || condition.ReasonCode != tt.reason {
				t.Fatalf("condition = %+v, want %s/%s", condition, tt.state, tt.reason)
			}
		})
	}
}

func TestBlock17CommittedFixturesHaveManagedAssessmentShape(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "examples", "block17-managed-harness")
	cases := block17FixtureCases()
	if os.Getenv("SDP_TRACE_UPDATE_BLOCK17_FIXTURES") == "1" {
		for _, fixture := range cases {
			result := fixture.result()
			payload, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				t.Fatalf("marshal %s: %v", fixture.name, err)
			}
			payload = append(payload, '\n')
			if err := os.WriteFile(filepath.Join(fixtureDir, fixture.name), payload, 0o644); err != nil {
				t.Fatalf("write %s: %v", fixture.name, err)
			}
		}
	}

	expected := map[string]block17FixtureCase{}
	for _, fixture := range cases {
		expected[fixture.name] = fixture
	}
	entries, err := os.ReadDir(fixtureDir)
	if err != nil {
		t.Fatalf("read fixture dir: %v", err)
	}
	seen := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		fixture, ok := expected[entry.Name()]
		if !ok {
			t.Fatalf("unexpected Block 17 fixture %s", entry.Name())
		}
		payload, err := os.ReadFile(filepath.Join(fixtureDir, entry.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", entry.Name(), err)
		}
		var result AssessmentResult
		if err := json.Unmarshal(payload, &result); err != nil {
			t.Fatalf("unmarshal %s: %v", entry.Name(), err)
		}
		if result.SchemaVersion != SchemaVersion || result.SelectedProfile != ProfileManagedHarness {
			t.Fatalf("%s schema/profile = %s/%s", entry.Name(), result.SchemaVersion, result.SelectedProfile)
		}
		if len(result.ManagedConditions) != len(managedConditionIDs) {
			t.Fatalf("%s condition count = %d", entry.Name(), len(result.ManagedConditions))
		}
		for i, id := range managedConditionIDs {
			if result.ManagedConditions[i].ID != id {
				t.Fatalf("%s condition %d = %s, want %s", entry.Name(), i, result.ManagedConditions[i].ID, id)
			}
		}
		if result.ManagedHarnessAssessment != topLevel(result.ManagedConditions) {
			t.Fatalf("%s assessment = %s, want %s from conditions", entry.Name(), result.ManagedHarnessAssessment, topLevel(result.ManagedConditions))
		}
		if result.ManagedHarnessAssessment != fixture.topLevel {
			t.Fatalf("%s assessment = %s, want %s", entry.Name(), result.ManagedHarnessAssessment, fixture.topLevel)
		}
		condition := conditionByID(result.ManagedConditions, fixture.conditionID)
		if condition.State != fixture.conditionState || condition.ReasonCode != fixture.reasonCode {
			t.Fatalf("%s condition %s = %s/%s, want %s/%s", entry.Name(), fixture.conditionID, condition.State, condition.ReasonCode, fixture.conditionState, fixture.reasonCode)
		}
		seen++
	}
	if seen != len(expected) {
		t.Fatalf("seen fixtures = %d, want %d", seen, len(expected))
	}
}

type block17FixtureCase struct {
	name           string
	mutate         func(*Input)
	topLevel       string
	conditionID    string
	conditionState string
	reasonCode     string
}

func (fixture block17FixtureCase) result() AssessmentResult {
	input := validInput()
	if fixture.mutate != nil {
		fixture.mutate(&input)
	}
	return Evaluate(input)
}

func block17FixtureCases() []block17FixtureCase {
	return []block17FixtureCase{
		{
			name:           "valid-managed-profile.assessment-result.json",
			topLevel:       StatePass,
			conditionID:    "managed_witness_bound",
			conditionState: StatePass,
			reasonCode:     "managed_witness_bound",
		},
		{
			name: "unmanaged-run-fail.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.ManagedBoundaryEnrolled = nil
			},
			topLevel:       StateFail,
			conditionID:    "managed_boundary_enrolled_before_run",
			conditionState: StateFail,
			reasonCode:     "run_not_managed",
		},
		{
			name: "late-enrollment-fail.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.ManagedBoundaryEnrolled.Sequence = 5
			},
			topLevel:       StateFail,
			conditionID:    "managed_boundary_enrolled_before_run",
			conditionState: StateFail,
			reasonCode:     "late_enrollment",
		},
		{
			name: "post-hoc-policy-fail.assessment-result.json",
			mutate: func(input *Input) {
				input.Policy.PolicyProvenance.Source = "run_local"
			},
			topLevel:       StateFail,
			conditionID:    "managed_policy_loaded",
			conditionState: StateFail,
			reasonCode:     "post_hoc_policy",
		},
		{
			name: "post-hoc-registry-fail.assessment-result.json",
			mutate: func(input *Input) {
				input.Registry.Provenance.Source = "run_local"
			},
			topLevel:       StateFail,
			conditionID:    "adapter_registry_loaded",
			conditionState: StateFail,
			reasonCode:     "post_hoc_registry",
		},
		{
			name: "unauthorized-adapter-fail.assessment-result.json",
			mutate: func(input *Input) {
				input.Registry.Adapters[0].IdentityState = IdentitySelfClaimed
			},
			topLevel:       StateFail,
			conditionID:    "adapter_identity_authorized",
			conditionState: StateFail,
			reasonCode:     "adapter_identity_unauthorized",
		},
		{
			name: "adapter-disconnect-fail.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.AdapterDisconnectObserved = true
			},
			topLevel:       StateFail,
			conditionID:    "adapter_connection_continuous",
			conditionState: StateFail,
			reasonCode:     "adapter_disconnect_observed",
		},
		{
			name: "missing-capability-cannot-verify.assessment-result.json",
			mutate: func(input *Input) {
				input.Registry.Adapters[0].Capabilities = []Capability{{ID: "harness-events", EventTypes: []string{"harness_lifecycle_observed"}, ProvenanceScope: "harness_observed"}}
			},
			topLevel:       StateCannotVerify,
			conditionID:    "adapter_capabilities_satisfy_contract",
			conditionState: StateCannotVerify,
			reasonCode:     "adapter_capability_missing",
		},
		{
			name: "missing-harness-event-cannot-verify.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.ObservedEvents = removeEvents(input.Run.ObservedEvents, "harness_lifecycle_observed")
			},
			topLevel:       StateCannotVerify,
			conditionID:    "required_harness_events_observed",
			conditionState: StateMissingTelemetry,
			reasonCode:     "harness_event_missing",
		},
		{
			name: "missing-tool-event-cannot-verify.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.ObservedEvents = removeEvents(input.Run.ObservedEvents, "tool_call_observed")
			},
			topLevel:       StateCannotVerify,
			conditionID:    "required_tool_events_observed",
			conditionState: StateMissingTelemetry,
			reasonCode:     "tool_event_missing",
		},
		{
			name: "missing-file-mutation-event-cannot-verify.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.ObservedEvents = removeEvents(input.Run.ObservedEvents, "file_mutation_observed")
			},
			topLevel:       StateCannotVerify,
			conditionID:    "required_file_mutations_observed",
			conditionState: StateMissingTelemetry,
			reasonCode:     "file_mutation_event_missing",
		},
		{
			name: "missing-test-telemetry-cannot-verify.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.TestEvidence = nil
			},
			topLevel:       StateCannotVerify,
			conditionID:    "test_provenance_not_agent_reported",
			conditionState: StateMissingTelemetry,
			reasonCode:     "test_provenance_missing",
		},
		{
			name: "agent-reported-test-evidence-fail.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.TestEvidence = []EvidenceEvent{{EventType: "test_observed", ProvenanceScope: "agent_reported"}}
			},
			topLevel:       StateFail,
			conditionID:    "test_provenance_not_agent_reported",
			conditionState: StateFail,
			reasonCode:     "agent_reported_test_not_executed",
		},
		{
			name: "policy-authorized-suppression-pass.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.ObservedEvents = removeEvents(input.Run.ObservedEvents, "tool_call_observed")
				input.Policy.SuppressionRules = []SuppressionRule{{
					EventGroup:                   "tool",
					AuthorityRef:                 "sig:policy",
					PolicyProvenanceSource:       "vcs",
					SuppressionMaySatisfyProfile: true,
				}}
				input.Run.SuppressedEventGroups = []SuppressedEventGroup{{EventGroup: "tool", AuthorizedByPolicy: true, PolicyProvenanceSource: "vcs"}}
			},
			topLevel:       StatePass,
			conditionID:    "required_tool_events_observed",
			conditionState: StatePass,
			reasonCode:     "tool_event_suppressed_by_policy",
		},
		{
			name: "suppression-without-policy-fail.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.ObservedEvents = removeEvents(input.Run.ObservedEvents, "tool_call_observed")
				input.Run.SuppressedEventGroups = []SuppressedEventGroup{{EventGroup: "tool", AuthorizedByPolicy: false, PolicyProvenanceSource: "run_local"}}
			},
			topLevel:       StateFail,
			conditionID:    "suppression_policy_valid",
			conditionState: StateFail,
			reasonCode:     "suppression_unverified",
		},
		{
			name: "witness-missing-cannot-verify.assessment-result.json",
			mutate: func(input *Input) {
				input.Witness = Witness{}
			},
			topLevel:       StateCannotVerify,
			conditionID:    "managed_witness_bound",
			conditionState: StateCannotVerify,
			reasonCode:     "managed_witness_missing",
		},
		{
			name: "stale-witness-cannot-verify.assessment-result.json",
			mutate: func(input *Input) {
				input.Witness.FreshnessState = StateCannotVerify
			},
			topLevel:       StateCannotVerify,
			conditionID:    "managed_witness_bound",
			conditionState: StateCannotVerify,
			reasonCode:     "managed_witness_missing",
		},
		{
			name: "witness-mismatch-fail.assessment-result.json",
			mutate: func(input *Input) {
				input.Witness.RunNonce = "wrong-nonce"
			},
			topLevel:       StateFail,
			conditionID:    "managed_witness_bound",
			conditionState: StateFail,
			reasonCode:     "managed_witness_mismatch",
		},
		{
			name: "override-present-non-upgrading-pass.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.OverridePresent = true
			},
			topLevel:       StatePass,
			conditionID:    "override_does_not_upgrade_managed_profile",
			conditionState: StatePass,
			reasonCode:     "override_present_non_upgrading",
		},
		{
			name: "override-upgrade-fail.assessment-result.json",
			mutate: func(input *Input) {
				input.Run.OverrideAttemptsTrustUpgrade = true
			},
			topLevel:       StateFail,
			conditionID:    "override_does_not_upgrade_managed_profile",
			conditionState: StateFail,
			reasonCode:     "override_upgrade_rejected",
		},
	}
}

func conditionByID(conditions []Condition, id string) Condition {
	for _, condition := range conditions {
		if condition.ID == id {
			return condition
		}
	}
	return Condition{}
}

func removeEvents(events []EvidenceEvent, eventType string) []EvidenceEvent {
	out := []EvidenceEvent{}
	for _, event := range events {
		if event.EventType != eventType {
			out = append(out, event)
		}
	}
	return out
}

func validInput() Input {
	return Input{
		Contract: Contract{
			RequiredEventTypes: []string{
				"harness_lifecycle_observed",
				"tool_call_observed",
				"file_mutation_observed",
				"test_observed",
			},
		},
		Policy: Policy{
			PolicyID: "managed-policy-v1",
			PolicyProvenance: Provenance{
				Source: "vcs",
				Digest: "policy-digest",
			},
			AuthorizedAdapters: []AuthorizedAdapter{{
				AdapterID:       "opencode-adapter",
				HarnessID:       "opencode",
				AuthorityRef:    "sig:adapter",
				CapabilityIDs:   []string{"harness-events", "tool-events", "file-events", "test-events"},
				DeploymentRef:   "vcs:adapter",
				VersionRequired: "1.0.0",
			}},
			RequiredEventGroups: []RequiredEventGroup{
				{ID: "harness", EventTypes: []string{"harness_lifecycle_observed"}, AcceptableProvenanceScopes: []string{"harness_observed"}},
				{ID: "tool", EventTypes: []string{"tool_call_observed"}, AcceptableProvenanceScopes: []string{"local_observed"}},
				{ID: "file", EventTypes: []string{"file_mutation_observed"}, AcceptableProvenanceScopes: []string{"local_observed"}},
				{ID: "test", EventTypes: []string{"test_observed"}, AcceptableProvenanceScopes: []string{"local_observed", "ci_witnessed"}},
			},
		},
		Registry: Registry{
			RegistryID: "adapter-registry-v1",
			Provenance: Provenance{
				Source: "vcs",
				Digest: "registry-digest",
			},
			Adapters: []Adapter{{
				AdapterID:      "opencode-adapter",
				HarnessID:      "opencode",
				Version:        "1.0.0",
				DeploymentRef:  "vcs:adapter",
				IdentityState:  IdentityVerified,
				AuthorityRef:   "sig:adapter",
				AllowedEvents:  []string{"harness_lifecycle_observed", "tool_call_observed", "file_mutation_observed", "test_observed"},
				CapabilityRefs: []string{"harness-events", "tool-events", "file-events", "test-events"},
				Capabilities: []Capability{
					{ID: "harness-events", EventTypes: []string{"harness_lifecycle_observed"}, ProvenanceScope: "harness_observed"},
					{ID: "tool-events", EventTypes: []string{"tool_call_observed"}, ProvenanceScope: "local_observed"},
					{ID: "file-events", EventTypes: []string{"file_mutation_observed"}, ProvenanceScope: "local_observed"},
					{ID: "test-events", EventTypes: []string{"test_observed"}, ProvenanceScope: "local_observed"},
				},
			}},
		},
		Run: RunEvidence{
			RunID:           "managed-run-1",
			RunNonce:        "nonce-1",
			SourceCommit:    "abc123",
			ChainHead:       "chain-head",
			EventCount:      8,
			OutputArtifacts: []ArtifactDigest{{Path: "run.json", SHA256: "run-digest"}},
			ManagedBoundaryEnrolled: &ManagedBoundaryEnrolled{
				Sequence:              1,
				ManagedPolicyDigest:   "policy-digest",
				AdapterRegistryDigest: "registry-digest",
				AdapterID:             "opencode-adapter",
				EnrollmentSource:      "vcs",
				ManagedProfileID:      "managed-harness",
				ParentRunID:           "managed-run-1",
				RunNonce:              "nonce-1",
				EventDigest:           "enroll-digest",
			},
			ChildLaunch: LaunchEvent{
				Sequence:    2,
				EventDigest: "launch-digest",
			},
			ObservedEvents: []EvidenceEvent{
				{EventType: "harness_lifecycle_observed", ProvenanceScope: "harness_observed"},
				{EventType: "tool_call_observed", ProvenanceScope: "local_observed"},
				{EventType: "file_mutation_observed", ProvenanceScope: "local_observed"},
				{EventType: "test_observed", ProvenanceScope: "local_observed"},
			},
			TestEvidence: []EvidenceEvent{{EventType: "test_observed", ProvenanceScope: "local_observed"}},
		},
		Witness: Witness{
			WitnessID:             "ci-witness-1",
			Status:                StatePass,
			RunID:                 "managed-run-1",
			RunNonce:              "nonce-1",
			SourceCommit:          "abc123",
			ManagedPolicyDigest:   "policy-digest",
			AdapterRegistryDigest: "registry-digest",
			AdapterIdentityDigest: "opencode-adapter:sig:adapter",
			EnrollmentEventDigest: "enroll-digest",
			LaunchEventDigest:     "launch-digest",
			ChainHead:             "chain-head",
			EventCount:            8,
			FreshnessState:        StatePass,
			ArtifactDigests:       []ArtifactDigest{{Path: "run.json", SHA256: "run-digest"}},
		},
	}
}
