package main

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/authority"
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
	"github.com/fall_out_bug/sdp-trace/internal/releaseproof"
	"github.com/fall_out_bug/sdp-trace/internal/telemetry"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
)

func TestStateExitCodeTables(t *testing.T) {
	tests := []struct {
		name string
		got  int
		want int
	}{
		{name: "harness pass", got: harnessStateExitCode(harnessobs.StatePass), want: 0},
		{name: "harness fail", got: harnessStateExitCode(harnessobs.StateFail), want: 1},
		{name: "harness not assessed", got: harnessStateExitCode(harnessobs.StateNotAssessed), want: exitCannotVerify},
		{name: "harness cannot verify", got: harnessStateExitCode(harnessobs.StateCannotVerify), want: exitCannotVerify},
		{name: "harness unknown", got: harnessStateExitCode("unknown"), want: exitCannotVerify},
		{name: "release pass", got: releaseProofExitCode(releaseproof.StatePass), want: 0},
		{name: "release fail", got: releaseProofExitCode(releaseproof.StateFail), want: 1},
		{name: "release unknown", got: releaseProofExitCode("unknown"), want: exitCannotVerify},
		{name: "authority within", got: authorityExitCode(authority.Result{AuthorityEvaluationState: authority.StateWithinAuthority}), want: 0},
		{name: "authority outside", got: authorityExitCode(authority.Result{AuthorityEvaluationState: authority.StateOutsideAuthority}), want: 1},
		{name: "authority not assessed", got: authorityExitCode(authority.Result{AuthorityEvaluationState: authority.StateNotAssessed}), want: exitCannotVerify},
		{name: "authority unknown", got: authorityExitCode(authority.Result{AuthorityEvaluationState: "unknown"}), want: exitCannotVerify},
		{name: "forensic pass", got: forensicExitCode(forensic.AssessmentResult{ForensicRetentionAssessment: forensic.StatePass}), want: 0},
		{name: "forensic fail", got: forensicExitCode(forensic.AssessmentResult{ForensicRetentionAssessment: forensic.StateFail}), want: 1},
		{name: "forensic unknown", got: forensicExitCode(forensic.AssessmentResult{ForensicRetentionAssessment: "unknown"}), want: exitCannotVerify},
		{name: "ci artifact pass", got: ciArtifactExitCode(ciartifact.ObservationResult{ArtifactObservationState: ciartifact.StatePass}), want: 0},
		{name: "ci artifact fail", got: ciArtifactExitCode(ciartifact.ObservationResult{ArtifactObservationState: ciartifact.StateFail}), want: 1},
		{name: "ci artifact unknown", got: ciArtifactExitCode(ciartifact.ObservationResult{ArtifactObservationState: "unknown"}), want: exitCannotVerify},
		{name: "verifier observed", got: verifierResultExitCode(trace.VerdictObserved), want: 0},
		{name: "verifier not assessed", got: verifierResultExitCode(trace.VerdictNotAssessed), want: 0},
		{name: "verifier fail", got: verifierResultExitCode(trace.VerdictFail), want: 1},
		{name: "verifier cannot verify", got: verifierResultExitCode(trace.VerdictCannotVerify), want: exitCannotVerify},
		{name: "verifier unknown", got: verifierResultExitCode("unknown"), want: 0},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Fatalf("exit code = %d want %d", test.got, test.want)
			}
		})
	}
}

func TestPreviewActionHelpers(t *testing.T) {
	tests := []struct {
		name string
		run  func(map[string]string) []string
		key  string
		want map[string][]string
	}{
		{
			name: "protected",
			run:  protectedPreviewActions,
			key:  "checkpoint",
			want: map[string][]string{
				"absent":             {"Supply checkpoint input before running protected gate."},
				"present_unreadable": {"Replace checkpoint input with readable JSON."},
				"present_malformed":  {"Replace checkpoint input with readable JSON."},
				"present_valid":      {},
			},
		},
		{
			name: "adapter capture",
			run:  adapterCapturePreviewActions,
			key:  "run",
			want: map[string][]string{
				"absent":              {"Supply run before adapter capture assessment."},
				"present_unreadable":  {"Fix run so it is a readable JSON run directory."},
				"present_malformed":   {"Fix run so it is a readable JSON run directory."},
				"present_valid":       nil,
				"present_json_object": nil,
			},
		},
		{
			name: "ci artifact",
			run:  ciArtifactPreviewActions,
			key:  "artifact_manifest",
			want: map[string][]string{
				"absent":             {"Supply artifact manifest before CI artifact observation assessment."},
				"present_unreadable": {"Fix artifact manifest so it is readable JSON."},
				"present_malformed":  {"Fix artifact manifest so it is readable JSON."},
				"present_valid":      nil,
			},
		},
		{
			name: "authority",
			run:  authorityPreviewActions,
			key:  "authority_package",
			want: map[string][]string{
				"absent":             {"Supply authority package before authority envelope assessment."},
				"present_unreadable": {"Fix authority package so it is readable JSON."},
				"present_malformed":  {"Fix authority package so it is readable JSON."},
				"present_valid":      nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for state, want := range test.want {
				got := test.run(map[string]string{test.key: state})
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("%s actions = %#v want %#v", state, got, want)
				}
			}
		})
	}
}

