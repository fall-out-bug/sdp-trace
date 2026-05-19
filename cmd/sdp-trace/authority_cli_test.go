package main

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/authority"
)

func TestAuthorityEnvelopeAssessExplainsAndPreviews(t *testing.T) {
	root := t.TempDir()
	packagePath := filepath.Join(root, "authority-package.json")
	writeTestJSON(t, packagePath, validAuthorityPackage())
	outPath := filepath.Join(root, "authority-result.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"assess",
		"--profile", "authority-envelope",
		"--authority-package", packagePath,
		"--out", outPath,
	}, &out, &errOut)
	if exit != 1 {
		t.Fatalf("authority assess exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	assertNoAuthorityLeak(t, out.String())
	var result authority.Result
	if err := json.Unmarshal(out.Bytes(), &result); err != nil {
		t.Fatalf("authority payload: %v", err)
	}
	if result.AuthorityEvaluationState != authority.StateOutsideAuthority || result.SchemaVersion != authority.ResultSchemaVersion {
		t.Fatalf("authority result = %+v", result)
	}

	out.Reset()
	errOut.Reset()
	exit = run([]string{"assess", "explain", "--assessment-result", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("authority explain exit %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "Authority evaluation: outside_authority") ||
		!strings.Contains(out.String(), "Observed action event-mutation-1: outside_authority") {
		t.Fatalf("authority explain missing fields: %s", out.String())
	}
	assertNoAuthorityLeak(t, out.String())

	out.Reset()
	errOut.Reset()
	exit = run([]string{"assess", "preview", "--profile", "authority-envelope", "--authority-package", packagePath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("authority preview exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"selected_profile": "authority-envelope"`) ||
		!strings.Contains(out.String(), `"policy_effects": "not_emitted"`) {
		t.Fatalf("authority preview missing fields: %s", out.String())
	}
	assertNoAuthorityLeak(t, out.String())
}

func TestAuthorityEnvelopeAssessDistinctExitStates(t *testing.T) {
	cases := []assessCLICase[authority.Package]{
		{
			name: "within-authority",
			mutate: func(pkg *authority.Package) {
				pkg.AuthorityEnvelopes[0].TargetRules[0].AllowedEvents = []string{"direct_mutation"}
				pkg.AuthorityEnvelopes[0].TargetRules[0].DeniedEvents = nil
			},
			wantExit: 0, wantState: authority.StateWithinAuthority, wantCode: "target_event_allowed",
		},
		{
			name: "not-assessed",
			mutate: func(pkg *authority.Package) {
				pkg.SelectedPolicyID = ""
			},
			wantExit: exitCannotVerify, wantState: authority.StateNotAssessed, wantCode: "policy_not_selected",
		},
		{
			name: "cannot-verify",
			mutate: func(pkg *authority.Package) {
				pkg.AuthorityEnvelopes[0].AllowedEvents = []string{"direct_mutation"}
				pkg.AuthorityEnvelopes[0].DeniedEvents = []string{"direct_mutation"}
			},
			wantExit: exitCannotVerify, wantState: authority.StateCannotVerify, wantCode: "allow_deny_event_conflict",
		},
	}
	runAssessCLICases(t, cases, validAuthorityPackage, "authority-package.json", "authority_evaluation_state", func(root, packagePath string) []string {
		return []string{
			"assess",
			"--profile", "authority-envelope",
			"--authority-package", packagePath,
			"--out", filepath.Join(root, "authority-result.json"),
		}
	}, assertNoAuthorityLeak)
}

func validAuthorityPackage() authority.Package {
	return authority.Package{
		SchemaVersion:    authority.PackageSchemaVersion,
		SelectedPolicyID: "policy-deny-ci",
		Actors: []authority.ActorDeclaration{{
			ActorID:      "agent-1",
			ActorType:    "ai_agent",
			DeclaredRole: "observer",
			Harness:      "generic-harness",
			OperationID:  "op-1",
		}},
		AuthorityEnvelopes: []authority.AuthorityEnvelope{{
			SchemaVersion:  "authority-envelope-v1",
			TaskID:         "task-1",
			PolicyID:       "policy-deny-ci",
			AuthorityScope: "task",
			ActorRef:       "agent-1",
			AllowedEvents:  []string{"review", "feedback"},
			DeniedEvents:   []string{},
			TargetRules: []authority.TargetRule{{
				RuleID:        "rule-ci-denied",
				TargetPattern: ".github/workflows/**",
				DeniedEvents:  []string{"direct_mutation"},
			}},
		}},
		ObservedActions: []authority.ObservedAction{{
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

func assertNoAuthorityLeak(t *testing.T, text string) {
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
