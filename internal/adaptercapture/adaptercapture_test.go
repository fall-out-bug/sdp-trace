package adaptercapture

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluatePassesForBoundGenericAdapterEvents(t *testing.T) {
	result := Evaluate(validInput())
	if result.AdapterCaptureAssessment != StatePass {
		t.Fatalf("assessment = %s reasons=%v", result.AdapterCaptureAssessment, result.Reasons)
	}
	for _, id := range adapterConditionIDs {
		if conditionByID(result.AdapterCaptureConditions, id).State != StatePass {
			t.Fatalf("condition %s = %+v", id, conditionByID(result.AdapterCaptureConditions, id))
		}
	}
}

func TestEvaluateCannotVerifyMissingRequiredAdapterEvent(t *testing.T) {
	input := validInput()
	input.Run.AdapterEvents = removeEvents(input.Run.AdapterEvents, "tool_call")
	result := Evaluate(input)
	if result.AdapterCaptureAssessment != StateCannotVerify {
		t.Fatalf("assessment = %s reasons=%v", result.AdapterCaptureAssessment, result.Reasons)
	}
	condition := conditionByID(result.AdapterCaptureConditions, "tool_call_depth_visible")
	if condition.State != StateMissingTelemetry || condition.ReasonCode != "tool_event_missing" {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestEvaluateFailsAgentReportedTestAsExecutedEvidence(t *testing.T) {
	input := validInput()
	for i := range input.Run.AdapterEvents {
		if input.Run.AdapterEvents[i].EventType == "test_observed" {
			input.Run.AdapterEvents[i].TestProvenance = "agent_reported"
			input.Run.AdapterEvents[i].ExecutedEvidenceClaimed = true
		}
	}
	result := Evaluate(input)
	if result.AdapterCaptureAssessment != StateFail {
		t.Fatalf("assessment = %s reasons=%v", result.AdapterCaptureAssessment, result.Reasons)
	}
	condition := conditionByID(result.AdapterCaptureConditions, "test_provenance_not_overclaimed")
	if condition.State != StateFail || condition.ReasonCode != "agent_reported_test_not_executed" {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestEvaluateCannotVerifyLateAdapterEvent(t *testing.T) {
	input := validInput()
	input.Run.AdapterEvents[1].Sequence = input.Run.RunClosedSequence + 1
	result := Evaluate(input)
	if result.AdapterCaptureAssessment != StateCannotVerify {
		t.Fatalf("assessment = %s reasons=%v", result.AdapterCaptureAssessment, result.Reasons)
	}
	condition := conditionByID(result.AdapterCaptureConditions, "run_binding_established")
	if condition.State != StateCannotVerify || condition.ReasonCode != "late_adapter_event" {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestEvaluateFailsUnsafeAdapterMetadata(t *testing.T) {
	input := validInput()
	input.Run.AdapterEvents[0].SensitiveMetadataPersisted = true
	input.Run.AdapterEvents[0].GatewayEvidenceRef = "https://example.invalid/review?token=secret-token"
	input.Run.AdapterEvents[0].ProviderRefs = []string{"https://review.invalid/7?access_token=secret-token"}
	result := Evaluate(input)
	if result.AdapterCaptureAssessment != StateFail {
		t.Fatalf("assessment = %s reasons=%v", result.AdapterCaptureAssessment, result.Reasons)
	}
	condition := conditionByID(result.AdapterCaptureConditions, "redaction_metadata_consistent")
	if condition.State != StateFail || condition.ReasonCode != "forbidden_adapter_metadata_persisted" {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestEvaluateCannotVerifyFileMutationWithoutSourceCorrelation(t *testing.T) {
	input := validInput()
	for i := range input.Run.AdapterEvents {
		if input.Run.AdapterEvents[i].EventType == "file_mutation" {
			input.Run.AdapterEvents[i].SourceBaseline = ""
		}
	}
	result := Evaluate(input)
	if result.AdapterCaptureAssessment != StateCannotVerify {
		t.Fatalf("assessment = %s reasons=%v", result.AdapterCaptureAssessment, result.Reasons)
	}
	condition := conditionByID(result.AdapterCaptureConditions, "file_mutation_correlated")
	if condition.State != StateCannotVerify || condition.ReasonCode != "file_mutation_source_missing" {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestEvaluateFailsEventLevelProviderRefSecret(t *testing.T) {
	input := validInput()
	input.Run.AdapterEvents[0].ProviderRefs = []string{"https://review.invalid/7?session_id=secret-token"}
	result := Evaluate(input)
	if result.AdapterCaptureAssessment != StateFail {
		t.Fatalf("assessment = %s reasons=%v", result.AdapterCaptureAssessment, result.Reasons)
	}
	condition := conditionByID(result.AdapterCaptureConditions, "provider_refs_portable")
	if condition.State != StateFail || condition.ReasonCode != "provider_ref_contains_secret" {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestEvaluateFailsCaptureDepthOverclaim(t *testing.T) {
	input := validInput()
	input.Run.EventFamilySummaries = []EventFamilyState{{
		EventFamily:     "tool_call",
		State:           StateMissingTelemetry,
		RetentionMode:   RetentionDigestOnly,
		Reconstructable: true,
	}}
	result := Evaluate(input)
	if result.AdapterCaptureAssessment != StateFail {
		t.Fatalf("assessment = %s reasons=%v", result.AdapterCaptureAssessment, result.Reasons)
	}
	condition := conditionByID(result.AdapterCaptureConditions, "capture_depth_not_overclaimed")
	if condition.State != StateFail || condition.ReasonCode != "capture_depth_overclaimed" {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestEvaluateCannotVerifyConflictingCorrelationKeys(t *testing.T) {
	input := validInput()
	conflict := input.Run.AdapterEvents[2]
	conflict.EventID = "evt-tool-conflict"
	conflict.Sequence = 9
	conflict.EventHash = "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	input.Run.AdapterEvents = append(input.Run.AdapterEvents, conflict)
	result := Evaluate(input)
	if result.AdapterCaptureAssessment != StateCannotVerify {
		t.Fatalf("assessment = %s reasons=%v", result.AdapterCaptureAssessment, result.Reasons)
	}
	condition := conditionByID(result.AdapterCaptureConditions, "adapter_event_contract_valid")
	if condition.State != StateCannotVerify || condition.ReasonCode != "conflicting_adapter_events" {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestBlock19CommittedFixturesHaveAdapterCaptureShape(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "examples", "block19-adapter-capture")
	cases := block19FixtureCases()
	if os.Getenv("SDP_TRACE_UPDATE_BLOCK19_FIXTURES") == "1" {
		if err := os.MkdirAll(fixtureDir, 0o755); err != nil {
			t.Fatalf("mkdir fixture dir: %v", err)
		}
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
	expected := map[string]block19FixtureCase{}
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
			t.Fatalf("unexpected Block 19 fixture %s", entry.Name())
		}
		payload, err := os.ReadFile(filepath.Join(fixtureDir, entry.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", entry.Name(), err)
		}
		var result AssessmentResult
		if err := json.Unmarshal(payload, &result); err != nil {
			t.Fatalf("unmarshal %s: %v", entry.Name(), err)
		}
		if result.SchemaVersion != SchemaVersion || result.SelectedProfile != ProfileAdapterCapture {
			t.Fatalf("%s schema/profile = %s/%s", entry.Name(), result.SchemaVersion, result.SelectedProfile)
		}
		if len(result.AdapterCaptureConditions) != len(adapterConditionIDs) {
			t.Fatalf("%s condition count = %d", entry.Name(), len(result.AdapterCaptureConditions))
		}
		if result.AdapterCaptureAssessment != topLevel(result.AdapterCaptureConditions) {
			t.Fatalf("%s assessment = %s, want %s from conditions", entry.Name(), result.AdapterCaptureAssessment, topLevel(result.AdapterCaptureConditions))
		}
		condition := conditionByID(result.AdapterCaptureConditions, fixture.conditionID)
		if condition.State != fixture.conditionState || condition.ReasonCode != fixture.reasonCode {
			t.Fatalf("%s condition %s = %s/%s, want %s/%s", entry.Name(), fixture.conditionID, condition.State, condition.ReasonCode, fixture.conditionState, fixture.reasonCode)
		}
		seen++
	}
	if seen != len(expected) {
		t.Fatalf("seen fixtures = %d, want %d", seen, len(expected))
	}
}

type block19FixtureCase struct {
	name           string
	mutate         func(*Input)
	conditionID    string
	conditionState string
	reasonCode     string
}

func (fixture block19FixtureCase) result() AssessmentResult {
	input := validInput()
	if fixture.mutate != nil {
		fixture.mutate(&input)
	}
	return Evaluate(input)
}

func block19FixtureCases() []block19FixtureCase {
	return []block19FixtureCase{
		{name: "valid-adapter-capture.assessment-result.json", conditionID: "run_binding_established", conditionState: StatePass, reasonCode: "run_binding_established"},
		{name: "adapter-bundle-binding-pass.assessment-result.json", mutate: useAdapterBundleBinding, conditionID: "run_binding_established", conditionState: StatePass, reasonCode: "run_binding_established"},
		{name: "missing-tool-event-cannot-verify.assessment-result.json", mutate: func(input *Input) {
			input.Run.AdapterEvents = removeEvents(input.Run.AdapterEvents, "tool_call")
		}, conditionID: "tool_call_depth_visible", conditionState: StateMissingTelemetry, reasonCode: "tool_event_missing"},
		{name: "unsupported-observer-cannot-verify.assessment-result.json", mutate: func(input *Input) {
			input.Run.AdapterEvents = removeEvents(input.Run.AdapterEvents, "tool_call")
			input.Run.UnsupportedEventTypes = []string{"tool_call"}
		}, conditionID: "tool_call_depth_visible", conditionState: StateUnsupported, reasonCode: "tool_event_unsupported"},
		{name: "gateway-not-integrated-cannot-verify.assessment-result.json", mutate: func(input *Input) {
			input.Run.GatewayIntegrated = false
			input.Run.GatewayEvidenceBound = false
			for i := range input.Run.AdapterEvents {
				if input.Run.AdapterEvents[i].EventType == "model_call_observed" {
					input.Run.AdapterEvents[i].ModelIdentityProvenance = "harness_observed"
					input.Run.AdapterEvents[i].IdentityBinding = IdentitySelfAsserted
				}
			}
		}, conditionID: "model_identity_not_overclaimed", conditionState: StateNotIntegrated, reasonCode: "gateway_not_integrated"},
		{name: "agent-reported-test-fail.assessment-result.json", mutate: func(input *Input) {
			for i := range input.Run.AdapterEvents {
				if input.Run.AdapterEvents[i].EventType == "test_observed" {
					input.Run.AdapterEvents[i].TestProvenance = "agent_reported"
					input.Run.AdapterEvents[i].ExecutedEvidenceClaimed = true
				}
			}
		}, conditionID: "test_provenance_not_overclaimed", conditionState: StateFail, reasonCode: "agent_reported_test_not_executed"},
		{name: "harness-observed-test-cannot-verify.assessment-result.json", mutate: func(input *Input) {
			for i := range input.Run.AdapterEvents {
				if input.Run.AdapterEvents[i].EventType == "test_observed" {
					input.Run.AdapterEvents[i].TestProvenance = "harness_observed"
					input.Run.AdapterEvents[i].ExecutedEvidenceClaimed = false
				}
			}
		}, conditionID: "test_provenance_not_overclaimed", conditionState: StateCannotVerify, reasonCode: "test_execution_unverified"},
		{name: "file-mutation-source-correlated-pass.assessment-result.json", conditionID: "file_mutation_correlated", conditionState: StatePass, reasonCode: "file_mutation_correlated"},
		{name: "task-supersession-attribution-pass.assessment-result.json", mutate: addBoundTaskSupersession, conditionID: "task_drift_visible", conditionState: StatePass, reasonCode: "task_supersessions_visible"},
		{name: "provider-neutral-refs-pass.assessment-result.json", conditionID: "provider_refs_portable", conditionState: StatePass, reasonCode: "provider_refs_portable"},
		{name: "late-adapter-event-cannot-verify.assessment-result.json", mutate: func(input *Input) {
			input.Run.AdapterEvents[1].Sequence = input.Run.RunClosedSequence + 1
		}, conditionID: "run_binding_established", conditionState: StateCannotVerify, reasonCode: "late_adapter_event"},
		{name: "conflicting-adapter-events-cannot-verify.assessment-result.json", mutate: func(input *Input) {
			conflict := input.Run.AdapterEvents[2]
			conflict.EventID = "evt-tool-conflict"
			conflict.Sequence = 9
			conflict.EventHash = "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
			input.Run.AdapterEvents = append(input.Run.AdapterEvents, conflict)
		}, conditionID: "adapter_event_contract_valid", conditionState: StateCannotVerify, reasonCode: "conflicting_adapter_events"},
		{name: "capture-depth-overclaimed-fail.assessment-result.json", mutate: func(input *Input) {
			input.Run.EventFamilySummaries = []EventFamilyState{{
				EventFamily:     "tool_call",
				State:           StateMissingTelemetry,
				RetentionMode:   RetentionDigestOnly,
				Reconstructable: true,
			}}
		}, conditionID: "capture_depth_not_overclaimed", conditionState: StateFail, reasonCode: "capture_depth_overclaimed"},
		{name: "unsafe-provider-ref-fail.assessment-result.json", mutate: func(input *Input) {
			input.Run.ProviderRefs = []ProviderRef{{ReviewRef: "https://review.invalid/7?token=secret-token"}}
		}, conditionID: "provider_refs_portable", conditionState: StateFail, reasonCode: "provider_ref_contains_secret"},
		{name: "unsafe-event-provider-ref-fail.assessment-result.json", mutate: func(input *Input) {
			input.Run.AdapterEvents[0].ProviderRefs = []string{"https://review.invalid/7?session_id=secret-token"}
		}, conditionID: "provider_refs_portable", conditionState: StateFail, reasonCode: "provider_ref_contains_secret"},
	}
}

func conditionByID(conditions []Condition, id string) Condition {
	for _, condition := range conditions {
		if condition.ID == id {
			return condition
		}
	}
	return Condition{ID: id}
}

func removeEvents(events []AdapterEvent, eventType string) []AdapterEvent {
	out := []AdapterEvent{}
	for _, event := range events {
		if event.EventType != eventType {
			out = append(out, event)
		}
	}
	return out
}

func useAdapterBundleBinding(input *Input) {
	head := "cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	input.Run.AdapterBundle = &AdapterBundle{
		BundleID:           "bundle-1",
		HeadDigest:         head,
		ReferencedSequence: input.Run.RunClosedSequence,
		EventCount:         len(input.Run.AdapterEvents),
	}
	for i := range input.Run.AdapterEvents {
		input.Run.AdapterEvents[i].BindingMode = BindingAdapterBundle
		input.Run.AdapterEvents[i].Sequence = 0
		input.Run.AdapterEvents[i].PrevEventHash = ""
		input.Run.AdapterEvents[i].EventHash = ""
		input.Run.AdapterEvents[i].AdapterBundleID = "bundle-1"
		input.Run.AdapterEvents[i].AdapterBundleHeadDigest = head
		input.Run.AdapterEvents[i].BundleSequence = i + 1
	}
}

func addBoundTaskSupersession(input *Input) {
	event := validEvent("evt-task-superseded", "task_superseded", 9, input.Run.RunID, input.Run.RunNonce, input.Run.SourceBaseline, input.Run.RedactionPolicyDigest)
	event.ActorAttributionState = "bound"
	input.Run.AdapterEvents = append(input.Run.AdapterEvents, event)
	input.Run.TaskSupersessionCount = 1
}
