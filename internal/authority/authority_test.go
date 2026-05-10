package authority

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluateOutsideAuthorityDeniedTarget(t *testing.T) {
	pkg := validPackage()
	result := Evaluate(pkg)
	if result.AuthorityEvaluationState != StateOutsideAuthority {
		t.Fatalf("state = %s reasons=%v", result.AuthorityEvaluationState, result.Reasons)
	}
	eval := result.Evaluations[0]
	if eval.State != StateOutsideAuthority || eval.ReasonCode != "target_event_denied" {
		t.Fatalf("evaluation = %+v", eval)
	}
	if eval.ActorAttribution != AttributionVerified || eval.ToolAttribution != AttributionVerified {
		t.Fatalf("attribution = %+v", eval)
	}
	if eval.ModelAttribution != AttributionNotAssessed {
		t.Fatalf("model attribution = %s", eval.ModelAttribution)
	}
}

func TestEvaluateGitOnlyKeepsAttributionNotAssessed(t *testing.T) {
	pkg := validPackage()
	pkg.AuthorityEnvelopes[0].TargetRules[0].DeniedEvents = nil
	pkg.AuthorityEnvelopes[0].TargetRules[0].AllowedEvents = []string{"direct_mutation"}
	pkg.ObservedActions[0].SourceType = "git"
	result := Evaluate(pkg)
	eval := result.Evaluations[0]
	if eval.State != StateWithinAuthority {
		t.Fatalf("state = %s reason=%s", eval.State, eval.ReasonCode)
	}
	for _, got := range []string{eval.ActorAttribution, eval.ToolAttribution, eval.ModelAttribution} {
		if got != AttributionNotAssessed {
			t.Fatalf("unexpected attribution state in %+v", eval)
		}
	}
}

func TestEvaluateMissingApprovalEvidenceIsOutsideAuthority(t *testing.T) {
	pkg := validPackage()
	pkg.AuthorityEnvelopes[0].TargetRules[0].DeniedEvents = nil
	pkg.AuthorityEnvelopes[0].TargetRules[0].AllowedEvents = []string{"direct_mutation"}
	pkg.AuthorityEnvelopes[0].ApprovalRequirements = []ApprovalRequirement{{
		RequirementID: "approval-1",
		EventType:     "direct_mutation",
		TargetRuleRef: "rule-ci-denied",
	}}
	result := Evaluate(pkg)
	if result.AuthorityEvaluationState != StateOutsideAuthority {
		t.Fatalf("state = %s reasons=%v", result.AuthorityEvaluationState, result.Reasons)
	}
	if result.Evaluations[0].ReasonCode != "approval_evidence_missing" {
		t.Fatalf("evaluation = %+v", result.Evaluations[0])
	}
}

func TestEvaluateGatewayWithoutBindingDoesNotAttributeMutationModel(t *testing.T) {
	pkg := validPackage()
	pkg.AuthorityEnvelopes[0].TargetRules[0].DeniedEvents = nil
	pkg.AuthorityEnvelopes[0].TargetRules[0].AllowedEvents = []string{"direct_mutation"}
	pkg.ObservedActions = append(pkg.ObservedActions, ObservedAction{
		EventID:               "event-gateway-1",
		TaskID:                "task-1",
		EventType:             "gateway_request",
		SourceType:            "llm_gateway",
		EvidenceRefs:          []string{"artifact:gateway#requests/1"},
		ObservedAt:            "2026-05-09T00:00:01Z",
		ObservationConfidence: "single_source",
	})
	result := Evaluate(pkg)
	for _, eval := range result.Evaluations {
		if eval.EventID == "event-mutation-1" && eval.ModelAttribution != AttributionNotAssessed {
			t.Fatalf("mutation model attribution = %s", eval.ModelAttribution)
		}
	}
}

