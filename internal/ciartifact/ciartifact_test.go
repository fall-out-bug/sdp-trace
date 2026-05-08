package ciartifact

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEvaluateValidCIUploadedBundlePasses(t *testing.T) {
	result := Evaluate(validManifest())
	if result.ArtifactObservationState != StatePass {
		t.Fatalf("state = %s reasons=%v", result.ArtifactObservationState, result.Reasons)
	}
	if result.ProducerScope != ProducerCIUploaded || result.ArtifactAccessState != AccessPresent {
		t.Fatalf("aggregate producer/access = %s/%s", result.ProducerScope, result.ArtifactAccessState)
	}
}

func TestEvaluateCheckedInOnlyCannotVerify(t *testing.T) {
	manifest := validManifest()
	manifest.ArtifactFamilies = []FamilyInput{{Family: "provenance", ProducerScope: ProducerCheckedIn, ArtifactAccessState: AccessPresent, BindingState: BindingMatched}}
	manifest.RequiredFamilies = []FamilyRequirement{{Family: "provenance", RequiredProducerScope: ProducerCIUploaded}}
	result := Evaluate(manifest)
	if result.ArtifactObservationState != StateCannotVerify {
		t.Fatalf("state = %s", result.ArtifactObservationState)
	}
	family := result.ArtifactFamilies[0]
	if family.ReasonCode != "checked_in_claim_contradicts_ci_artifacts" || family.BindingState != BindingMatched {
		t.Fatalf("family = %+v", family)
	}
	if result.Bindings.SourceBindingState != BindingMatched || result.Bindings.ProducerBindingState != BindingMismatch {
		t.Fatalf("bindings = %+v", result.Bindings)
	}
}

func TestEvaluateContradictionsFail(t *testing.T) {
	for name, mutate := range map[string]func(*Manifest){
		"artifact-index-self-reference": func(m *Manifest) { m.ArtifactIndex.State = IndexSelfReference },
		"artifact-digest-mismatch":      func(m *Manifest) { m.ArtifactIndex.State = IndexDigestMismatch },
		"source-run-binding-mismatch":   func(m *Manifest) { m.ArtifactFamilies[0].BindingState = BindingMismatch },
		"unsafe-output": func(m *Manifest) {
			m.OutputSafety.State = StateFail
			m.OutputSafety.UnsafeClasses = []string{"jwt_token", "raw-secret-value"}
		},
	} {
		t.Run(name, func(t *testing.T) {
			manifest := validManifest()
			mutate(&manifest)
			result := Evaluate(manifest)
			if result.ArtifactObservationState != StateFail {
				t.Fatalf("state = %s reasons=%v", result.ArtifactObservationState, result.Reasons)
			}
			assertNoRawMarker(t, result)
		})
	}
}

func TestEvaluateCannotVerifyCases(t *testing.T) {
	for name, mutate := range map[string]func(*Manifest){
		"ci-bundle-absent":           func(m *Manifest) { m.ArtifactFamilies = nil },
		"ci-bundle-partial":          func(m *Manifest) { m.ArtifactFamilies[0].ArtifactAccessState = AccessPartial },
		"source-run-binding-missing": func(m *Manifest) { m.ArtifactFamilies[0].BindingState = BindingAbsent },
		"artifact-expired":           func(m *Manifest) { m.ArtifactFamilies[0].ArtifactAccessState = AccessExpired },
		"external-artifact-ref-unverifiable": func(m *Manifest) {
			m.ArtifactFamilies[0].ProducerScope = ProducerExternalArtifactRef
		},
		"agent-reported-happy-path": func(m *Manifest) {
			m.ArtifactFamilies[0].ProducerScope = ProducerAgentReported
		},
	} {
		t.Run(name, func(t *testing.T) {
			manifest := validManifest()
			mutate(&manifest)
			result := Evaluate(manifest)
			if result.ArtifactObservationState != StateCannotVerify {
				t.Fatalf("state = %s reasons=%v", result.ArtifactObservationState, result.Reasons)
			}
		})
	}
}

func TestEvaluateTopLevelNotAssessedAndMixedAggregates(t *testing.T) {
	manifest := validManifest()
	manifest.RequiredFamilies = nil
	manifest.ArtifactFamilies = nil
	result := Evaluate(manifest)
	if result.ArtifactObservationState != StateNotAssessed {
		t.Fatalf("state = %s", result.ArtifactObservationState)
	}
	payload, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	for _, field := range []string{`"required_families":[]`, `"artifact_families":[]`} {
		if !strings.Contains(string(payload), field) {
			t.Fatalf("empty array field %s missing from %s", field, string(payload))
		}
	}

	manifest = validManifest()
	manifest.ArtifactFamilies[0].ProducerScope = ProducerCheckedIn
	result = Evaluate(manifest)
	if result.ProducerScope != "mixed" {
		t.Fatalf("producer scope = %s", result.ProducerScope)
	}
}

