package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
)

func TestCIArtifactObservationAssessPassesExplainsAndPreviews(t *testing.T) {
	root := t.TempDir()
	manifestPath := filepath.Join(root, "manifest.json")
	writeTestJSON(t, manifestPath, validCIArtifactManifest())
	outPath := filepath.Join(root, "observation.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"assess",
		"--profile", "ci-artifact-observation",
		"--artifact-manifest", manifestPath,
		"--out", outPath,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("ci artifact assess exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	assertNoCIArtifactLeak(t, out.String())
	var result ciartifact.ObservationResult
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("observation payload: %v", err)
	}
	if result.ArtifactObservationState != ciartifact.StatePass || result.SchemaVersion != ciartifact.SchemaVersion {
		t.Fatalf("observation = %+v", result)
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"assess", "explain", "--assessment-result", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("ci artifact explain exit %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "CI artifact observation: pass") ||
		!strings.Contains(out.String(), "Artifact family provenance: pass") {
		t.Fatalf("ci artifact explain missing fields: %s", out.String())
	}
	assertNoCIArtifactLeak(t, out.String())

	out.Reset()
	errOut.Reset()
	previewOut := filepath.Join(root, "preview.json")
	exit = run([]string{"assess", "preview", "--profile", "ci-artifact-observation", "--artifact-manifest", manifestPath, "--out", previewOut}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("ci artifact preview exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"selected_profile": "ci_artifact_observation"`) ||
		!strings.Contains(out.String(), `"network_fetch": "not_performed"`) {
		t.Fatalf("preview missing fields: %s", out.String())
	}
	if _, err := os.Stat(previewOut); !os.IsNotExist(err) {
		t.Fatalf("ci artifact preview wrote artifact")
	}
}

func TestCIArtifactObservationAssessDistinctExitStates(t *testing.T) {
	cases := []struct {
		name      string
		mutate    func(*ciartifact.Manifest)
		wantExit  int
		wantState string
		wantCode  string
	}{
		{
			name: "checked-in-only",
			mutate: func(m *ciartifact.Manifest) {
				m.RequiredFamilies = []ciartifact.FamilyRequirement{{Family: "provenance", RequiredProducerScope: ciartifact.ProducerCIUploaded}}
				m.ArtifactFamilies = []ciartifact.FamilyInput{{Family: "provenance", ProducerScope: ciartifact.ProducerCheckedIn, ArtifactAccessState: ciartifact.AccessPresent, BindingState: ciartifact.BindingMatched}}
			},
			wantExit: exitCannotVerify, wantState: ciartifact.StateCannotVerify, wantCode: "checked_in_claim_contradicts_ci_artifacts",
		},
		{
			name: "agent-reported",
			mutate: func(m *ciartifact.Manifest) {
				m.ArtifactFamilies[0].ProducerScope = ciartifact.ProducerAgentReported
			},
			wantExit: exitCannotVerify, wantState: ciartifact.StateCannotVerify, wantCode: "agent_reported_claim_without_observed_family",
		},
		{
			name: "index-self-reference",
			mutate: func(m *ciartifact.Manifest) {
				m.ArtifactIndex.State = ciartifact.IndexSelfReference
			},
			wantExit: 1, wantState: ciartifact.StateFail, wantCode: "artifact_index_self_reference",
		},
		{
			name: "unsafe-output",
			mutate: func(m *ciartifact.Manifest) {
				m.OutputSafety.State = ciartifact.StateFail
				m.OutputSafety.UnsafeClasses = []string{"jwt_token"}
			},
			wantExit: 1, wantState: ciartifact.StateFail, wantCode: "unsafe_artifact_output",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			manifest := validCIArtifactManifest()
			tc.mutate(&manifest)
			manifestPath := filepath.Join(root, "manifest.json")
			writeTestJSON(t, manifestPath, manifest)
			var out bytes.Buffer
			var errOut bytes.Buffer
			exit := run([]string{
				"assess",
				"--profile", "ci-artifact-observation",
				"--artifact-manifest", manifestPath,
				"--out", filepath.Join(root, "observation.json"),
			}, &out, &errOut)
			if exit != tc.wantExit {
				t.Fatalf("exit = %d want %d err=%s out=%s", exit, tc.wantExit, errOut.String(), out.String())
			}
			if !strings.Contains(out.String(), `"artifact_observation_state": "`+tc.wantState+`"`) ||
				!strings.Contains(out.String(), `"reason_code": "`+tc.wantCode+`"`) {
				t.Fatalf("output missing state/code: %s", out.String())
			}
			assertNoCIArtifactLeak(t, out.String())
		})
	}
}

func validCIArtifactManifest() ciartifact.Manifest {
	required := []ciartifact.FamilyRequirement{}
	families := []ciartifact.FamilyInput{}
	for _, family := range []string{"run", "report", "witness", "provenance", "evidence", "trace", "artifact_index", "redaction_scan", "review", "change_ci"} {
		required = append(required, ciartifact.FamilyRequirement{Family: family, RequiredProducerScope: ciartifact.ProducerCIUploaded})
		families = append(families, ciartifact.FamilyInput{Family: family, ProducerScope: ciartifact.ProducerCIUploaded, ArtifactAccessState: ciartifact.AccessPresent, BindingState: ciartifact.BindingMatched})
	}
	return ciartifact.Manifest{
		SchemaVersion:    "block26-ci-artifact-observation-input-v1",
		SelectedProfile:  ciartifact.ProfileCIArtifactObservation,
		AuthorityScope:   ciartifact.AuthorityScopeObservation,
		SelectedSource:   ciartifact.SourceIdentity{Repository: "example/repo", Ref: "refs/heads/main", CommitSHA: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
		SelectedRun:      ciartifact.RunIdentity{Provider: "generic-ci", RunID: "run-1", RunAttempt: "1", WorkflowID: "verify", JobID: "job-1"},
		RequiredFamilies: required,
		ArtifactFamilies: families,
		ArtifactIndex:    ciartifact.ArtifactIndexInput{State: ciartifact.IndexValid},
		OutputSafety:     ciartifact.OutputSafetyInput{State: ciartifact.StatePass},
		SafetyRuleset:    ciartifact.SafetyRuleset{ID: ciartifact.SafetyRulesetDefault, SHA256: "1111111111111111111111111111111111111111111111111111111111111111"},
	}
}

func assertNoCIArtifactLeak(t *testing.T, text string) {
	t.Helper()
	for _, marker := range []string{
		"Bearer ",
		"access_token=",
		"oidc_token",
		"BEGIN PRIVATE KEY",
		"raw prompt",
		"raw response",
		"raw_job_log",
		"private_artifact_url",
		"raw-secret-value",
	} {
		if strings.Contains(text, marker) {
			t.Fatalf("output leaked sensitive marker %q: %s", marker, text)
		}
	}
}
