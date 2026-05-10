package forensic

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEvaluatePassesWithSanitizedAndRawReferences(t *testing.T) {
	result := Evaluate(validInput())
	if result.ForensicRetentionAssessment != StatePass {
		t.Fatalf("assessment = %s reasons=%v", result.ForensicRetentionAssessment, result.Reasons)
	}
	for _, id := range forensicConditionIDs {
		if conditionByID(result.ForensicConditions, id).State != StatePass {
			t.Fatalf("condition %s = %+v", id, conditionByID(result.ForensicConditions, id))
		}
	}
}

func TestEvaluateRejectsDigestOnlyCriticalEvidenceWithCap(t *testing.T) {
	input := validInput()
	input.Run.Events[0].RetentionMode = RetentionModeDigestOnly
	input.Run.Events[0].RawReference = nil
	result := Evaluate(input)
	if result.ForensicRetentionAssessment != StateFail {
		t.Fatalf("assessment = %s reasons=%v", result.ForensicRetentionAssessment, result.Reasons)
	}
	condition := conditionByID(result.ForensicConditions, "critical_evidence_reconstructable")
	if condition.State != StateFail || condition.ReasonCode != "critical_evidence_digest_only" || condition.CappedToRetentionMode != RetentionModeDigestOnly {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestEvaluateCannotVerifyUnverifiableRawReferenceAccess(t *testing.T) {
	input := validInput()
	input.Run.Events[1].RawReference.AccessState = AccessStateRestricted
	input.Run.Events[1].RawReference.AccessStateLastVerified = ""
	result := Evaluate(input)
	if result.ForensicRetentionAssessment != StateCannotVerify {
		t.Fatalf("assessment = %s reasons=%v", result.ForensicRetentionAssessment, result.Reasons)
	}
	condition := conditionByID(result.ForensicConditions, "raw_reference_bound")
	if condition.State != StateCannotVerify || condition.ReasonCode != "access_unverifiable" {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestEvaluateRawReferenceReasonCodes(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*RawReference)
		state  string
		reason string
	}{
		{
			name: "invalid reference type",
			mutate: func(ref *RawReference) {
				ref.ReferenceType = RetentionModeDigestOnly
			},
			state:  StateFail,
			reason: "raw_reference_type_invalid",
		},
		{
			name: "missing reference uri",
			mutate: func(ref *RawReference) {
				ref.ReferenceURI = ""
			},
			state:  StateCannotVerify,
			reason: "missing_reference",
		},
		{
			name: "encrypted key custody unknown",
			mutate: func(ref *RawReference) {
				ref.ReferenceType = RetentionModeEncryptedRawRef
				ref.KeyCustodyState = KeyCustodyUnknown
			},
			state:  StateCannotVerify,
			reason: "key_custody_unverifiable",
		},
		{
			name: "retention lifecycle expired",
			mutate: func(ref *RawReference) {
				ref.RetentionLifecycle.State = RetentionLifecycleExpired
			},
			state:  StateCannotVerify,
			reason: "retention_lifecycle_unverifiable",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := validInput()
			tt.mutate(input.Run.Events[1].RawReference)
			result := Evaluate(input)
			condition := conditionByID(result.ForensicConditions, "raw_reference_bound")
			if condition.State != tt.state || condition.ReasonCode != tt.reason {
				t.Fatalf("condition = %+v", condition)
			}
		})
	}
}

func TestEvaluateFailsWeakDigestAndCannotVerifySelfClaimedAuthority(t *testing.T) {
	input := validInput()
	input.Run.Events[1].RawReference.Digest.Algorithm = "sha1"
	input.Run.Events[0].RedactionAuthority.VerificationState = AuthoritySelfClaimed
	result := Evaluate(input)
	if result.ForensicRetentionAssessment != StateFail {
		t.Fatalf("assessment = %s reasons=%v", result.ForensicRetentionAssessment, result.Reasons)
	}
	if conditionByID(result.ForensicConditions, "raw_reference_bound").ReasonCode != "weak_digest" {
		t.Fatalf("raw condition = %+v", conditionByID(result.ForensicConditions, "raw_reference_bound"))
	}
	policyCondition := conditionByID(result.ForensicConditions, "redaction_policy_bound")
	if policyCondition.State != StateCannotVerify || policyCondition.ReasonCode != "authority_self_claimed" {
		t.Fatalf("policy condition = %+v", conditionByID(result.ForensicConditions, "redaction_policy_bound"))
	}
}

func TestEvaluateWithholdMapsToNotAssessedWithoutRawReference(t *testing.T) {
	input := validInput()
	input.Run.Events[0].RedactionAction = RedactionActionWithhold
	input.Run.Events[0].RedactionRuleRefs = []string{"withhold-privacy-v1"}
	input.Run.Events[0].RetentionMode = RetentionModeNotAssessed
	input.Run.Events[0].Withholding = &Withholding{
		Authority:     AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
		Requestor:     AuthorityRef{ActorID: "human:privacy-owner", VerificationState: AuthorityVerified},
		ReasonCode:    "privacy_request",
		Justification: "privacy legal hold",
	}
	input.Run.Events[0].RawReference = nil
	result := Evaluate(input)
	if result.ForensicRetentionAssessment != StateFail {
		t.Fatalf("assessment = %s reasons=%v", result.ForensicRetentionAssessment, result.Reasons)
	}
	condition := conditionByID(result.ForensicConditions, "critical_evidence_reconstructable")
	if condition.State != StateFail || condition.ReasonCode != "critical_evidence_not_assessed" || condition.CappedToRetentionMode != RetentionModeNotAssessed {
		t.Fatalf("condition = %+v", condition)
	}
}

func TestEvaluateCannotVerifyMissingRuleRefsAndWithholdingAudit(t *testing.T) {
	input := validInput()
	input.Run.Events[0].RedactionRuleRefs = nil
	input.Run.Events[1].RedactionAction = RedactionActionWithhold
	input.Run.Events[1].RedactionRuleRefs = []string{"withhold-privacy-v1"}
	input.Run.Events[1].Withholding = nil
	result := Evaluate(input)
	if result.ForensicRetentionAssessment != StateCannotVerify {
		t.Fatalf("assessment = %s reasons=%v", result.ForensicRetentionAssessment, result.Reasons)
	}
	if conditionByID(result.ForensicConditions, "redaction_prewrite_applied").ReasonCode != "redaction_rule_refs_missing" {
		t.Fatalf("prewrite condition = %+v", conditionByID(result.ForensicConditions, "redaction_prewrite_applied"))
	}
	if conditionByID(result.ForensicConditions, "redaction_unresolved_visible").ReasonCode != "withholding_audit_missing" {
		t.Fatalf("unresolved condition = %+v", conditionByID(result.ForensicConditions, "redaction_unresolved_visible"))
	}
}

func TestEvaluateFailsUnknownRuleRefAndPolicyDisallowedRetentionMode(t *testing.T) {
	input := validInput()
	input.Run.Events[0].RedactionRuleRefs = []string{"missing-rule"}
	input.Run.Events[1].RetentionMode = RetentionModeDigestOnly
	input.Policy.AllowedRetentionModes = []string{RetentionModeSanitizedExcerpt, RetentionModeEncryptedRawRef, RetentionModeExternalArtifactRef, RetentionModeNotAssessed}
	result := Evaluate(input)
	if result.ForensicRetentionAssessment != StateFail {
		t.Fatalf("assessment = %s reasons=%v", result.ForensicRetentionAssessment, result.Reasons)
	}
	if conditionByID(result.ForensicConditions, "redaction_prewrite_applied").ReasonCode != "redaction_rule_unknown" {
		t.Fatalf("prewrite condition = %+v", conditionByID(result.ForensicConditions, "redaction_prewrite_applied"))
	}
	if conditionByID(result.ForensicConditions, "retention_mode_declared").ReasonCode != "retention_mode_not_policy_allowed" {
		t.Fatalf("retention condition = %+v", conditionByID(result.ForensicConditions, "retention_mode_declared"))
	}
}

func TestBlock18CommittedFixturesHaveForensicAssessmentShape(t *testing.T) {
	fixtureDir := filepath.Join("..", "..", "examples", "block18-forensic-retention")
	cases := block18FixtureCases()
	if os.Getenv("SDP_TRACE_UPDATE_BLOCK18_FIXTURES") == "1" {
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
	expected := map[string]block18FixtureCase{}
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
			t.Fatalf("unexpected Block 18 fixture %s", entry.Name())
		}
		payload, err := os.ReadFile(filepath.Join(fixtureDir, entry.Name()))
		if err != nil {
			t.Fatalf("read %s: %v", entry.Name(), err)
		}
		var result AssessmentResult
		if err := json.Unmarshal(payload, &result); err != nil {
			t.Fatalf("unmarshal %s: %v", entry.Name(), err)
		}
		if result.SchemaVersion != SchemaVersion || result.SelectedProfile != ProfileForensicRetention {
			t.Fatalf("%s schema/profile = %s/%s", entry.Name(), result.SchemaVersion, result.SelectedProfile)
		}
		if len(result.ForensicConditions) != len(forensicConditionIDs) {
			t.Fatalf("%s condition count = %d", entry.Name(), len(result.ForensicConditions))
		}
		if result.ForensicRetentionAssessment != topLevel(result.ForensicConditions) {
			t.Fatalf("%s assessment = %s, want %s from conditions", entry.Name(), result.ForensicRetentionAssessment, topLevel(result.ForensicConditions))
		}
		condition := conditionByID(result.ForensicConditions, fixture.conditionID)
		if condition.State != fixture.conditionState || condition.ReasonCode != fixture.reasonCode {
			t.Fatalf("%s condition %s = %s/%s, want %s/%s", entry.Name(), fixture.conditionID, condition.State, condition.ReasonCode, fixture.conditionState, fixture.reasonCode)
		}
		seen++
	}
	if seen != len(expected) {
		t.Fatalf("seen fixtures = %d, want %d", seen, len(expected))
	}
}