func TestEvaluateUnknownStatesBecomeCannotVerify(t *testing.T) {
	manifest := validManifest()
	manifest.ArtifactIndex.State = "new-index-state"
	manifest.OutputSafety.State = "new-safety-state"
	result := Evaluate(manifest)
	if result.ArtifactObservationState != StateCannotVerify {
		t.Fatalf("state = %s reasons=%v", result.ArtifactObservationState, result.Reasons)
	}
	if result.ArtifactIndex.Result != StateCannotVerify || result.OutputSafety.State != StateCannotVerify {
		t.Fatalf("index/safety = %+v %+v", result.ArtifactIndex, result.OutputSafety)
	}
}

func TestEvaluateManifestMutationRegressions(t *testing.T) {
	for name, tc := range map[string]struct {
		mutate        func(*Manifest)
		wantState     string
		wantReason    string
		wantSourceRun string
		wantProducer  string
	}{
		"checked-in-producer-does-not-poison-source-run-binding": {
			mutate: func(m *Manifest) {
				m.RequiredFamilies = []FamilyRequirement{{Family: "provenance", RequiredProducerScope: ProducerCIUploaded}}
				m.ArtifactFamilies = []FamilyInput{{Family: "provenance", ProducerScope: ProducerCheckedIn, ArtifactAccessState: AccessPresent, BindingState: BindingMatched}}
			},
			wantState:     StateCannotVerify,
			wantReason:    "checked_in_claim_contradicts_ci_artifacts",
			wantSourceRun: BindingMatched,
			wantProducer:  BindingMismatch,
		},
		"unknown-access-state-is-cannot-verify-not-pass": {
			mutate: func(m *Manifest) {
				m.ArtifactFamilies[0].ArtifactAccessState = "new-access-state"
			},
			wantState:  StateCannotVerify,
			wantReason: "artifact_access_cannot_verify",
		},
		"unsafe-required-producer-is-normalized-before-output": {
			mutate: func(m *Manifest) {
				m.RequiredFamilies[0].RequiredProducerScope = "Bearer raw-secret-value"
			},
			wantState:     StatePass,
			wantSourceRun: BindingMatched,
			wantProducer:  BindingMatched,
		},
	} {
		t.Run(name, func(t *testing.T) {
			manifest := validManifest()
			tc.mutate(&manifest)
			result := Evaluate(manifest)
			if result.ArtifactObservationState != tc.wantState {
				t.Fatalf("state = %s want %s reasons=%v", result.ArtifactObservationState, tc.wantState, result.Reasons)
			}
			if tc.wantReason != "" && !hasReasonCode(result, tc.wantReason) {
				t.Fatalf("missing reason code %s in %+v", tc.wantReason, result)
			}
			if tc.wantSourceRun != "" && result.Bindings.SourceBindingState != tc.wantSourceRun {
				t.Fatalf("source binding = %s want %s", result.Bindings.SourceBindingState, tc.wantSourceRun)
			}
			if tc.wantProducer != "" && result.Bindings.ProducerBindingState != tc.wantProducer {
				t.Fatalf("producer binding = %s want %s", result.Bindings.ProducerBindingState, tc.wantProducer)
			}
			assertNoRawMarker(t, result)
		})
	}
}

func TestEvaluateDefaultsIndexAndSafetyToNotAssessed(t *testing.T) {
	manifest := validManifest()
	manifest.ArtifactIndex = ArtifactIndexInput{}
	manifest.OutputSafety = OutputSafetyInput{}
	result := Evaluate(manifest)
	if result.ArtifactObservationState != StatePass {
		t.Fatalf("state = %s reasons=%v", result.ArtifactObservationState, result.Reasons)
	}
	if result.ArtifactIndex.State != IndexNotAssessed || result.ArtifactIndex.Result != StateNotAssessed {
		t.Fatalf("artifact index = %+v", result.ArtifactIndex)
	}
	if result.OutputSafety.State != StateNotAssessed {
		t.Fatalf("output safety = %+v", result.OutputSafety)
	}
}