func TestEvaluateGatewaySourcedMutationWithoutBindingDoesNotVerifyModel(t *testing.T) {
	pkg := validPackage()
	pkg.AuthorityEnvelopes[0].TargetRules[0].DeniedEvents = nil
	pkg.AuthorityEnvelopes[0].TargetRules[0].AllowedEvents = []string{"direct_mutation"}
	pkg.ObservedActions[0].SourceType = "llm_gateway"
	result := Evaluate(pkg)
	eval := result.Evaluations[0]
	if eval.State != StateWithinAuthority {
		t.Fatalf("state = %s reason=%s", eval.State, eval.ReasonCode)
	}
	if eval.ModelAttribution != AttributionNotAssessed {
		t.Fatalf("model attribution = %s", eval.ModelAttribution)
	}
}

func TestEvaluateHarnessGatewayBindingCanAttributeModel(t *testing.T) {
	pkg := validPackage()
	pkg.AuthorityEnvelopes[0].TargetRules[0].DeniedEvents = nil
	pkg.AuthorityEnvelopes[0].TargetRules[0].AllowedEvents = []string{"direct_mutation"}
	pkg.ObservedActions[0].SourceType = "harness_log"
	pkg.ObservedActions = append(pkg.ObservedActions, ObservedAction{
		EventID:               "event-gateway-1",
		TaskID:                "task-1",
		EventType:             "gateway_request",
		SourceType:            "llm_gateway",
		EvidenceRefs:          []string{"artifact:gateway#requests/1"},
		ObservedAt:            "2026-05-09T00:00:01Z",
		ObservationConfidence: "single_source",
	})
	pkg.EvidenceBindings = []EvidenceBindingInput{{
		BindingID:     "binding-1",
		LeftEventID:   "event-mutation-1",
		RightEventID:  "event-gateway-1",
		BindingType:   "same_gateway_request",
		BindingState:  BindingVerified,
		MatchedFields: []string{"operation_id"},
		EvidenceRef:   "artifact:harness#binding/1",
	}}
	result := Evaluate(pkg)
	for _, eval := range result.Evaluations {
		if eval.EventID == "event-mutation-1" && eval.ModelAttribution != AttributionVerified {
			t.Fatalf("mutation model attribution = %s", eval.ModelAttribution)
		}
	}
}

func TestEvaluateMissingPolicyIsNotAssessed(t *testing.T) {
	pkg := validPackage()
	pkg.SelectedPolicyID = ""
	result := Evaluate(pkg)
	if result.AuthorityEvaluationState != StateNotAssessed {
		t.Fatalf("state = %s", result.AuthorityEvaluationState)
	}
	if result.Evaluations[0].ReasonCode != "policy_not_selected" {
		t.Fatalf("evaluation = %+v", result.Evaluations[0])
	}
}

func TestEvaluateConflictingEnvelopeCannotVerify(t *testing.T) {
	pkg := validPackage()
	pkg.AuthorityEnvelopes[0].AllowedEvents = []string{"direct_mutation"}
	pkg.AuthorityEnvelopes[0].DeniedEvents = []string{"direct_mutation"}
	result := Evaluate(pkg)
	if result.AuthorityEvaluationState != StateCannotVerify {
		t.Fatalf("state = %s", result.AuthorityEvaluationState)
	}
	if result.Evaluations[0].ReasonCode != "allow_deny_event_conflict" {
		t.Fatalf("evaluation = %+v", result.Evaluations[0])
	}
}

func TestEvaluateFailedBindingCannotVerify(t *testing.T) {
	pkg := validPackage()
	pkg.EvidenceBindings = []EvidenceBindingInput{{
		BindingID:     "binding-1",
		LeftEventID:   "event-mutation-1",
		RightEventID:  "event-gateway-1",
		BindingType:   "same_gateway_request",
		BindingState:  BindingCannotVerify,
		MatchedFields: []string{"operation_id"},
		EvidenceRef:   "artifact:harness#binding/1",
	}}
	pkg.ObservedActions = append(pkg.ObservedActions, ObservedAction{
		EventID:               "event-gateway-1",
		TaskID:                "task-1",
		EventType:             "gateway_request",
		SourceType:            "llm_gateway",
		EvidenceRefs:          []string{"artifact:gateway#requests/1"},
		ObservedAt:            "2026-05-09T00:00:01Z",
		ObservationConfidence: "single_source",
	})
	result := Evaluate(pkg)
	if result.AuthorityEvaluationState != StateCannotVerify {
		t.Fatalf("state = %s reasons=%v", result.AuthorityEvaluationState, result.Reasons)
	}
}