func TestInteractionCommandFailurePaths(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	if code := runInteractionImportTranscript(nil, &out, &errOut); code != exitUsage {
		t.Fatalf("import transcript missing args code = %d", code)
	}
	if code := runInteractionSummarize(nil, &out, &errOut); code != exitUsage {
		t.Fatalf("summarize missing args code = %d", code)
	}
	opts, code, ok := parseInteractionRelayArgs([]string{
		"--task-id", "task-1",
		"--out", filepath.Join(t.TempDir(), "trace.json"),
		"--", "echo", "ok",
	}, &errOut)
	if !ok || code != 0 {
		t.Fatalf("relay args code=%d ok=%v err=%s", code, ok, errOut.String())
	}
	if got := opts.rest(); !reflect.DeepEqual(got, []string{"echo", "ok"}) {
		t.Fatalf("relay rest = %#v", got)
	}
	if _, code, ok := parseInteractionRelayArgs([]string{"--task-id", "task-1", "--out", "trace.json"}, &errOut); ok || code != exitUsage {
		t.Fatalf("relay missing rest code=%d ok=%v", code, ok)
	}
}

func TestJSONOutputAndAssessmentHelpers(t *testing.T) {
	tempDir := t.TempDir()
	if err := writeOptionalJSONFile("", map[string]string{"ignored": "true"}); err != nil {
		t.Fatalf("empty optional json path: %v", err)
	}
	outPath := filepath.Join(tempDir, "summary.json")
	if err := writeOptionalJSONFile(outPath, map[string]string{"state": "pass"}); err != nil {
		t.Fatalf("write optional json: %v", err)
	}
	if data, err := os.ReadFile(outPath); err != nil || !strings.Contains(string(data), `"state": "pass"`) {
		t.Fatalf("optional json data=%q err=%v", data, err)
	}

	var out bytes.Buffer
	explainMissingAuditEvidence([]string{"review"}, &out)
	explainOverrideRequests([]demo.OverrideRequest{{OverrideID: "override-1", State: demo.GateCannotVerify}}, &out)
	if got := out.String(); !strings.Contains(got, "Missing audit evidence: review") || !strings.Contains(got, "Override override-1") {
		t.Fatalf("explanation output = %q", got)
	}

	if !optionalStringMatches("", "actual") || !optionalStringMatches("same", "same") || optionalStringMatches("expected", "actual") {
		t.Fatalf("optionalStringMatches returned unexpected values")
	}
}

func TestManagedAndRepoObserverCommandFailurePaths(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer

	emptyManaged := &flagSet{data: map[string]string{}}
	if code := runManagedAssess(emptyManaged, &out, &errOut); code != exitUsage {
		t.Fatalf("managed missing inputs code = %d", code)
	}
	invalidManaged := &flagSet{data: map[string]string{
		"out":              filepath.Join(t.TempDir(), "managed.json"),
		"run":              "missing-run.json",
		"adapter-registry": "missing-adapters.json",
		"managed-policy":   "missing-policy.json",
		"managed-witness":  "missing-witness.json",
	}}
	if code := runManagedAssess(invalidManaged, &out, &errOut); code != exitUsage {
		t.Fatalf("managed invalid inputs code = %d", code)
	}

	if code := runDoctor(context.Background(), []string{"--help"}, &out, &errOut); code != 0 {
		t.Fatalf("doctor help code = %d", code)
	}
	if code := runDoctor(context.Background(), []string{"unexpected"}, &out, &errOut); code != exitUsage {
		t.Fatalf("doctor rest code = %d", code)
	}
	if code := runInstall(context.Background(), nil, &out, &errOut); code != exitUsage {
		t.Fatalf("install missing subcommand code = %d", code)
	}
}