func TestEvaluateSafeTokenLengthMatchesSchema(t *testing.T) {
	longSafeToken := strings.Repeat("a", 129)
	manifest := validManifest()
	manifest.AuthorityScope = longSafeToken
	manifest.SafetyRuleset.ID = longSafeToken
	result := Evaluate(manifest)
	if result.AuthorityScope != AuthorityScopeObservation {
		t.Fatalf("authority scope = %q", result.AuthorityScope)
	}
	if result.SafetyRuleset.ID != SafetyRulesetDefault {
		t.Fatalf("safety ruleset id = %q", result.SafetyRuleset.ID)
	}
	assertNoRawMarker(t, result)

	manifest = validManifest()
	manifest.AuthorityScope = ""
	manifest.SafetyRuleset.ID = ""
	result = Evaluate(manifest)
	if result.AuthorityScope != AuthorityScopeObservation {
		t.Fatalf("empty authority scope = %q", result.AuthorityScope)
	}
	if result.SafetyRuleset.ID != SafetyRulesetDefault {
		t.Fatalf("empty safety ruleset id = %q", result.SafetyRuleset.ID)
	}
}

func TestFixtureMatrixScenarios(t *testing.T) {
	matrixPath := filepath.Join("..", "..", "examples", "block26-ci-artifact-observation", "fixture-matrix.json")
	data, err := os.ReadFile(matrixPath)
	if err != nil {
		t.Fatalf("read fixture matrix: %v", err)
	}
	var matrix struct {
		Scenarios []struct {
			Name          string `json:"name"`
			ExpectedState string `json:"expected_state"`
			ReasonCode    string `json:"reason_code"`
		} `json:"scenarios"`
	}
	if err := json.Unmarshal(data, &matrix); err != nil {
		t.Fatalf("decode fixture matrix: %v", err)
	}
	for _, scenario := range matrix.Scenarios {
		t.Run(scenario.Name, func(t *testing.T) {
			manifest := manifestForScenario(scenario.Name)
			result := Evaluate(manifest)
			if result.ArtifactObservationState != scenario.ExpectedState {
				t.Fatalf("state = %s want %s reasons=%v", result.ArtifactObservationState, scenario.ExpectedState, result.Reasons)
			}
			if scenario.ExpectedState != StatePass && !hasReasonCode(result, scenario.ReasonCode) {
				t.Fatalf("missing reason code %s in %+v", scenario.ReasonCode, result)
			}
		})
	}
}

func TestOutputSafetyDoesNotEchoForbiddenMarkers(t *testing.T) {
	markers := []string{
		"Bearer secret-token",
		"access_token=secret-token",
		"oidc_token",
		"BEGIN PRIVATE KEY",
		"raw prompt",
		"raw response",
		"raw_job_log",
		"private_artifact_url",
		"/Users/private/path",
		"eyJhbGciOi",
	}
	for _, marker := range markers {
		t.Run(marker, func(t *testing.T) {
			manifest := validManifest()
			manifest.SelectedSource.Repository = marker
			manifest.OutputSafety.State = StateFail
			manifest.OutputSafety.UnsafeClasses = []string{"jwt_token", "private_key", "raw_job_log"}
			result := Evaluate(manifest)
			assertNoRawMarker(t, result)
		})
	}
}

func TestEvaluateSanitizesUnsafeIdentityAndEnums(t *testing.T) {
	manifest := validManifest()
	manifest.AuthorityScope = "Bearer raw-secret-value"
	manifest.SafetyRuleset.ID = "access_token=raw-secret-value"
	manifest.SafetyRuleset.SHA256 = "raw-secret-value"
	manifest.RequiredFamilies[0].RequiredProducerScope = "Bearer raw-secret-value"
	manifest.SelectedSource.Repository = "example/repo?token=raw-secret-value"
	manifest.SelectedSource.Ref = "refs/heads/main raw-secret-value"
	manifest.SelectedRun.RunID = "run-1?access_token=raw-secret-value"
	manifest.ArtifactFamilies = append(manifest.ArtifactFamilies, FamilyInput{
		Family:              "raw-secret-value",
		ProducerScope:       "Bearer raw-secret-value",
		ArtifactAccessState: "raw-secret-value",
		BindingState:        "raw-secret-value",
	})
	result := Evaluate(manifest)
	if result.ArtifactObservationState != StateCannotVerify {
		t.Fatalf("state = %s reasons=%v", result.ArtifactObservationState, result.Reasons)
	}
	if result.SelectedSource.Repository != "" || result.SelectedSource.Ref != "" || result.SelectedRun.RunID != "" {
		t.Fatalf("unsafe identity copied into result: %+v %+v", result.SelectedSource, result.SelectedRun)
	}
	if result.AuthorityScope != AuthorityScopeObservation {
		t.Fatalf("unsafe authority scope copied into result: %q", result.AuthorityScope)
	}
	if result.SafetyRuleset.ID != SafetyRulesetDefault || result.SafetyRuleset.SHA256 == "" || strings.Contains(result.SafetyRuleset.SHA256, "raw-secret-value") {
		t.Fatalf("unsafe safety ruleset copied into result: %+v", result.SafetyRuleset)
	}
	if result.RequiredFamilies[0].RequiredProducerScope != ProducerCIUploaded {
		t.Fatalf("unsafe required producer copied into result: %+v", result.RequiredFamilies[0])
	}
	for _, family := range result.ArtifactFamilies {
		if family.Family == "raw-secret-value" || family.ProducerScope == "Bearer raw-secret-value" || family.RequiredProducer == "Bearer raw-secret-value" {
			t.Fatalf("unsafe enum copied into family: %+v", family)
		}
	}
	for _, reason := range result.Reasons {
		if reason == "raw-secret-value" {
			t.Fatalf("unsafe reason copied: %v", result.Reasons)
		}
	}
}

