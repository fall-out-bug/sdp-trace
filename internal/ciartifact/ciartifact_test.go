package ciartifact

import "testing"

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
	if family.ReasonCode != "checked_in_claim_contradicts_ci_artifacts" || family.BindingState != BindingMismatch {
		t.Fatalf("family = %+v", family)
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
			for _, reason := range result.Reasons {
				if reason == "raw-secret-value" {
					t.Fatalf("unsafe reason leaked raw marker: %v", result.Reasons)
				}
			}
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