func TestEvaluateOverlappingTargetRulesCannotVerify(t *testing.T) {
	pkg := validPackage()
	pkg.AuthorityEnvelopes[0].AllowedEvents = nil
	pkg.AuthorityEnvelopes[0].TargetRules = []TargetRule{
		{
			RuleID:        "rule-docs-allow",
			TargetPattern: "docs/**",
			AllowedEvents: []string{"direct_mutation"},
		},
		{
			RuleID:        "rule-md-deny",
			TargetPattern: "**/*.md",
			DeniedEvents:  []string{"direct_mutation"},
		},
	}
	pkg.ObservedActions[0].Target = "docs/authority-envelope.md"
	result := Evaluate(pkg)
	if result.AuthorityEvaluationState != StateCannotVerify {
		t.Fatalf("state = %s reasons=%v", result.AuthorityEvaluationState, result.Reasons)
	}
	if result.Evaluations[0].ReasonCode != "overlapping_target_rules_conflict" {
		t.Fatalf("evaluation = %+v", result.Evaluations[0])
	}
}

func TestValidateEnvelopeReasons(t *testing.T) {
	base := validPackage().AuthorityEnvelopes[0]
	tests := []struct {
		name   string
		mutate func(*AuthorityEnvelope)
		want   string
	}{
		{
			name:   "missing required field",
			mutate: func(env *AuthorityEnvelope) { env.ActorRef = "" },
			want:   "authority_envelope_malformed",
		},
		{
			name:   "unsupported top level event",
			mutate: func(env *AuthorityEnvelope) { env.AllowedEvents = []string{"unknown"} },
			want:   "unsupported_event_type",
		},
		{
			name: "top level allow deny conflict",
			mutate: func(env *AuthorityEnvelope) {
				env.AllowedEvents = []string{"review"}
				env.DeniedEvents = []string{"review"}
			},
			want: "allow_deny_event_conflict",
		},
		{
			name:   "target rule malformed",
			mutate: func(env *AuthorityEnvelope) { env.TargetRules[0].RuleID = "" },
			want:   "target_rule_malformed",
		},
		{
			name:   "target unsupported event",
			mutate: func(env *AuthorityEnvelope) { env.TargetRules[0].DeniedEvents = []string{"unknown"} },
			want:   "target_rule_conflict",
		},
		{
			name: "target allow deny conflict",
			mutate: func(env *AuthorityEnvelope) {
				env.TargetRules[0].AllowedEvents = []string{"direct_mutation"}
				env.TargetRules[0].DeniedEvents = []string{"direct_mutation"}
			},
			want: "target_rule_conflict",
		},
		{
			name: "target allowed conflicts with top level denial",
			mutate: func(env *AuthorityEnvelope) {
				env.AllowedEvents = nil
				env.DeniedEvents = []string{"direct_mutation"}
				env.TargetRules[0].AllowedEvents = []string{"direct_mutation"}
				env.TargetRules[0].DeniedEvents = nil
			},
			want: "target_rule_conflicts_with_top_level",
		},
		{
			name: "target denied conflicts with top level allow",
			mutate: func(env *AuthorityEnvelope) {
				env.AllowedEvents = []string{"direct_mutation"}
				env.DeniedEvents = nil
				env.TargetRules[0].AllowedEvents = nil
				env.TargetRules[0].DeniedEvents = []string{"direct_mutation"}
			},
			want: "target_rule_conflicts_with_top_level",
		},
		{
			name: "overlapping target rules conflict",
			mutate: func(env *AuthorityEnvelope) {
				env.AllowedEvents = nil
				env.TargetRules = append(env.TargetRules, TargetRule{
					RuleID:        "rule-ci-allow",
					TargetPattern: ".github/workflows/**",
					AllowedEvents: []string{"direct_mutation"},
				})
			},
			want: "overlapping_target_rules_conflict",
		},
		{
			name: "valid envelope",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := base
			env.TargetRules = append([]TargetRule(nil), base.TargetRules...)
			if tt.mutate != nil {
				tt.mutate(&env)
			}
			if got := validateEnvelope(env); got != tt.want {
				t.Fatalf("reason = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEvaluateExternalAndStaleEvidenceCannotVerify(t *testing.T) {
	for name, mutate := range map[string]func(*Package){
		"unresolved-external": func(pkg *Package) {
			pkg.ObservedActions[0].EvidenceRefs = []string{"external:ticket-1"}
		},
		"stale-file-ref": func(pkg *Package) {
			pkg.EvidenceResolution.StaleRefs = []string{"file:evidence/git-diff.json"}
		},
	} {
		t.Run(name, func(t *testing.T) {
			pkg := validPackage()
			mutate(&pkg)
			result := Evaluate(pkg)
			if result.AuthorityEvaluationState != StateCannotVerify {
				t.Fatalf("state = %s reasons=%v", result.AuthorityEvaluationState, result.Reasons)
			}
		})
	}
}

func TestEvaluateActionReasonOrdering(t *testing.T) {
	pkg := validPackage()
	env := pkg.AuthorityEnvelopes[0]
	action := pkg.ObservedActions[0]
	tests := []struct {
		name        string
		envState    string
		envReason   string
		mutate      func(*AuthorityEnvelope, *ObservedAction)
		bindings    []EvidenceBinding
		resolution  map[string]string
		wantState   string
		wantReason  string
		matchedRule string
	}{
		{name: "envelope state first", envState: StateCannotVerify, envReason: "selected_policy_ambiguous", wantState: StateCannotVerify, wantReason: "selected_policy_ambiguous"},
		{name: "unsupported event", mutate: func(_ *AuthorityEnvelope, action *ObservedAction) { action.EventType = "unknown" }, wantState: StateCannotVerify, wantReason: "unsupported_event_type"},
		{name: "task not assessed", mutate: func(_ *AuthorityEnvelope, action *ObservedAction) { action.TaskID = "" }, wantState: StateNotAssessed, wantReason: "task_not_assessed"},
		{name: "task outside envelope", mutate: func(_ *AuthorityEnvelope, action *ObservedAction) { action.TaskID = "other" }, wantState: StateNotAssessed, wantReason: "task_outside_selected_envelope"},
		{name: "external evidence", mutate: func(_ *AuthorityEnvelope, action *ObservedAction) {
			action.EvidenceRefs = []string{"external:ticket-1"}
		}, wantState: StateCannotVerify, wantReason: "external_evidence_unresolved"},
		{
			name:       "binding cannot verify",
			bindings:   []EvidenceBinding{{BindingState: BindingCannotVerify}},
			wantState:  StateCannotVerify,
			wantReason: "evidence_binding_cannot_verify",
		},
		{
			name:        "matched denial",
			wantState:   StateOutsideAuthority,
			wantReason:  "target_event_denied",
			matchedRule: "rule-ci-denied",
		},
		{
			name: "approval evidence missing",
			mutate: func(env *AuthorityEnvelope, _ *ObservedAction) {
				env.TargetRules[0].DeniedEvents = nil
				env.TargetRules[0].AllowedEvents = []string{"direct_mutation"}
				env.ApprovalRequirements = []ApprovalRequirement{{
					RequirementID: "approval-1",
					EventType:     "direct_mutation",
					TargetRuleRef: "rule-ci-denied",
				}}
			},
			wantState:   StateOutsideAuthority,
			wantReason:  "approval_evidence_missing",
			matchedRule: "rule-ci-denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testAction := action
			testEnv := env
			if tt.mutate != nil {
				tt.mutate(&testEnv, &testAction)
			}
			eval := evaluateAction("eval-1", pkg.SelectedPolicyID, testEnv, tt.envState, tt.envReason, testAction, tt.bindings, tt.resolution)
			if eval.State != tt.wantState || eval.ReasonCode != tt.wantReason {
				t.Fatalf("evaluation = %+v, want %s/%s", eval, tt.wantState, tt.wantReason)
			}
			if tt.matchedRule != "" && eval.MatchedRuleRef != tt.matchedRule {
				t.Fatalf("matched rule = %q, want %q", eval.MatchedRuleRef, tt.matchedRule)
			}
		})
	}
}

func TestEvaluateSameActionDifferentPolicies(t *testing.T) {
	pkg := validPackage()
	allow := pkg.AuthorityEnvelopes[0]
	allow.PolicyID = "policy-allow"
	allow.TargetRules = append([]TargetRule(nil), allow.TargetRules...)
	allow.TargetRules[0].AllowedEvents = []string{"direct_mutation"}
	allow.TargetRules[0].DeniedEvents = nil
	deny := pkg.AuthorityEnvelopes[0]
	deny.PolicyID = "policy-deny"
	deny.TargetRules = append([]TargetRule(nil), deny.TargetRules...)
	pkg.AuthorityEnvelopes = []AuthorityEnvelope{allow, deny}

	pkg.SelectedPolicyID = "policy-allow"
	if got := Evaluate(pkg).AuthorityEvaluationState; got != StateWithinAuthority {
		t.Fatalf("allow policy state = %s", got)
	}
	pkg.SelectedPolicyID = "policy-deny"
	if got := Evaluate(pkg).AuthorityEvaluationState; got != StateOutsideAuthority {
		t.Fatalf("deny policy state = %s", got)
	}
}

func TestFixtureMatrixScenarios(t *testing.T) {
	root := filepath.Join("..", "..", "examples", "authority-envelope-basic")
	data, err := os.ReadFile(filepath.Join(root, "fixture-matrix.json"))
	if err != nil {
		t.Fatalf("read fixture matrix: %v", err)
	}
	var matrix struct {
		Scenarios []struct {
			Name          string `json:"name"`
			Path          string `json:"path"`
			ExpectedState string `json:"expected_state"`
			ReasonCode    string `json:"reason_code"`
		} `json:"scenarios"`
	}
	if err := json.Unmarshal(data, &matrix); err != nil {
		t.Fatalf("decode fixture matrix: %v", err)
	}
	for _, scenario := range matrix.Scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			pkg, err := ReadPackage(filepath.Join(root, scenario.Path))
			if err != nil {
				t.Fatalf("read package: %v", err)
			}
			result := Evaluate(pkg)
			if result.AuthorityEvaluationState != scenario.ExpectedState {
				t.Fatalf("state = %s want %s reasons=%v", result.AuthorityEvaluationState, scenario.ExpectedState, result.Reasons)
			}
			if scenario.ReasonCode != "" && !hasReason(result, scenario.ReasonCode) {
				t.Fatalf("missing reason %s in %+v", scenario.ReasonCode, result)
			}
		})
	}
}