type block18FixtureCase struct {
	name           string
	mutate         func(*Input)
	topLevel       string
	conditionID    string
	conditionState string
	reasonCode     string
}

func (fixture block18FixtureCase) result() AssessmentResult {
	input := validInput()
	if fixture.mutate != nil {
		fixture.mutate(&input)
	}
	return Evaluate(input)
}

func block18FixtureCases() []block18FixtureCase {
	return []block18FixtureCase{
		{name: "valid-forensic-retention.assessment-result.json", topLevel: StatePass, conditionID: "raw_reference_bound", conditionState: StatePass, reasonCode: "raw_reference_bound"},
		{name: "digest-only-critical-fail.assessment-result.json", mutate: func(input *Input) {
			input.Run.Events[0].RetentionMode = RetentionModeDigestOnly
			input.Run.Events[0].RawReference = nil
		}, topLevel: StateFail, conditionID: "critical_evidence_reconstructable", conditionState: StateFail, reasonCode: "critical_evidence_digest_only"},
		{name: "external-access-unverifiable-cannot-verify.assessment-result.json", mutate: func(input *Input) {
			input.Run.Events[1].RawReference.AccessState = AccessStateRestricted
			input.Run.Events[1].RawReference.AccessStateLastVerified = ""
		}, topLevel: StateCannotVerify, conditionID: "raw_reference_bound", conditionState: StateCannotVerify, reasonCode: "access_unverifiable"},
		{name: "weak-digest-fail.assessment-result.json", mutate: func(input *Input) {
			input.Run.Events[1].RawReference.Digest.Algorithm = "sha1"
		}, topLevel: StateFail, conditionID: "raw_reference_bound", conditionState: StateFail, reasonCode: "weak_digest"},
		{name: "authority-self-claim-cannot-verify.assessment-result.json", mutate: func(input *Input) {
			input.Run.Events[0].RedactionAuthority.VerificationState = AuthoritySelfClaimed
		}, topLevel: StateCannotVerify, conditionID: "redaction_policy_bound", conditionState: StateCannotVerify, reasonCode: "authority_self_claimed"},
		{name: "withhold-not-assessed-fail.assessment-result.json", mutate: func(input *Input) {
			input.Run.Events[0].RedactionAction = RedactionActionWithhold
			input.Run.Events[0].RedactionRuleRefs = []string{"withhold-privacy-v1"}
			input.Run.Events[0].RetentionMode = RetentionModeNotAssessed
			input.Run.Events[0].Withholding = &Withholding{
				Authority:     AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
				Requestor:     AuthorityRef{ActorID: "human:privacy-owner", VerificationState: AuthorityVerified},
				ReasonCode:    "privacy_request",
				Justification: "privacy legal hold",
			}
			input.Run.Events[0].RawReference = nil
		}, topLevel: StateFail, conditionID: "critical_evidence_reconstructable", conditionState: StateFail, reasonCode: "critical_evidence_not_assessed"},
		{name: "missing-policy-cannot-verify.assessment-result.json", mutate: func(input *Input) {
			input.Policy = Policy{}
		}, topLevel: StateCannotVerify, conditionID: "redaction_policy_bound", conditionState: StateCannotVerify, reasonCode: "missing_redaction_policy"},
	}
}