func manifestForScenario(name string) Manifest {
	manifest := validManifest()
	switch name {
	case "checked-in-only-claim":
		manifest.RequiredFamilies = []FamilyRequirement{{Family: "provenance", RequiredProducerScope: ProducerCIUploaded}}
		manifest.ArtifactFamilies = []FamilyInput{{Family: "provenance", ProducerScope: ProducerCheckedIn, ArtifactAccessState: AccessPresent, BindingState: BindingMatched}}
	case "ci-bundle-absent":
		manifest.ArtifactFamilies = nil
	case "ci-bundle-partial":
		manifest.ArtifactFamilies[0].ArtifactAccessState = AccessPartial
	case "artifact-index-self-reference":
		manifest.ArtifactIndex.State = IndexSelfReference
	case "artifact-digest-mismatch":
		manifest.ArtifactIndex.State = IndexDigestMismatch
	case "source-run-binding-missing":
		manifest.ArtifactFamilies[0].BindingState = BindingAbsent
	case "source-run-binding-mismatch":
		manifest.ArtifactFamilies[0].BindingState = BindingMismatch
	case "agent-reported-happy-path":
		manifest.ArtifactFamilies[0].ProducerScope = ProducerAgentReported
	case "unsafe-artifact-output":
		manifest.OutputSafety.State = StateFail
		manifest.OutputSafety.UnsafeClasses = []string{"jwt_token"}
	case "artifact-expired":
		manifest.ArtifactFamilies[0].ArtifactAccessState = AccessExpired
	case "external-artifact-ref-unverifiable":
		manifest.ArtifactFamilies[0].ProducerScope = ProducerExternalArtifactRef
	}
	return manifest
}

func hasReasonCode(result ObservationResult, code string) bool {
	for _, family := range result.ArtifactFamilies {
		if family.ReasonCode == code {
			return true
		}
	}
	if result.ArtifactIndex.ReasonCode == code || result.OutputSafety.ReasonCode == code {
		return true
	}
	for _, reason := range result.Reasons {
		if strings.HasPrefix(reason, code+":") {
			return true
		}
	}
	return false
}

func assertNoRawMarker(t *testing.T, result ObservationResult) {
	t.Helper()
	payload, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal result: %v", err)
	}
	for _, marker := range []string{
		"raw-secret-value",
		"Bearer secret-token",
		"access_token=secret-token",
		"oidc_token",
		"BEGIN PRIVATE KEY",
		"raw prompt",
		"raw response",
		"private_artifact_url",
		"/Users/private/path",
	} {
		if strings.Contains(string(payload), marker) {
			t.Fatalf("output leaked raw marker %q: %s", marker, string(payload))
		}
	}
}

func validManifest() Manifest {
	required := []FamilyRequirement{}
	families := []FamilyInput{}
	for _, family := range familyOrder {
		required = append(required, FamilyRequirement{Family: family, RequiredProducerScope: ProducerCIUploaded})
		families = append(families, FamilyInput{Family: family, ProducerScope: ProducerCIUploaded, ArtifactAccessState: AccessPresent, BindingState: BindingMatched})
	}
	return Manifest{
		SchemaVersion:    "block26-ci-artifact-observation-input-v1",
		SelectedProfile:  ProfileCIArtifactObservation,
		AuthorityScope:   AuthorityScopeObservation,
		SelectedSource:   SourceIdentity{Repository: "example/repo", Ref: "refs/heads/main", CommitSHA: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		SelectedRun:      RunIdentity{Provider: "generic-ci", RunID: "run-1", RunAttempt: "1", WorkflowID: "verify", JobID: "job-1"},
		RequiredFamilies: required,
		ArtifactFamilies: families,
		ArtifactIndex:    ArtifactIndexInput{State: IndexValid},
		OutputSafety:     OutputSafetyInput{State: StatePass},
		SafetyRuleset:    SafetyRuleset{ID: SafetyRulesetDefault, SHA256: "1111111111111111111111111111111111111111111111111111111111111111"},
	}
}