func TestGateAndFixtureHelpers(t *testing.T) {
	tests := []struct {
		name   string
		states []string
		want   int
	}{
		{name: "pass", states: []string{demo.GatePass}, want: 0},
		{name: "fail", states: []string{demo.GateFail}, want: 1},
		{name: "missing telemetry", states: []string{demo.GateMissingTelemetry}, want: 1},
		{name: "cannot verify", states: []string{demo.GateCannotVerify}, want: exitCannotVerify},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := gateStateExitCode(test.states); got != test.want {
				t.Fatalf("gateStateExitCode = %d want %d", got, test.want)
			}
		})
	}

	if !unexpectedFixtureResultFailed(trace.VerifierResult{Result: trace.VerdictFail}) {
		t.Fatalf("fail verdict should be unexpected failure")
	}
	if !unexpectedFixtureResultFailed(trace.VerifierResult{Result: trace.VerdictCannotVerify}) {
		t.Fatalf("cannot_verify verdict should be unexpected failure")
	}
	if unexpectedFixtureResultFailed(trace.VerifierResult{Result: trace.VerdictObserved}) {
		t.Fatalf("observed verdict should not be unexpected failure")
	}
	if expectedFixtureResultFailed("run", trace.VerifierResult{Result: trace.VerdictFail}, fixtureExpectation{ExpectedResult: string(trace.VerdictFail)}, io.Discard) {
		t.Fatalf("matching fixture expectation failed")
	}
	if !expectedFixtureResultFailed("run", trace.VerifierResult{Result: trace.VerdictFail}, fixtureExpectation{ExpectedResult: string(trace.VerdictObserved)}, io.Discard) {
		t.Fatalf("mismatched fixture expectation did not fail")
	}
}

func TestTelemetryWitnessAndFlagHelpers(t *testing.T) {
	opts := &flagSet{data: map[string]string{
		"profile":            telemetry.ProfilePrometheusTextV1,
		"cross-repo-posture": "posture.json",
		"out":                "metrics.prom",
	}}
	if err := requireTelemetryExportInputs(opts); err != nil {
		t.Fatalf("telemetry inputs: %v", err)
	}
	for _, key := range []string{"profile", "cross-repo-posture", "out"} {
		bad := &flagSet{data: map[string]string{
			"profile":            telemetry.ProfilePrometheusTextV1,
			"cross-repo-posture": "posture.json",
			"out":                "metrics.prom",
		}}
		bad.data[key] = ""
		if err := requireTelemetryExportInputs(bad); err == nil {
			t.Fatalf("missing %s did not fail", key)
		}
	}

	if _, err := buildWitnessRecord(witnessOptions{kind: "unsupported"}); err == nil {
		t.Fatalf("unsupported witness kind did not fail")
	}

	emptyFlags := &flagSet{}
	if got := emptyFlags.stringValue("missing"); got != "" {
		t.Fatalf("empty string flag = %q", got)
	}
	if got := emptyFlags.boolValue("missing"); got {
		t.Fatalf("empty bool flag = true")
	}
	flags := &flagSet{}
	flags.setString("name", "default")
	flags.setBool("enabled", false)
	if err := flags.parse([]string{"--name", "value", "--enabled", "true", "--", "rest"}); err != nil {
		t.Fatalf("parse flags: %v", err)
	}
	if flags.stringValue("name") != "value" || !flags.boolValue("enabled") || !reflect.DeepEqual(flags.rest(), []string{"rest"}) {
		t.Fatalf("parsed flags name=%q enabled=%v rest=%#v", flags.stringValue("name"), flags.boolValue("enabled"), flags.rest())
	}

	for _, args := range [][]string{{"--help"}, {"-h"}, {"help"}} {
		if !isHelp(args) {
			t.Fatalf("isHelp(%#v) = false", args)
		}
	}
	if isHelp(nil) || isHelp([]string{"help", "extra"}) {
		t.Fatalf("isHelp accepted invalid args")
	}
	for _, value := range []string{"true", "false", "1", "0", "TRUE"} {
		if !isBoolLiteral(value) {
			t.Fatalf("isBoolLiteral(%q) = false", value)
		}
	}
	if isBoolLiteral("yes") {
		t.Fatalf("isBoolLiteral accepted yes")
	}
}

func TestRunOptionalSubcommand(t *testing.T) {
	handlers := map[string]subcommandHandler{
		"known": func(args []string, stdout, stderr io.Writer) int {
			if !reflect.DeepEqual(args, []string{"flag"}) {
				t.Fatalf("handler args = %#v", args)
			}
			return 17
		},
	}
	var out bytes.Buffer
	var errOut bytes.Buffer

	if code, ok := runOptionalSubcommand(nil, &out, &errOut, handlers); ok || code != 0 {
		t.Fatalf("empty optional subcommand code=%d ok=%v", code, ok)
	}
	if code, ok := runOptionalSubcommand([]string{"missing"}, &out, &errOut, handlers); ok || code != 0 {
		t.Fatalf("missing optional subcommand code=%d ok=%v", code, ok)
	}
	if code, ok := runOptionalSubcommand([]string{"known", "flag"}, &out, &errOut, handlers); !ok || code != 17 {
		t.Fatalf("known optional subcommand code=%d ok=%v", code, ok)
	}
}