func validInput() Input {
	policyDigest := "1111111111111111111111111111111111111111111111111111111111111111"
	return Input{
		Policy: Policy{
			PolicyID:              "customer-forensic-policy-v1",
			SchemaVersion:         "1.0.0",
			PolicyDigest:          policyDigest,
			PolicyProvenance:      Provenance{Source: "vcs", Digest: policyDigest},
			AllowedRetentionModes: []string{RetentionModeDigestOnly, RetentionModeSanitizedExcerpt, RetentionModeEncryptedRawRef, RetentionModeExternalArtifactRef, RetentionModeNotAssessed},
			RedactionActions:      []string{RedactionActionApplyRule, RedactionActionWithhold, RedactionActionMarkUnavailable},
			ForbiddenPersistenceClasses: []string{
				"credentials",
				"tokens",
				"raw_prompts",
				"raw_model_responses",
				"source_snippets",
				"stdout_stderr_bodies",
				"oidc_tokens",
				"adapter_secrets",
				"gateway_tokens",
				"checkpoint_key_material",
			},
			CriticalEventFamilies: []string{"command_finished", "test_output_observed"},
			Authority:             AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
			ProfileMappings: []ProfileMapping{
				{
					EventFamily:            "command_finished",
					RequiredRetentionModes: []string{RetentionModeSanitizedExcerpt, RetentionModeEncryptedRawRef, RetentionModeExternalArtifactRef},
					Critical:               true,
					Authority:              AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
				},
				{
					EventFamily:            "test_output_observed",
					RequiredRetentionModes: []string{RetentionModeSanitizedExcerpt, RetentionModeEncryptedRawRef, RetentionModeExternalArtifactRef},
					Critical:               true,
					Authority:              AuthorityRef{ActorID: "human:security-owner", VerificationState: AuthorityVerified},
				},
			},
			UnresolvedRedactionImpact: "fail_forensic_retention",
			Rules: []Rule{
				{RuleID: "secret-token-v1", DetectorFamily: "secret", RuleVersion: "1.0.0", Action: RedactionActionApplyRule, RetentionMode: RetentionModeSanitizedExcerpt},
				{RuleID: "withhold-privacy-v1", DetectorFamily: "privacy", RuleVersion: "1.0.0", Action: RedactionActionWithhold, RetentionMode: RetentionModeNotAssessed},
			},
		},
		Run: RunEvidence{
			RunID:                 "forensic-run-1",
			SelectedProfile:       ProfileForensicRetention,
			RedactionPolicyDigest: policyDigest,
			ProfileSelection: ProfileSelection{
				ActorID:                 "human:security-owner",
				SelectedProfile:         ProfileForensicRetention,
				RedactionPolicyDigest:   policyDigest,
				Justification:           "incident review",
				AuthorityVerified:       true,
				SelectionEvidenceDigest: "2222222222222222222222222222222222222222222222222222222222222222",
			},
			Events: []EventRetention{
				{
					EventType:              "command_finished",
					RetentionMode:          RetentionModeSanitizedExcerpt,
					ForensicImportance:     "critical",
					RedactionPolicyDigest:  policyDigest,
					RedactionInputDigest:   "3333333333333333333333333333333333333333333333333333333333333333",
					RedactedPayloadDigest:  "4444444444444444444444444444444444444444444444444444444444444444",
					RedactionAction:        RedactionActionApplyRule,
					RedactionRuleRefs:      []string{"secret-token-v1"},
					RedactionUnresolved:    false,
					SecretLikeValuePresent: false,
					RedactionAuthority: AuthorityRef{
						ActorID:           "human:security-owner",
						VerificationState: AuthorityVerified,
					},
				},
				{
					EventType:              "test_output_observed",
					RetentionMode:          RetentionModeExternalArtifactRef,
					ForensicImportance:     "critical",
					RedactionPolicyDigest:  policyDigest,
					RedactionInputDigest:   "5555555555555555555555555555555555555555555555555555555555555555",
					RedactedPayloadDigest:  "6666666666666666666666666666666666666666666666666666666666666666",
					RedactionAction:        RedactionActionApplyRule,
					RedactionRuleRefs:      []string{"secret-token-v1"},
					RedactionUnresolved:    false,
					SecretLikeValuePresent: false,
					RedactionAuthority: AuthorityRef{
						ActorID:           "human:security-owner",
						VerificationState: AuthorityVerified,
					},
					RawReference: &RawReference{
						ReferenceType:           RetentionModeExternalArtifactRef,
						ReferenceURI:            "artifact://ci/run-1/test-output",
						Digest:                  Digest{Algorithm: "sha256", Value: "7777777777777777777777777777777777777777777777777777777777777777"},
						AccessState:             AccessStateVerifiedAvailable,
						AccessStateLastVerified: "2026-05-07T10:00:00Z",
						KeyCustodyState:         KeyCustodyNotApplicable,
						RetentionLifecycle:      RetentionLifecycle{State: RetentionLifecycleActive, PolicyRef: "policy:retain-30d", ExpiresAt: "2026-06-06T10:00:00Z"},
						UnavailableReason:       UnavailableReasonNotAssessed,
					},
				},
			},
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