func hasReason(result Result, reason string) bool {
	for _, got := range result.Reasons {
		if got == reason {
			return true
		}
	}
	return false
}

func validPackage() Package {
	return Package{
		SchemaVersion:    PackageSchemaVersion,
		SelectedPolicyID: "policy-deny-ci",
		Actors: []ActorDeclaration{{
			ActorID:      "agent-1",
			ActorType:    "ai_agent",
			DeclaredRole: "observer",
			Harness:      "generic-harness",
			OperationID:  "op-1",
		}},
		AuthorityEnvelopes: []AuthorityEnvelope{{
			SchemaVersion:  "authority-envelope-v1",
			TaskID:         "task-1",
			PolicyID:       "policy-deny-ci",
			AuthorityScope: "task",
			ActorRef:       "agent-1",
			AllowedEvents:  []string{"review", "feedback"},
			DeniedEvents:   []string{},
			TargetRules: []TargetRule{{
				RuleID:        "rule-ci-denied",
				TargetPattern: ".github/workflows/**",
				AllowedEvents: nil,
				DeniedEvents:  []string{"direct_mutation"},
			}},
		}},
		ObservedActions: []ObservedAction{{
			EventID:               "event-mutation-1",
			TaskID:                "task-1",
			EventType:             "direct_mutation",
			Target:                ".github/workflows/verify.yml",
			SourceType:            "harness_log",
			EvidenceRefs:          []string{"file:evidence/git-diff.json"},
			ActorID:               "agent-1",
			OperationID:           "op-1",
			ObservedAt:            "2026-05-09T00:00:00Z",
			ObservationConfidence: "single_source",
		}},
	}
}
