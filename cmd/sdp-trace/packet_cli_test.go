package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/fall_out_bug/sdp-trace/internal/packet"
)

func TestPacketHandlersExposeExpectedSubcommands(t *testing.T) {
	want := map[string]subcommandHandler{
		"build-pr":     runPacketBuildPR,
		"build-github": runPacketBuildGitHub,
		"validate":     runPacketValidate,
		"check-demo":   runPacketCheckDemo,
		"render":       runPacketRender,
	}
	if len(packetHandlers) != len(want) {
		t.Fatalf("packetHandlers length = %d, want %d", len(packetHandlers), len(want))
	}
	for name, wantHandler := range want {
		gotHandler, ok := packetHandlers[name]
		if !ok {
			t.Fatalf("packetHandlers missing %s", name)
		}
		if functionName(gotHandler) != functionName(wantHandler) {
			t.Fatalf("packetHandlers[%s] = %s, want %s", name, functionName(gotHandler), functionName(wantHandler))
		}
	}
}

func TestPacketRequiredFlagsKeepNamesAndDiagnostics(t *testing.T) {
	cases := []struct {
		name string
		got  []requiredCLIFlag
		want []requiredCLIFlag
	}{
		{
			name: "build-pr",
			got:  packetBuildPRRequiredFlags,
			want: []requiredCLIFlag{{"out", "packet build-pr requires --out"}},
		},
		{
			name: "build-github",
			got:  packetBuildGitHubRequiredFlags,
			want: []requiredCLIFlag{
				{"github-input", "packet build-github requires --github-input"},
				{"out", "packet build-github requires --out"},
			},
		},
		{
			name: "validate",
			got:  packetValidateRequiredFlags,
			want: []requiredCLIFlag{{"bundle", "packet validate requires --bundle"}},
		},
		{
			name: "check-demo",
			got:  packetCheckDemoRequiredFlags,
			want: []requiredCLIFlag{{"bundle", "packet check-demo requires --bundle"}},
		},
		{
			name: "render",
			got:  packetRenderRequiredFlags,
			want: []requiredCLIFlag{
				{"bundle", "packet render requires --bundle"},
				{"out", "packet render requires --out"},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.got, tc.want) {
				t.Fatalf("required flags = %#v, want %#v", tc.got, tc.want)
			}
		})
	}
}

func functionName(fn subcommandHandler) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}

func TestPacketValidateAndRenderCLI(t *testing.T) {
	root := t.TempDir()
	bundlePath := filepath.Join(root, "bundle.json")
	writeTestJSON(t, bundlePath, validPacketBundleForCLI())
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "validate", "--bundle", bundlePath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("validate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"state": "pass"`) {
		t.Fatalf("validate output missing pass: %s", out.String())
	}

	out.Reset()
	errOut.Reset()
	packetPath := filepath.Join(root, "change-evidence-packet.md")
	exit = run([]string{"packet", "render", "--bundle", bundlePath, "--out", packetPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("render exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	raw, err := os.ReadFile(packetPath)
	if err != nil {
		t.Fatalf("read rendered packet: %v", err)
	}
	if !strings.Contains(string(raw), "## Required Rows") ||
		!strings.Contains(string(raw), "| PC-THEATER | pass | No P0 theater finding triggered.") {
		t.Fatalf("rendered packet missing required content:\n%s", string(raw))
	}
}

func TestPacketCommandAppearsInTopLevelHelp(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"--help"}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("help exit %d err=%s", exit, errOut.String())
	}
	if !strings.Contains(out.String(), "sdp-trace packet validate --bundle <file>") ||
		!strings.Contains(out.String(), "sdp-trace packet build-pr --source <github-actions|github-fixture> --out <dir>") ||
		!strings.Contains(out.String(), "sdp-trace packet build-github --github-input <file> --out <file>") ||
		!strings.Contains(out.String(), "sdp-trace packet check-demo --bundle <file>") ||
		!strings.Contains(out.String(), "sdp-trace packet render --bundle <file> --out <file>") {
		t.Fatalf("help missing packet commands:\n%s", out.String())
	}
}

func TestPacketDispatchRequiresKnownPacketSubcommand(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet"}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("packet exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(errOut.String(), "packet requires build-pr, build-github, validate, check-demo, or render") {
		t.Fatalf("packet usage diagnostic changed: %s", errOut.String())
	}
}

func TestParsePacketBuildPROptionsKeepsFlagDefaultsAndRequiredOut(t *testing.T) {
	var errOut bytes.Buffer
	opts, code, ok := parsePacketBuildPROptions([]string{"--out", "packet-out"}, &errOut)
	if !ok || code != 0 {
		t.Fatalf("parse build-pr opts ok=%v code=%d err=%s", ok, code, errOut.String())
	}
	if opts.stringValue("source") != "github-actions" ||
		opts.stringValue("github-event") != "" ||
		opts.stringValue("checks-json") != "" ||
		opts.stringValue("artifacts-json") != "" ||
		opts.stringValue("route-manifest") != "" ||
		opts.stringValue("github-api-url") != "" ||
		opts.stringValue("out") != "packet-out" {
		t.Fatalf("build-pr defaults changed: source=%q event=%q checks=%q artifacts=%q route=%q api=%q out=%q",
			opts.stringValue("source"),
			opts.stringValue("github-event"),
			opts.stringValue("checks-json"),
			opts.stringValue("artifacts-json"),
			opts.stringValue("route-manifest"),
			opts.stringValue("github-api-url"),
			opts.stringValue("out"))
	}

	errOut.Reset()
	opts, code, ok = parsePacketBuildPROptions(nil, &errOut)
	if ok || opts != nil || code != exitUsage {
		t.Fatalf("missing --out parse ok=%v code=%d opts=%v err=%s", ok, code, opts, errOut.String())
	}
	if !strings.Contains(errOut.String(), "packet build-pr requires --out") {
		t.Fatalf("missing --out diagnostic changed: %s", errOut.String())
	}
}

func TestBuildPacketPRResultKeepsPathsAndAggregatesGateErrors(t *testing.T) {
	result, _ := buildPacketPRResult(packet.GitHubPREvidenceInput{}, "packet-out")
	if result.State != packet.StateCannotVerify {
		t.Fatalf("build-pr result state = %s, want cannot_verify", result.State)
	}
	if result.BundlePath != filepath.Join("packet-out", "bundle.json") ||
		result.PacketPath != filepath.Join("packet-out", "change-evidence-packet.md") ||
		result.ResultPath != filepath.Join("packet-out", "build-pr-result.json") {
		t.Fatalf("build-pr paths changed: %+v", result)
	}
	errors := strings.Join(result.Errors, "\n")
	if !strings.Contains(errors, "PC-AGENT-ROUTE cannot verify live route proof") ||
		!strings.Contains(errors, "PC-VERIFICATION cannot verify live CI evidence") {
		t.Fatalf("build-pr result missing live-gate aggregation: %+v", result.Errors)
	}
}

func TestPacketBuildPRGateErrorsPreserveRouteAndVerificationOrder(t *testing.T) {
	bundle := validPacketBundleForCLI()
	setPacketCLIRowState(&bundle, "PC-AGENT-ROUTE", packet.StateCannotVerify, nil, "route missing")
	setPacketCLIRowState(&bundle, "PC-VERIFICATION", packet.StateCannotVerify, nil, "ci missing")
	errors := packetBuildPRGateErrors(bundle)
	want := []string{
		"PC-AGENT-ROUTE cannot verify live route proof: route missing",
		"PC-VERIFICATION cannot verify live CI evidence: ci missing",
	}
	if !reflect.DeepEqual(errors, want) {
		t.Fatalf("gate errors = %#v, want %#v", errors, want)
	}
}

func TestPacketBuildPRRouteErrorsAcceptPassAndPartial(t *testing.T) {
	for _, state := range []string{packet.StatePass, packet.StatePartial} {
		rows := map[string]packet.Row{"PC-AGENT-ROUTE": {State: state, Reason: "kept"}}
		if got := packetBuildPRRouteErrors(rows); got != nil {
			t.Fatalf("route errors for %s = %#v, want nil", state, got)
		}
	}
	rows := map[string]packet.Row{"PC-AGENT-ROUTE": {State: packet.StateCannotVerify, Reason: "route missing"}}
	want := []string{"PC-AGENT-ROUTE cannot verify live route proof: route missing"}
	if got := packetBuildPRRouteErrors(rows); !reflect.DeepEqual(got, want) {
		t.Fatalf("route errors = %#v, want %#v", got, want)
	}
}

func TestPacketBuildPRVerificationErrorsRequirePass(t *testing.T) {
	rows := map[string]packet.Row{"PC-VERIFICATION": {State: packet.StatePass, Reason: "kept"}}
	if got := packetBuildPRVerificationErrors(rows); got != nil {
		t.Fatalf("verification errors for pass = %#v, want nil", got)
	}
	rows = map[string]packet.Row{"PC-VERIFICATION": {State: packet.StatePartial, Reason: "ci partial"}}
	want := []string{"PC-VERIFICATION cannot verify live CI evidence: ci partial"}
	if got := packetBuildPRVerificationErrors(rows); !reflect.DeepEqual(got, want) {
		t.Fatalf("verification errors = %#v, want %#v", got, want)
	}
}

func TestValidPRInputSourceAcceptsOnlyKnownSources(t *testing.T) {
	cases := map[string]bool{
		"github-actions": true,
		"github-fixture": true,
		"":               false,
		"local":          false,
		"github":         false,
	}
	for source, want := range cases {
		if got := validPRInputSource(source); got != want {
			t.Fatalf("validPRInputSource(%q) = %v, want %v", source, got, want)
		}
	}
}

func TestPREventPathUsesActionsEnvOnlyWhenEventPathMissing(t *testing.T) {
	getenv := func(name string) string {
		if name != "GITHUB_EVENT_PATH" {
			t.Fatalf("unexpected getenv key %s", name)
		}
		return "env-event.json"
	}
	if got := prEventPath("github-actions", "", getenv); got != "env-event.json" {
		t.Fatalf("actions missing event path = %q", got)
	}
	if got := prEventPath("github-actions", "explicit-event.json", getenv); got != "explicit-event.json" {
		t.Fatalf("actions explicit event path = %q", got)
	}
	if got := prEventPath("github-fixture", "fixture-event.json", getenv); got != "fixture-event.json" {
		t.Fatalf("fixture explicit event path = %q", got)
	}
	if got := prEventPath("github-fixture", "", func(string) string {
		t.Fatal("fixture mode must not read GITHUB_EVENT_PATH")
		return "env-event.json"
	}); got != "" {
		t.Fatalf("fixture missing event path = %q, want empty", got)
	}
}

func TestLoadPRInputSourceEventRejectsUnsupportedAndMissingEvent(t *testing.T) {
	opts := packetBuildPROptionsForTest("unsupported", "")
	source, _, err := loadPRInputSourceEvent(opts)
	if err == nil || !strings.Contains(err.Error(), `unsupported packet build-pr source "unsupported"`) {
		t.Fatalf("unsupported source=%q err=%v", source, err)
	}

	opts = packetBuildPROptionsForTest("github-fixture", "")
	source, _, err = loadPRInputSourceEvent(opts)
	if err == nil || err.Error() != "missing GitHub event JSON" {
		t.Fatalf("missing event source=%q err=%v", source, err)
	}

	eventPath := writePRFixtureEventForTest(t, t.TempDir())
	opts = packetBuildPROptionsForTest("github-fixture", eventPath)
	source, event, err := loadPRInputSourceEvent(opts)
	if err != nil {
		t.Fatalf("load fixture event: %v", err)
	}
	if source != "github-fixture" || event.WorkflowRunID != "1001" || event.PullRequest.Number != 38 {
		t.Fatalf("loaded source=%q event=%+v", source, event)
	}
}

func TestReadOptionalPREvidenceKeepsErrorPrefixes(t *testing.T) {
	root := t.TempDir()
	badChecks := filepath.Join(root, "bad-checks.json")
	if err := os.WriteFile(badChecks, []byte(`{`), 0o644); err != nil {
		t.Fatalf("write bad checks: %v", err)
	}
	opts := packetBuildPROptionsForTest("github-fixture", "")
	opts.setString("checks-json", badChecks)
	var input packet.GitHubPREvidenceInput
	if err := readOptionalPREvidence(opts, &input); err == nil || !strings.Contains(err.Error(), "read checks json:") {
		t.Fatalf("checks error = %v", err)
	}

	badArtifacts := filepath.Join(root, "bad-artifacts.json")
	if err := os.WriteFile(badArtifacts, []byte(`{`), 0o644); err != nil {
		t.Fatalf("write bad artifacts: %v", err)
	}
	opts = packetBuildPROptionsForTest("github-fixture", "")
	opts.setString("artifacts-json", badArtifacts)
	if err := readOptionalPREvidence(opts, &input); err == nil || !strings.Contains(err.Error(), "read artifacts json:") {
		t.Fatalf("artifacts error = %v", err)
	}
}

func TestBuildPRInputFromOptionsAppliesOptionalEvidenceAndRoute(t *testing.T) {
	root := t.TempDir()
	eventPath := writePRFixtureEventForTest(t, root)
	checksPath := filepath.Join(root, "checks.json")
	artifactsPath := filepath.Join(root, "artifacts.json")
	routePath := filepath.Join(root, "route.json")
	writeTestJSON(t, checksPath, []packet.GitHubCheck{{
		Name: "ci", URL: "https://github.com/example/repo/actions/runs/1001", Conclusion: "success",
	}})
	writeTestJSON(t, artifactsPath, []packet.GitHubArtifact{{
		Name: "packet", Resolver: "https://github.com/example/repo/actions/runs/1001/artifacts/2002", RetainedForm: "external_ref",
	}})
	writeTestJSON(t, routePath, packet.GitHubPREvidenceInput{
		AgentRouteRefs:         []string{"recorder:run-1"},
		AgentRouteComponents:   []string{"opencode", "gsd", "minimax-m2.5"},
		AgentRouteDigest:       "sha256:route",
		AgentRouteEvidenceKind: "harness_route_observation",
		PromptBoundary:         packet.PromptBoundary{Text: "Implement and verify."},
	})

	opts := packetBuildPROptionsForTest("github-fixture", eventPath)
	opts.setString("checks-json", checksPath)
	opts.setString("artifacts-json", artifactsPath)
	opts.setString("route-manifest", routePath)
	input, err := buildPRInputFromOptions(opts)
	if err != nil {
		t.Fatalf("build PR input: %v", err)
	}
	if input.WorkflowRunID != "1001" || input.PR.Number != 38 || input.CommitRange.Base != "base-sha" || input.CommitRange.Head != "head-sha" {
		t.Fatalf("event-derived input changed: %+v", input)
	}
	if len(input.Checks) != 1 || input.Checks[0].Name != "ci" {
		t.Fatalf("checks not applied: %+v", input.Checks)
	}
	if len(input.Artifacts) != 1 || input.Artifacts[0].Name != "packet" {
		t.Fatalf("artifacts not applied: %+v", input.Artifacts)
	}
	if !reflect.DeepEqual(input.AgentRouteRefs, []string{"recorder:run-1"}) ||
		input.AgentRouteDigest != "sha256:route" ||
		input.PromptBoundary.Text != "Implement and verify." {
		t.Fatalf("route not applied: %+v", input)
	}

	badRoute := filepath.Join(root, "bad-route.json")
	if err := os.WriteFile(badRoute, []byte(`{`), 0o644); err != nil {
		t.Fatalf("write bad route: %v", err)
	}
	opts = packetBuildPROptionsForTest("github-fixture", eventPath)
	opts.setString("route-manifest", badRoute)
	if _, err := buildPRInputFromOptions(opts); err == nil || !strings.Contains(err.Error(), "read route manifest:") {
		t.Fatalf("route manifest error = %v", err)
	}
}

func TestGitHubPRInputFromEventUsesActionsEnvRunID(t *testing.T) {
	event := prFixtureEventForTest(t)
	input := githubPRInputFromEvent(event, "github-actions", func(name string) string {
		if name != "GITHUB_RUN_ID" {
			t.Fatalf("unexpected getenv key %s", name)
		}
		return "actions-run-99"
	})
	if input.SchemaVersion != "github-pr-evidence-input.v0" ||
		!input.RequirePromptBoundary ||
		input.WorkflowRunID != "actions-run-99" {
		t.Fatalf("actions input metadata changed: %+v", input)
	}
}

func TestGitHubPRInputFromEventUsesFixtureRunID(t *testing.T) {
	event := prFixtureEventForTest(t)
	input := githubPRInputFromEvent(event, "github-fixture", func(string) string {
		return "actions-run-99"
	})
	if input.WorkflowRunID != "1001" || !input.RequirePromptBoundary {
		t.Fatalf("fixture input metadata changed: %+v", input)
	}
}

func TestGitHubPRFromEventMapsPRFields(t *testing.T) {
	pr := githubPRFromEvent(prFixtureEventForTest(t))
	if pr.Number != 38 ||
		pr.URL != "https://github.com/example/repo/pull/38" ||
		pr.Title != "Demo feature" ||
		pr.BodyRef != "https://github.com/example/repo/pull/38" ||
		pr.Author != "developer" ||
		pr.BaseRef != "main" ||
		pr.HeadRef != "feature" ||
		pr.HeadSHA != "head-sha" {
		t.Fatalf("PR mapping changed: %+v", pr)
	}
}

func TestGitHubCommitRangeFromEventMapsSHAsAndDiffURL(t *testing.T) {
	commitRange := githubCommitRangeFromEvent(prFixtureEventForTest(t))
	if commitRange.Base != "base-sha" ||
		commitRange.Head != "head-sha" ||
		commitRange.ChangedFilesRef != "https://github.com/example/repo/pull/38/files" {
		t.Fatalf("commit range mapping changed: %+v", commitRange)
	}
}

func TestRunPacketBuildPRWritesCannotVerifyJSONForInputFailure(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := runPacketBuildPR([]string{"--source", "github-fixture", "--out", t.TempDir()}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("build-pr exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"state": "cannot_verify"`) ||
		!strings.Contains(out.String(), "missing GitHub event JSON") {
		t.Fatalf("build-pr input failure did not emit cannot_verify JSON: %s", out.String())
	}
}

func TestRenderPacketPRMarkdownDowngradesResultOnFailure(t *testing.T) {
	result := packet.BuildPRResult{State: packet.StatePass}
	markdown, ok := renderPacketPRMarkdown(packet.Bundle{}, &result)
	if ok || markdown != "" {
		t.Fatalf("render invalid bundle ok=%v markdown=%q", ok, markdown)
	}
	if result.State != packet.StateCannotVerify || len(result.Errors) == 0 {
		t.Fatalf("render failure did not downgrade result: %+v", result)
	}
}

func TestWritePacketPRArtifactsWritesCannotVerifyJSONOnRenderFailure(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	result := packet.BuildPRResult{State: packet.StatePass}
	exit := writePacketPRArtifacts(t.TempDir(), packet.Bundle{}, result, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("write artifacts exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"state": "cannot_verify"`) {
		t.Fatalf("render failure did not emit cannot_verify JSON: %s", out.String())
	}
}

func TestWritePacketPRFilesCreatesOutputDirectoryAndArtifacts(t *testing.T) {
	outDir := filepath.Join(t.TempDir(), "packet-out")
	result := packet.BuildPRResult{
		BundlePath: filepath.Join(outDir, "bundle.json"),
		PacketPath: filepath.Join(outDir, "change-evidence-packet.md"),
		ResultPath: filepath.Join(outDir, "build-pr-result.json"),
	}
	var errOut bytes.Buffer
	if !writePacketPRFiles(outDir, validPacketBundleForCLI(), result, "packet markdown", &errOut) {
		t.Fatalf("write packet files failed: %s", errOut.String())
	}
	info, err := os.Stat(outDir)
	if err != nil {
		t.Fatalf("stat output dir: %v", err)
	}
	if got := info.Mode().Perm(); got&0o700 != 0o700 {
		t.Fatalf("output dir mode = %#o, want owner rwx bits preserved", got)
	}
	files := packetPRArtifactFiles(validPacketBundleForCLI(), result, "packet markdown")
	wantLabels := []string{"write packet bundle", "write packet markdown", "write packet result"}
	for i, want := range wantLabels {
		if files[i].label != want {
			t.Fatalf("artifact label[%d] = %q, want %q", i, files[i].label, want)
		}
	}
	for _, name := range []string{"bundle.json", "change-evidence-packet.md", "build-pr-result.json"} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected artifact %s: %v", name, err)
		}
	}
}

func TestWritePacketPRArtifactFilesStopsAfterFirstFailure(t *testing.T) {
	calls := []string{}
	files := []packetPRArtifactFile{
		{label: "first", write: func() error {
			calls = append(calls, "first")
			return nil
		}},
		{label: "second", write: func() error {
			calls = append(calls, "second")
			return os.ErrPermission
		}},
		{label: "third", write: func() error {
			calls = append(calls, "third")
			return nil
		}},
	}
	var errOut bytes.Buffer
	if writePacketPRArtifactFiles(files, &errOut) {
		t.Fatal("write artifact files succeeded after a failing artifact")
	}
	if !reflect.DeepEqual(calls, []string{"first", "second"}) {
		t.Fatalf("artifact write calls = %#v, want first two only", calls)
	}
	if !strings.Contains(errOut.String(), "second:") {
		t.Fatalf("stderr missing failing artifact label: %s", errOut.String())
	}
}

func TestPacketBuildGitHubCLI(t *testing.T) {
	root := repoRootForTest(t)
	inputPath := filepath.Join(root, "examples", "change-evidence-packet", "github-input.json")
	outPath := filepath.Join(t.TempDir(), "bundle.json")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "build-github", "--github-input", inputPath, "--out", outPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("build-github exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	var bundle packet.Bundle
	readTestJSON(t, outPath, &bundle)
	if bundle.Packet.AuthoringMethod != packet.AuthoringToolGenerated ||
		bundle.Packet.BundleRef != bundle.Manifest.BundleID ||
		bundle.Manifest.PacketDigest == "" {
		t.Fatalf("generated bundle not bound/tool-generated: %+v", bundle)
	}
}

func TestPacketBuildPRFixtureCLIWritesArtifacts(t *testing.T) {
	root := t.TempDir()
	eventPath := filepath.Join(root, "event.json")
	checksPath := filepath.Join(root, "checks.json")
	artifactsPath := filepath.Join(root, "artifacts.json")
	routePath := filepath.Join(root, "route.json")
	outDir := filepath.Join(root, "packet-out")
	writeTestJSON(t, eventPath, map[string]any{
		"workflow_run_id": "1001",
		"pull_request": map[string]any{
			"number":   38,
			"html_url": "https://github.com/example/repo/pull/38",
			"title":    "Demo feature",
			"body_ref": "https://github.com/example/repo/pull/38",
			"diff_url": "https://github.com/example/repo/pull/38/files",
			"user":     map[string]any{"login": "developer"},
			"base":     map[string]any{"ref": "main", "sha": "base-sha"},
			"head":     map[string]any{"ref": "feature", "sha": "head-sha"},
		},
	})
	writeTestJSON(t, checksPath, []packet.GitHubCheck{{
		Name:         "ci",
		URL:          "https://github.com/example/repo/actions/runs/1001",
		Conclusion:   "success",
		ArtifactRefs: []string{"packet-artifacts"},
	}})
	writeTestJSON(t, artifactsPath, []packet.GitHubArtifact{{
		Name:         "packet-artifacts",
		Resolver:     "https://github.com/example/repo/actions/runs/1001/artifacts/2002",
		RetainedForm: "external_ref",
		Digest:       "sha256:artifact",
	}})
	writeTestJSON(t, routePath, packet.GitHubPREvidenceInput{
		AgentRouteRefs:         []string{"recorder:run-1"},
		AgentRouteComponents:   []string{"opencode", "gsd", "minimax-m2.5"},
		AgentRouteDigest:       "sha256:route",
		AgentRouteEvidenceKind: "harness_route_observation",
		PromptBoundary:         packet.PromptBoundary{Text: "Implement the feature and run tests."},
		Reviews:                []packet.GitHubReview{{Reviewer: "pi", Resolver: "external:review", State: packet.StatePass}},
	})
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"packet", "build-pr",
		"--source", "github-fixture",
		"--github-event", eventPath,
		"--checks-json", checksPath,
		"--artifacts-json", artifactsPath,
		"--route-manifest", routePath,
		"--out", outDir,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("build-pr exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	for _, name := range []string{"bundle.json", "change-evidence-packet.md", "build-pr-result.json"} {
		if _, err := os.Stat(filepath.Join(outDir, name)); err != nil {
			t.Fatalf("expected %s: %v", name, err)
		}
	}
	var bundle packet.Bundle
	readTestJSON(t, filepath.Join(outDir, "bundle.json"), &bundle)
	if got := rowStateForCLI(bundle, "PC-VERIFICATION"); got != packet.StatePass {
		t.Fatalf("PC-VERIFICATION state = %s", got)
	}
	if !strings.Contains(bundle.Packet.Rows[4].Summary, "1001") {
		t.Fatalf("verification row does not bind workflow run: %+v", bundle.Packet.Rows[4])
	}
	entry := packetEntryForCLI(bundle, "artifact:packet-artifacts")
	if entry.Actor != "ci_packet_builder" || entry.WriteAuthority != "ci_generated" || entry.SourceRef != "1001" {
		t.Fatalf("artifact authority metadata missing: %+v", entry)
	}
}

func TestPacketBuildPRFixtureFailsClosedWithoutPromptBoundary(t *testing.T) {
	root := t.TempDir()
	eventPath := filepath.Join(root, "event.json")
	checksPath := filepath.Join(root, "checks.json")
	artifactsPath := filepath.Join(root, "artifacts.json")
	writeTestJSON(t, eventPath, map[string]any{
		"workflow_run_id": "1001",
		"pull_request": map[string]any{
			"number":   38,
			"html_url": "https://github.com/example/repo/pull/38",
			"title":    "Demo feature",
			"user":     map[string]any{"login": "developer"},
			"base":     map[string]any{"ref": "main", "sha": "base-sha"},
			"head":     map[string]any{"ref": "feature", "sha": "head-sha"},
		},
	})
	writeTestJSON(t, checksPath, []packet.GitHubCheck{{Name: "ci", URL: "https://github.com/example/repo/actions/runs/1001", Conclusion: "success", ArtifactRefs: []string{"packet-artifacts"}}})
	writeTestJSON(t, artifactsPath, []packet.GitHubArtifact{{Name: "packet-artifacts", Resolver: "https://github.com/example/repo/actions/runs/1001/artifacts/2002", RetainedForm: "external_ref", Digest: "sha256:artifact"}})
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"packet", "build-pr",
		"--source", "github-fixture",
		"--github-event", eventPath,
		"--checks-json", checksPath,
		"--artifacts-json", artifactsPath,
		"--out", filepath.Join(root, "packet-out"),
	}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("build-pr exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "prompt boundary evidence missing") {
		t.Fatalf("missing cannot_verify diagnostic: %s", out.String())
	}
}

func TestPacketBuildPRActionsDiscoversArtifacts(t *testing.T) {
	root := t.TempDir()
	eventPath := filepath.Join(root, "event.json")
	checksPath := filepath.Join(root, "checks.json")
	routePath := filepath.Join(root, "route.json")
	outDir := filepath.Join(root, "packet-out")
	writeTestJSON(t, eventPath, map[string]any{
		"pull_request": map[string]any{
			"number":   38,
			"html_url": "https://github.com/example/repo/pull/38",
			"title":    "Demo feature",
			"body_ref": "https://github.com/example/repo/pull/38",
			"diff_url": "https://github.com/example/repo/pull/38/files",
			"user":     map[string]any{"login": "developer"},
			"base":     map[string]any{"ref": "main", "sha": "base-sha"},
			"head":     map[string]any{"ref": "feature", "sha": "head-sha"},
		},
	})
	writeTestJSON(t, checksPath, []packet.GitHubCheck{{
		Name:         "build-and-test",
		URL:          "https://github.com/example/repo/actions/runs/1001",
		Conclusion:   "success",
		ArtifactRefs: []string{"evidence-bundles"},
	}})
	writeTestJSON(t, routePath, packet.GitHubPREvidenceInput{
		AgentRouteRefs:         []string{"recorder:run-1"},
		AgentRouteComponents:   []string{"opencode", "gsd", "minimax-m2.5"},
		AgentRouteDigest:       "sha256:route",
		AgentRouteEvidenceKind: "harness_route_observation",
		PromptBoundary:         packet.PromptBoundary{Text: "Implement the feature and run tests."},
		Reviews:                []packet.GitHubReview{{Reviewer: "pi", Resolver: "external:review", State: packet.StatePass}},
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "" {
			t.Fatalf("loopback test server received bearer token: %q", r.Header.Get("Authorization"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"artifacts":[{"id":2002,"name":"evidence-bundles","expired":false,"expires_at":"2026-08-10T00:00:00Z","archive_download_url":"https://api.github.com/artifacts/2002/zip"}]}`))
	}))
	defer server.Close()
	t.Setenv("GITHUB_EVENT_PATH", eventPath)
	t.Setenv("GITHUB_REPOSITORY", "example/repo")
	t.Setenv("GITHUB_RUN_ID", "1001")
	t.Setenv("GITHUB_JOB", "build-and-test")
	t.Setenv("GITHUB_SERVER_URL", "https://github.com")
	t.Setenv("GITHUB_TOKEN", "test-token")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"packet", "build-pr",
		"--source", "github-actions",
		"--github-api-url", server.URL,
		"--checks-json", checksPath,
		"--route-manifest", routePath,
		"--out", outDir,
	}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("build-pr exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	var bundle packet.Bundle
	readTestJSON(t, filepath.Join(outDir, "bundle.json"), &bundle)
	if got := rowStateForCLI(bundle, "PC-VERIFICATION"); got != packet.StatePass {
		t.Fatalf("PC-VERIFICATION state = %s", got)
	}
	entry := packetEntryForCLI(bundle, "artifact:evidence-bundles")
	if entry.Resolver == "" || entry.SourceRef != "1001" {
		t.Fatalf("discovered artifact entry missing binding: %+v", entry)
	}
}

func TestPacketBuildPRActionsFailsClosedWithoutCheckEvidence(t *testing.T) {
	root := t.TempDir()
	eventPath := filepath.Join(root, "event.json")
	routePath := filepath.Join(root, "route.json")
	writeTestJSON(t, eventPath, map[string]any{
		"pull_request": map[string]any{
			"number":   38,
			"html_url": "https://github.com/example/repo/pull/38",
			"title":    "Demo feature",
			"body_ref": "https://github.com/example/repo/pull/38",
			"diff_url": "https://github.com/example/repo/pull/38/files",
			"user":     map[string]any{"login": "developer"},
			"base":     map[string]any{"ref": "main", "sha": "base-sha"},
			"head":     map[string]any{"ref": "feature", "sha": "head-sha"},
		},
	})
	writeTestJSON(t, routePath, packet.GitHubPREvidenceInput{
		AgentRouteRefs:         []string{"recorder:run-1"},
		AgentRouteComponents:   []string{"opencode", "gsd", "minimax-m2.5"},
		AgentRouteDigest:       "sha256:route",
		AgentRouteEvidenceKind: "harness_route_observation",
		PromptBoundary:         packet.PromptBoundary{Text: "Implement the feature and run tests."},
		Reviews:                []packet.GitHubReview{{Reviewer: "pi", Resolver: "external:review", State: packet.StatePass}},
	})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"artifacts":[{"id":2002,"name":"evidence-bundles","expired":false,"expires_at":"2026-08-10T00:00:00Z","archive_download_url":"https://api.github.com/artifacts/2002/zip"}]}`))
	}))
	defer server.Close()
	t.Setenv("GITHUB_EVENT_PATH", eventPath)
	t.Setenv("GITHUB_REPOSITORY", "example/repo")
	t.Setenv("GITHUB_RUN_ID", "1001")
	t.Setenv("GITHUB_SERVER_URL", "https://github.com")
	t.Setenv("GITHUB_TOKEN", "test-token")

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"packet", "build-pr",
		"--source", "github-actions",
		"--github-api-url", server.URL,
		"--route-manifest", routePath,
		"--out", filepath.Join(root, "packet-out"),
	}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("build-pr exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "missing GitHub check or workflow run evidence") {
		t.Fatalf("missing check-evidence diagnostic: %s", out.String())
	}
}

func TestPacketBuildPRRejectsUnsafeGitHubAPIURL(t *testing.T) {
	root := t.TempDir()
	eventPath := filepath.Join(root, "event.json")
	outDir := filepath.Join(root, "packet")
	writeTestJSON(t, eventPath, map[string]any{
		"pull_request": map[string]any{
			"number":   38,
			"html_url": "https://github.com/example/repo/pull/38",
			"body_ref": "https://github.com/example/repo/pull/38",
			"diff_url": "https://github.com/example/repo/pull/38/files",
			"user":     map[string]any{"login": "developer"},
			"base":     map[string]any{"ref": "main", "sha": "base-sha"},
			"head":     map[string]any{"ref": "feature", "sha": "head-sha"},
		},
	})
	t.Setenv("GITHUB_EVENT_PATH", eventPath)
	t.Setenv("GITHUB_REPOSITORY", "example/repo")
	t.Setenv("GITHUB_RUN_ID", "1001")
	t.Setenv("GITHUB_TOKEN", "test-token")
	t.Setenv("GITHUB_SERVER_URL", "https://github.com")

	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{
		"packet", "build-pr",
		"--source", "github-actions",
		"--github-api-url", "http://evil.example/api",
		"--out", outDir,
	}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("build-pr exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "HTTPS is required") {
		t.Fatalf("expected unsafe URL error in JSON output, got err=%s out=%s", errOut.String(), out.String())
	}
}

func TestHydrateGitHubActionsEvidenceSkipsFixtureSource(t *testing.T) {
	input := packet.GitHubPREvidenceInput{}
	err := hydrateGitHubActionsEvidence("github-fixture", "://bad", &input, func(name string) string {
		t.Fatalf("fixture source should not read env var %s", name)
		return ""
	})
	if err != nil {
		t.Fatalf("hydrateGitHubActionsEvidence() error = %v", err)
	}
	if len(input.Artifacts) != 0 {
		t.Fatalf("fixture source artifacts = %+v, want none", input.Artifacts)
	}
}

func TestHydrateGitHubActionsArtifactsKeepsExplicitArtifacts(t *testing.T) {
	input := packet.GitHubPREvidenceInput{
		Artifacts: []packet.GitHubArtifact{{
			Name:         "explicit-evidence",
			Resolver:     "file://artifacts/evidence.zip",
			RetainedForm: "archive",
		}},
	}
	err := hydrateGitHubActionArtifacts("://bad", &input, func(name string) string {
		t.Fatalf("explicit artifacts should not read env var %s", name)
		return ""
	})
	if err != nil {
		t.Fatalf("hydrateGitHubActionArtifacts() error = %v", err)
	}
	if got := input.Artifacts[0].Name; got != "explicit-evidence" {
		t.Fatalf("artifact name = %q, want explicit-evidence", got)
	}
}

func TestHydrateGitHubActionsArtifactsBackfillsDiscoveredArtifacts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"artifacts":[{"id":2002,"name":"evidence-bundles","expired":false,"expires_at":"2026-08-10T00:00:00Z","archive_download_url":"https://api.github.com/artifacts/2002/zip"}]}`))
	}))
	defer server.Close()

	input := packet.GitHubPREvidenceInput{}
	err := hydrateGitHubActionArtifacts(server.URL, &input, githubActionsArtifactTestEnv)
	if err != nil {
		t.Fatalf("hydrateGitHubActionArtifacts() error = %v", err)
	}
	if len(input.Artifacts) != 1 {
		t.Fatalf("artifact count = %d, want 1", len(input.Artifacts))
	}
	if got := input.Artifacts[0].Name; got != "evidence-bundles" {
		t.Fatalf("artifact name = %q, want evidence-bundles", got)
	}
	if input.Artifacts[0].Resolver == "" {
		t.Fatalf("artifact resolver is empty: %+v", input.Artifacts[0])
	}
}

func TestHydrateGitHubActionsEvidencePropagatesLiveArtifactErrors(t *testing.T) {
	input := packet.GitHubPREvidenceInput{}
	err := hydrateGitHubActionsEvidence("github-actions", "://bad", &input, githubActionsArtifactTestEnv)
	if err == nil {
		t.Fatal("hydrateGitHubActionsEvidence() error = nil, want URL validation error")
	}
	if !strings.Contains(err.Error(), "unsafe GitHub API URL") {
		t.Fatalf("error = %q, want unsafe GitHub API URL", err.Error())
	}
}

func TestReadOptionalPRRouteKeepsOptionalJSONBehavior(t *testing.T) {
	root := t.TempDir()
	emptyRoute, err := readOptionalPRRoute("")
	if err != nil {
		t.Fatalf("read empty optional route: %v", err)
	}
	if !reflect.DeepEqual(emptyRoute, packet.GitHubPREvidenceInput{}) {
		t.Fatalf("empty route = %+v, want empty", emptyRoute)
	}

	_, err = readOptionalPRRoute(filepath.Join(root, "missing.json"))
	if err == nil {
		t.Fatal("read missing explicit route error = nil")
	}

	routePath := filepath.Join(root, "route.json")
	writeTestJSON(t, routePath, packet.GitHubPREvidenceInput{
		AgentRouteRefs: []string{"recorder:run-1"},
		Reviews:        []packet.GitHubReview{{Reviewer: "pi", Resolver: "external:review", State: packet.StatePass}},
	})
	route, err := readOptionalPRRoute(routePath)
	if err != nil {
		t.Fatalf("read route: %v", err)
	}
	if !reflect.DeepEqual(route.AgentRouteRefs, []string{"recorder:run-1"}) {
		t.Fatalf("route refs = %+v", route.AgentRouteRefs)
	}
	if len(route.Reviews) != 1 || route.Reviews[0].Reviewer != "pi" {
		t.Fatalf("route reviews = %+v", route.Reviews)
	}
}

func TestApplyPRRouteOnlyOverwritesRouteAndReviewFields(t *testing.T) {
	input := packet.GitHubPREvidenceInput{
		SchemaVersion: "github-pr-evidence-input.v0",
		PR: packet.GitHubPR{
			Number:  38,
			URL:     "https://github.com/example/repo/pull/38",
			Title:   "Demo",
			Author:  "developer",
			BaseRef: "main",
			HeadRef: "feature",
			HeadSHA: "head-sha",
		},
		CommitRange:            packet.GitHubCommitRange{Base: "base-sha", Head: "head-sha", ChangedFilesRef: "https://github.com/example/repo/pull/38/files"},
		Checks:                 []packet.GitHubCheck{{Name: "verify", URL: "https://github.com/example/repo/actions/runs/1001", Conclusion: "success"}},
		Artifacts:              []packet.GitHubArtifact{{Name: "evidence-bundles", Resolver: "https://api.github.com/artifacts/1/zip", RetainedForm: "external_ref"}},
		AgentRouteRefs:         []string{"old-route"},
		AgentRouteComponents:   []string{"old-component"},
		AgentRouteDigest:       "sha256:old",
		AgentRouteEvidenceKind: "old-kind",
		PromptBoundary:         packet.PromptBoundary{Text: "old prompt"},
		IntegrationActions:     []packet.IntegrationAction{{Kind: "old", Actor: "old", Resolver: "old"}},
		Reviews:                []packet.GitHubReview{{Reviewer: "old", Resolver: "old", State: packet.StateFail}},
	}
	identityBefore := input.PR
	commitBefore := input.CommitRange
	checksBefore := append([]packet.GitHubCheck(nil), input.Checks...)
	artifactsBefore := append([]packet.GitHubArtifact(nil), input.Artifacts...)

	route := packet.GitHubPREvidenceInput{
		PR:                     packet.GitHubPR{Number: 99, URL: "https://github.com/wrong/repo/pull/99"},
		Checks:                 []packet.GitHubCheck{{Name: "wrong"}},
		Artifacts:              []packet.GitHubArtifact{{Name: "wrong"}},
		AgentRouteRefs:         []string{"recorder:run-1"},
		AgentRouteComponents:   []string{"opencode", "minimax"},
		AgentRouteDigest:       "sha256:route",
		AgentRouteEvidenceKind: "harness_route_observation",
		PromptBoundary:         packet.PromptBoundary{Text: "new prompt", Digest: "sha256:prompt"},
		IntegrationActions:     []packet.IntegrationAction{{Kind: "merge", Actor: "bot", Resolver: "external:action"}},
		Reviews:                []packet.GitHubReview{{Reviewer: "pi", Resolver: "external:review", State: packet.StatePass}},
	}

	applyPRRoute(&input, route)

	if !reflect.DeepEqual(input.PR, identityBefore) {
		t.Fatalf("PR identity changed: %+v", input.PR)
	}
	if !reflect.DeepEqual(input.CommitRange, commitBefore) {
		t.Fatalf("commit range changed: %+v", input.CommitRange)
	}
	if !reflect.DeepEqual(input.Checks, checksBefore) {
		t.Fatalf("checks changed: %+v", input.Checks)
	}
	if !reflect.DeepEqual(input.Artifacts, artifactsBefore) {
		t.Fatalf("artifacts changed: %+v", input.Artifacts)
	}
	if !reflect.DeepEqual(input.AgentRouteRefs, route.AgentRouteRefs) ||
		!reflect.DeepEqual(input.AgentRouteComponents, route.AgentRouteComponents) ||
		input.AgentRouteDigest != route.AgentRouteDigest ||
		input.AgentRouteEvidenceKind != route.AgentRouteEvidenceKind ||
		!reflect.DeepEqual(input.PromptBoundary, route.PromptBoundary) ||
		!reflect.DeepEqual(input.IntegrationActions, route.IntegrationActions) ||
		!reflect.DeepEqual(input.Reviews, route.Reviews) {
		t.Fatalf("route-controlled fields not applied: %+v", input)
	}
}

func TestGitHubActionsArtifactsBackfillsRetainedArtifacts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"artifacts":[{"id":2002,"name":"evidence-bundles","expired":false,"expires_at":"2026-08-10T00:00:00Z","archive_download_url":"https://api.github.com/artifacts/2002/zip"},{"id":2003,"name":"expired-evidence","expired":true,"expires_at":"2026-01-01T00:00:00Z","archive_download_url":"https://api.github.com/artifacts/2003/zip"}]}`))
	}))
	defer server.Close()

	artifacts, err := githubActionsArtifacts(server.URL, githubActionsArtifactTestEnv)
	if err != nil {
		t.Fatalf("githubActionsArtifacts() error = %v", err)
	}
	if len(artifacts) != 1 {
		t.Fatalf("artifact count = %d, want 1 retained artifact", len(artifacts))
	}
	if artifacts[0].Name != "evidence-bundles" || artifacts[0].RetainedForm != "external_ref" {
		t.Fatalf("retained artifact = %+v", artifacts[0])
	}
	if artifacts[0].Resolver == "" || artifacts[0].ExpiresAt != "2026-08-10T00:00:00Z" {
		t.Fatalf("retained artifact binding = %+v", artifacts[0])
	}
}

func TestGitHubActionsArtifactsFailsClosedWithoutRetainedArtifacts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"artifacts":[{"id":2003,"name":"expired-evidence","expired":true,"expires_at":"2026-01-01T00:00:00Z","archive_download_url":"https://api.github.com/artifacts/2003/zip"}]}`))
	}))
	defer server.Close()

	_, err := githubActionsArtifacts(server.URL, githubActionsArtifactTestEnv)
	if err == nil {
		t.Fatal("githubActionsArtifacts() error = nil, want empty retained artifact error")
	}
	if !strings.Contains(err.Error(), "GitHub Actions artifact discovery returned no retained artifacts") {
		t.Fatalf("error = %q, want empty retained artifact diagnostic", err.Error())
	}
}

func TestGitHubActionsArtifactsInvalidContextStopsBeforeFetch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("artifact fetch should not run when context validation fails")
	}))
	defer server.Close()

	_, err := githubActionsArtifacts(server.URL, func(name string) string {
		switch name {
		case "GITHUB_REPOSITORY":
			return ""
		case "GITHUB_RUN_ID":
			return "1001"
		case "GITHUB_TOKEN":
			return "test-token"
		case "GITHUB_SERVER_URL":
			return "https://github.com"
		default:
			return ""
		}
	})
	if err == nil {
		t.Fatal("githubActionsArtifacts() error = nil, want context validation error")
	}
	if !strings.Contains(err.Error(), "missing GITHUB_REPOSITORY or GITHUB_RUN_ID") {
		t.Fatalf("error = %q, want repository/run validation error", err.Error())
	}
}

func TestGitHubActionsArtifactsPropagatesFetchErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unavailable", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	_, err := githubActionsArtifacts(server.URL, githubActionsArtifactTestEnv)
	if err == nil {
		t.Fatal("githubActionsArtifacts() error = nil, want fetch error")
	}
	if !strings.Contains(err.Error(), "list GitHub Actions artifacts: HTTP 503") {
		t.Fatalf("error = %q, want propagated HTTP status", err.Error())
	}
}

func TestNewGitHubActionsArtifactContextBuildsValidatedContext(t *testing.T) {
	ctx, err := newGitHubActionsArtifactContext("http://127.0.0.1:8080/", func(name string) string {
		switch name {
		case "GITHUB_REPOSITORY":
			return "example/repo"
		case "GITHUB_RUN_ID":
			return "1001"
		case "GITHUB_TOKEN":
			return ""
		case "GH_TOKEN":
			return "fallback-token"
		case "GITHUB_SERVER_URL":
			return "https://github.com"
		default:
			return ""
		}
	})
	if err != nil {
		t.Fatalf("newGitHubActionsArtifactContext() error = %v", err)
	}
	if ctx.repo != "example/repo" || ctx.runID != "1001" || ctx.token != "fallback-token" || ctx.apiURL != "http://127.0.0.1:8080" {
		t.Fatalf("context = %+v", ctx)
	}
}

func TestGitHubAPIURLSelectionKeepsPrecedenceAndTrimming(t *testing.T) {
	tests := []struct {
		name   string
		flag   string
		envURL string
		want   string
	}{
		{name: "flag wins and trims", flag: "http://127.0.0.1:8080/", envURL: "https://github.example.com/api/v3", want: "http://127.0.0.1:8080"},
		{name: "env used when flag empty", envURL: "https://github.example.com/api/v3/", want: "https://github.example.com/api/v3"},
		{name: "default public API", want: "https://api.github.com"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := githubAPIURL(tt.flag, func(name string) string {
				switch name {
				case "GITHUB_API_URL":
					return tt.envURL
				case "GITHUB_SERVER_URL":
					if tt.envURL != "" {
						return "https://github.example.com"
					}
					return "https://github.com"
				default:
					return ""
				}
			})
			if err != nil {
				t.Fatalf("githubAPIURL() error = %v", err)
			}
			if got != tt.want {
				t.Fatalf("githubAPIURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGitHubTokenPrefersGitHubTokenThenFallsBack(t *testing.T) {
	got := githubToken(func(name string) string {
		switch name {
		case "GITHUB_TOKEN":
			return "primary-token"
		case "GH_TOKEN":
			return "fallback-token"
		default:
			return ""
		}
	})
	if got != "primary-token" {
		t.Fatalf("githubToken() = %q, want primary-token", got)
	}

	got = githubToken(func(name string) string {
		if name == "GH_TOKEN" {
			return "fallback-token"
		}
		return ""
	})
	if got != "fallback-token" {
		t.Fatalf("githubToken() fallback = %q, want fallback-token", got)
	}
}

func TestValidateGitHubActionsArtifactContextKeepsDiagnostics(t *testing.T) {
	err := validateGitHubActionsArtifactContext(githubActionsArtifactContext{
		token: "token",
	})
	if err == nil || !strings.Contains(err.Error(), "missing GITHUB_REPOSITORY or GITHUB_RUN_ID") {
		t.Fatalf("identity error = %v", err)
	}

	err = validateGitHubActionsArtifactContext(githubActionsArtifactContext{
		repo:  "example/repo",
		runID: "1001",
	})
	if err == nil || !strings.Contains(err.Error(), "missing GITHUB_TOKEN or GH_TOKEN") {
		t.Fatalf("token error = %v", err)
	}
}

func githubActionsArtifactTestEnv(name string) string {
	switch name {
	case "GITHUB_REPOSITORY":
		return "example/repo"
	case "GITHUB_RUN_ID":
		return "1001"
	case "GITHUB_TOKEN":
		return "test-token"
	case "GITHUB_SERVER_URL":
		return "https://github.com"
	default:
		return ""
	}
}

func TestGitHubAPIURLPolicy(t *testing.T) {
	tests := []struct {
		name      string
		apiURL    string
		serverURL string
		wantErr   string
	}{
		{name: "github dot com", apiURL: "https://api.github.com"},
		{name: "github enterprise same host", apiURL: "https://github.example.com/api/v3", serverURL: "https://github.example.com"},
		{name: "loopback test server", apiURL: "http://127.0.0.1:8080"},
		{name: "reject public API for enterprise server", apiURL: "https://api.github.com", serverURL: "https://github.example.com", wantErr: "not the configured GitHub host"},
		{name: "reject plain http remote", apiURL: "http://evil.example/api", wantErr: "HTTPS is required"},
		{name: "reject embedded credentials", apiURL: "https://user:secret@api.github.com", wantErr: "embedded credentials are not allowed"},
		{name: "reject different https host", apiURL: "https://evil.example/api", serverURL: "https://github.example.com", wantErr: "not the configured GitHub host"},
		{name: "reject malformed", apiURL: "://bad", wantErr: "unsafe GitHub API URL"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGitHubAPIURL(tt.apiURL, tt.serverURL)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("validateGitHubAPIURL() error = %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("validateGitHubAPIURL() error = %v, want %q", err, tt.wantErr)
			}
		})
	}
}

func TestGitHubAPIURLPolicyHelpersKeepTrustBoundaries(t *testing.T) {
	parsed, err := parseGitHubAPIURL("http://user:secret@example.invalid/api")
	if err != nil {
		t.Fatalf("parseGitHubAPIURL() error = %v", err)
	}
	err = validateParsedGitHubAPIURL(parsed, "http://user:secret@example.invalid/api", "https://github.com")
	if err == nil || !strings.Contains(err.Error(), "embedded credentials are not allowed") {
		t.Fatalf("mixed-invalid URL error = %v, want embedded credentials", err)
	}

	parsed, err = parseGitHubAPIURL("http://127.0.0.1:8080/api")
	if err != nil {
		t.Fatalf("parse loopback API URL: %v", err)
	}
	if !localHTTPGitHubAPI(parsed) || !loopbackHost("localhost") || !loopbackHost("127.0.0.1") || !loopbackHost("::1") {
		t.Fatalf("loopback helpers rejected local API target")
	}
	if err := validateGitHubAPIURLTrustTarget(parsed, "http://127.0.0.1:8080/api", "https://github.com"); err != nil {
		t.Fatalf("validate loopback API URL: %v", err)
	}

	parsed, err = parseGitHubAPIURL("http://evil.example/api")
	if err != nil {
		t.Fatalf("parse nonlocal API URL: %v", err)
	}
	if localHTTPGitHubAPI(parsed) {
		t.Fatalf("nonlocal HTTP API reported as local")
	}
	if err := requireHTTPSGitHubAPI(parsed, "http://evil.example/api"); err == nil || !strings.Contains(err.Error(), "HTTPS is required") {
		t.Fatalf("requireHTTPSGitHubAPI() error = %v", err)
	}
}

func TestParseGitHubAPIURLKeepsSyntaxDiagnostics(t *testing.T) {
	for _, apiURL := range []string{"", "://bad", "https://"} {
		t.Run(apiURL, func(t *testing.T) {
			_, err := parseGitHubAPIURL(apiURL)
			if err == nil || !strings.Contains(err.Error(), "unsafe GitHub API URL") {
				t.Fatalf("parseGitHubAPIURL(%q) error = %v", apiURL, err)
			}
		})
	}
}

func TestGitHubServerHostKeepsFallbackBehavior(t *testing.T) {
	if got := githubServerHost(""); got != "" {
		t.Fatalf("empty server host = %q, want empty", got)
	}
	if got := githubServerHost("://bad"); got != "" {
		t.Fatalf("malformed server host = %q, want empty", got)
	}
	if got := githubServerHost("https://GitHub.Example.com/org"); got != "github.example.com" {
		t.Fatalf("server host = %q, want github.example.com", got)
	}

	if !publicGitHubServerHost("") || !publicGitHubServerHost("github.com") {
		t.Fatalf("publicGitHubServerHost rejected public GitHub host")
	}
	if !githubAPIHostAllowed("api.github.com", "https://github.com") {
		t.Fatalf("public GitHub API host not allowed")
	}
	if githubAPIHostAllowed("github.com", "https://github.com") {
		t.Fatalf("public GitHub server host should not replace api.github.com")
	}
	if !githubAPIHostAllowed("github.example.com", "https://github.example.com") {
		t.Fatalf("enterprise exact host not allowed")
	}
	if githubAPIHostAllowed("api.github.com", "https://github.example.com") {
		t.Fatalf("public API host allowed for enterprise server")
	}
}

func TestGitHubActionsArtifactsRequestOnlySendsTokenOverHTTPS(t *testing.T) {
	tests := []struct {
		name     string
		apiURL   string
		wantAuth string
	}{
		{name: "https github", apiURL: "https://api.github.com", wantAuth: "Bearer test-token"},
		{name: "http loopback", apiURL: "http://127.0.0.1:8080", wantAuth: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := githubActionsArtifactsRequest(githubActionsArtifactContext{
				apiURL: tt.apiURL,
				repo:   "example/repo",
				runID:  "1001",
				token:  "test-token",
			})
			if err != nil {
				t.Fatalf("githubActionsArtifactsRequest() error = %v", err)
			}
			if got := req.Header.Get("Authorization"); got != tt.wantAuth {
				t.Fatalf("Authorization = %q, want %q", got, tt.wantAuth)
			}
		})
	}
}

func TestGitHubActionsArtifactsRequestKeepsHTTPContract(t *testing.T) {
	req, err := githubActionsArtifactsRequest(githubActionsArtifactContext{
		apiURL: "https://api.github.com/",
		repo:   "example/repo",
		runID:  "1001",
		token:  "test-token",
	})
	if err != nil {
		t.Fatalf("githubActionsArtifactsRequest() error = %v", err)
	}
	if got := req.Method; got != http.MethodGet {
		t.Fatalf("method = %q, want %q", got, http.MethodGet)
	}
	if got := req.URL.String(); got != "https://api.github.com/repos/example/repo/actions/runs/1001/artifacts" {
		t.Fatalf("url = %q", got)
	}
	if got := req.Header.Get("Accept"); got != "application/vnd.github+json" {
		t.Fatalf("Accept = %q", got)
	}
	if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
		t.Fatalf("Authorization = %q", got)
	}
}

func TestGitHubActionsArtifactsAuthorizationFailClosed(t *testing.T) {
	tests := []struct {
		name   string
		apiURL string
	}{
		{name: "malformed", apiURL: "://bad"},
		{name: "plain http", apiURL: "http://github.example.com"},
		{name: "loopback http", apiURL: "http://127.0.0.1:8080"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := githubActionsArtifactsAuthorization(tt.apiURL, "test-token"); got != "" {
				t.Fatalf("authorization = %q, want empty", got)
			}
		})
	}
}

func TestSuccessfulHTTPStatusAcceptsOnly2xx(t *testing.T) {
	for _, status := range []int{200, 204, 299} {
		if !successfulHTTPStatus(status) {
			t.Fatalf("status %d rejected", status)
		}
	}
	for _, status := range []int{199, 300, 404, 500} {
		if successfulHTTPStatus(status) {
			t.Fatalf("status %d accepted", status)
		}
	}
}

func TestDecodeGitHubActionsArtifactsKeepsDiagnostics(t *testing.T) {
	payload, err := decodeGitHubActionsArtifacts(strings.NewReader(`{"artifacts":[{"id":42,"name":"packet","expired":false}]}`))
	if err != nil {
		t.Fatalf("decodeGitHubActionsArtifacts() error = %v", err)
	}
	if len(payload.Artifacts) != 1 || payload.Artifacts[0].ID != 42 || payload.Artifacts[0].Name != "packet" {
		t.Fatalf("decoded payload = %+v", payload)
	}
	_, err = decodeGitHubActionsArtifacts(strings.NewReader(`{`))
	if err == nil || !strings.Contains(err.Error(), "decode GitHub Actions artifacts") {
		t.Fatalf("decode malformed error = %v", err)
	}
}

func TestRetainedGitHubArtifactsKeepsResolverPolicy(t *testing.T) {
	payload := githubActionsArtifactPayload{Artifacts: []githubActionsArtifact{
		{ID: 1, Name: "expired", URL: "https://example.invalid/expired", Expired: true},
		{ID: 2, Name: "direct", URL: "https://example.invalid/direct", Expired: false},
		{ID: 3, Name: "synthesized", Expired: false},
		{Name: "missing-id", Expired: false},
	}}
	artifacts := retainedGitHubArtifacts(payload, githubActionsArtifactContext{
		apiURL: "https://api.github.com/",
		repo:   "example/repo",
	})
	if len(artifacts) != 3 {
		t.Fatalf("retained artifacts = %+v", artifacts)
	}
	if artifacts[0].Name != "direct" || artifacts[0].Resolver != "https://example.invalid/direct" {
		t.Fatalf("direct artifact = %+v", artifacts[0])
	}
	if artifacts[1].Name != "synthesized" || artifacts[1].Resolver != "https://api.github.com/repos/example/repo/actions/artifacts/3/zip" {
		t.Fatalf("synthesized artifact = %+v", artifacts[1])
	}
	if artifacts[2].Name != "missing-id" || artifacts[2].Resolver != "" {
		t.Fatalf("missing-id artifact = %+v", artifacts[2])
	}
	for _, artifact := range artifacts {
		if artifact.RetainedForm != "external_ref" {
			t.Fatalf("retained form = %q", artifact.RetainedForm)
		}
	}
}

func TestPacketBuildPRDoesNotAcceptCheckedInPacketAuthority(t *testing.T) {
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "build-pr", "--out", t.TempDir(), "--bundle", "checked-in-bundle.json"}, &out, &errOut)
	if exit != exitUsage {
		t.Fatalf("build-pr exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(errOut.String(), "unknown flag --bundle") {
		t.Fatalf("expected stale packet input flag rejection, got err=%s out=%s", errOut.String(), out.String())
	}
}

func TestPacketCheckDemoCLIRequiresFirstPacketRouteEvidence(t *testing.T) {
	root := t.TempDir()
	bundle := validPacketBundleForCLI()
	bundlePath := filepath.Join(root, "bundle.json")
	writeTestJSON(t, bundlePath, bundle)
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "check-demo", "--bundle", bundlePath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("check-demo exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), `"state": "pass"`) {
		t.Fatalf("check-demo output missing pass: %s", out.String())
	}

	setPacketCLIRowState(&bundle, "PC-AGENT-ROUTE", packet.StateNotAssessed, nil, "route evidence missing")
	bundle.Manifest.PacketDigest = packet.PacketDigest(bundle.Packet)
	writeTestJSON(t, bundlePath, bundle)
	out.Reset()
	errOut.Reset()
	exit = run([]string{"packet", "check-demo", "--bundle", bundlePath}, &out, &errOut)
	if exit != exitFail {
		t.Fatalf("check-demo missing route exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "PC-AGENT-ROUTE must be pass or partial") {
		t.Fatalf("check-demo output missing route error: %s", out.String())
	}
}

func TestPacketCheckDemoCLIRejectsWeakDemoProof(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name  string
		edit  func(*packet.Bundle)
		error string
	}{
		{
			name: "hand authored",
			edit: func(bundle *packet.Bundle) {
				bundle.Packet.AuthoringMethod = "hand_authored_before_tooling"
			},
			error: "tool_generated authoring_method",
		},
		{
			name: "generic harness route",
			edit: func(bundle *packet.Bundle) {
				setPacketCLIEntryResolver(bundle, "harness:route", "generic harness observation artifact")
				setPacketCLIEntryRouteObservation(bundle, "harness:route", "", nil)
			},
			error: "structured OpenCode/GSD/MiniMax harness route observation",
		},
		{
			name: "expired change evidence access",
			edit: func(bundle *packet.Bundle) {
				setPacketCLIEntryAccess(bundle, "git:change", "expired")
			},
			error: "PC-CHANGE evidence ref",
		},
		{
			name: "expired partial change evidence",
			edit: func(bundle *packet.Bundle) {
				setPacketCLIRowState(bundle, "PC-CHANGE", packet.StatePartial, []string{"git:change"}, "change evidence is partial")
				setPacketCLIEntryExpiresAt(bundle, "git:change", "2026-05-10T12:00:00Z")
			},
			error: "PC-CHANGE evidence ref",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bundle := validPacketBundleForCLI()
			tt.edit(&bundle)
			bundle.Manifest.PacketDigest = packet.PacketDigest(bundle.Packet)
			bundlePath := filepath.Join(root, strings.ReplaceAll(tt.name, " ", "-")+".json")
			writeTestJSON(t, bundlePath, bundle)
			var out bytes.Buffer
			var errOut bytes.Buffer
			exit := run([]string{"packet", "check-demo", "--bundle", bundlePath}, &out, &errOut)
			if exit != exitFail {
				t.Fatalf("check-demo exit %d err=%s out=%s", exit, errOut.String(), out.String())
			}
			if !strings.Contains(out.String(), tt.error) {
				t.Fatalf("check-demo output missing %q: %s", tt.error, out.String())
			}
		})
	}
}

func TestPacketCheckDemoCLIRejectsMinimumBarRegressions(t *testing.T) {
	root := t.TempDir()
	tests := []struct {
		name  string
		edit  func(*packet.Bundle)
		error string
	}{
		{
			name: "too few pass partial rows",
			edit: func(bundle *packet.Bundle) {
				for _, id := range []string{"PC-INITIATOR", "PC-VERIFICATION", "PC-REVIEW", "PC-THEATER", "PC-DECISION", "PC-RESIDUAL-GAPS"} {
					setPacketCLIRowState(bundle, id, packet.StateNotAssessed, nil, "not assessed in regression fixture")
					addPacketCLIGap(bundle, id, packet.StateNotAssessed, "not assessed in regression fixture")
				}
			},
			error: "at least 4 pass or partial rows",
		},
		{
			name: "missing change evidence",
			edit: func(bundle *packet.Bundle) {
				setPacketCLIRowState(bundle, "PC-CHANGE", packet.StatePartial, nil, "missing retained change evidence")
			},
			error: "PC-CHANGE retained evidence refs",
		},
		{
			name: "too many unclosed cannot verify rows",
			edit: func(bundle *packet.Bundle) {
				setPacketCLIRowState(bundle, "PC-AUTHORITY", packet.StateCannotVerify, nil, "cannot verify authority")
				setPacketCLIRowState(bundle, "PC-ATTESTATION", packet.StateCannotVerify, nil, "cannot verify attestation")
				bundle.Packet.ResidualGaps = withoutPacketCLIGaps(bundle.Packet.ResidualGaps, "PC-AUTHORITY", "PC-ATTESTATION")
			},
			error: "at most one cannot_verify row without closure path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bundle := validPacketBundleForCLI()
			tt.edit(&bundle)
			bundle.Manifest.PacketDigest = packet.PacketDigest(bundle.Packet)
			bundlePath := filepath.Join(root, strings.ReplaceAll(tt.name, " ", "-")+".json")
			writeTestJSON(t, bundlePath, bundle)
			var out bytes.Buffer
			var errOut bytes.Buffer
			exit := run([]string{"packet", "check-demo", "--bundle", bundlePath}, &out, &errOut)
			if exit != exitFail {
				t.Fatalf("check-demo exit %d err=%s out=%s", exit, errOut.String(), out.String())
			}
			if !strings.Contains(out.String(), tt.error) {
				t.Fatalf("check-demo output missing %q: %s", tt.error, out.String())
			}
		})
	}
}

func TestPacketCLIValidatesCommittedHappyPathFixture(t *testing.T) {
	root := repoRootForTest(t)
	bundlePath := filepath.Join(root, "examples", "change-evidence-packet", "happy-path.bundle.json")
	packetPath := filepath.Join(t.TempDir(), "packet.md")
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "validate", "--bundle", bundlePath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("validate fixture exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	out.Reset()
	errOut.Reset()
	exit = run([]string{"packet", "render", "--bundle", bundlePath, "--out", packetPath}, &out, &errOut)
	if exit != 0 {
		t.Fatalf("render fixture exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	rendered, err := os.ReadFile(packetPath)
	if err != nil {
		t.Fatalf("read rendered fixture packet: %v", err)
	}
	if !strings.Contains(string(rendered), "PC-THEATER | pass") {
		t.Fatalf("rendered fixture missing clean theater row:\n%s", string(rendered))
	}
}

func TestPacketValidateCLIRejectsInvalidBundle(t *testing.T) {
	root := t.TempDir()
	bundle := validPacketBundleForCLI()
	bundle.Packet.Rows[0].EvidenceRefs = []string{"missing:ref"}
	bundlePath := filepath.Join(root, "bundle.json")
	writeTestJSON(t, bundlePath, bundle)
	var out bytes.Buffer
	var errOut bytes.Buffer
	exit := run([]string{"packet", "validate", "--bundle", bundlePath}, &out, &errOut)
	if exit != exitCannotVerify {
		t.Fatalf("validate exit %d err=%s out=%s", exit, errOut.String(), out.String())
	}
	if !strings.Contains(out.String(), "missing:ref") {
		t.Fatalf("validate output missing error: %s", out.String())
	}
}

func repoRootForTest(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", ".."))
}

func validPacketBundleForCLI() packet.Bundle {
	refs := map[string]string{
		"git:change":        "git",
		"task:pr-body":      "change_host",
		"harness:route":     "harness",
		"ci:run":            "ci",
		"review:packet":     "review",
		"authority:profile": "manual",
		"theater:clean":     "witness",
		"gap:residual":      "manual",
		"decision:owner":    "manual",
	}
	entries := []packet.BundleEntry{}
	for ref, sourceClass := range refs {
		entries = append(entries, packet.BundleEntry{
			Ref:             ref,
			SourceClass:     sourceClass,
			Digest:          "sha256:1111111111111111111111111111111111111111111111111111111111111111",
			RetainedForm:    "raw",
			RedactionStatus: "not_needed",
			Resolver:        "examples/change-evidence-packet/" + ref,
		})
	}
	for i := range entries {
		if entries[i].Ref == "harness:route" {
			entries[i].Resolver = "OpenCode/GSD MiniMax-M2.5 observation artifact"
			entries[i].EvidenceKind = "harness_route_observation"
			entries[i].ObservedComponents = []string{"opencode", "gsd", "minimax-m2.5"}
		}
	}
	bundle := packet.Bundle{
		Packet: packet.Packet{
			PacketVersion:   packet.PacketSchemaVersion,
			PacketID:        "packet-cli-test",
			SourceChange:    packet.SourceChange{Repository: "fall-out-bug/sdp-trace", ChangeID: "PR-38", URL: "https://github.com/fall-out-bug/sdp-trace/pull/38"},
			GeneratedAt:     "2026-05-11T12:00:00Z",
			AuthoringMethod: packet.AuthoringToolGenerated,
			SelectedProfile: "change-host-rich-v0",
			RedactionPolicy: "not_needed",
			BundleRef:       "bundle-cli-test",
			PacketState:     "review_ready",
			Projection:      packet.Projection{Kind: packet.ProjectionCanonical, Canonical: true, ArtifactRef: "packet:markdown"},
			Rows: []packet.Row{
				packetCLIRow("PC-CHANGE", packet.StatePass, "Change metadata is retained.", []string{"git:change"}),
				packetCLIRow("PC-INITIATOR", packet.StatePartial, "PR body task source is retained.", []string{"task:pr-body"}),
				packetCLIRow("PC-AGENT-ROUTE", packet.StatePartial, "OpenCode/GSD route is retained.", []string{"harness:route"}),
				packetCLIRow("PC-MUTATION", packet.StatePass, "Commit range and changed files are retained.", []string{"git:change"}),
				packetCLIRow("PC-VERIFICATION", packet.StatePass, "CI run evidence is retained.", []string{"ci:run"}),
				packetCLIRow("PC-REVIEW", packet.StatePartial, "Review evidence is retained.", []string{"review:packet"}),
				packetCLIRow("PC-AUTHORITY", packet.StateNotAssessed, "Authority was not assessed.", nil),
				packetCLIRow("PC-THEATER", packet.StatePass, "No P0 theater finding triggered.", []string{"theater:clean"}),
				packetCLIRow("PC-ATTESTATION", packet.StateNotAssessed, "Signed or external attestation was not assessed.", nil),
				packetCLIRow("PC-DECISION", packet.StatePass, "Decision owner is retained.", []string{"decision:owner"}),
				packetCLIRow("PC-RESIDUAL-GAPS", packet.StatePartial, "Partial rows remain explicit.", []string{"gap:residual"}),
			},
			ResidualGaps: []packet.ResidualGap{
				{RowID: "PC-INITIATOR", State: packet.StatePartial, Reason: "PR body is weaker than issue binding", ClosureEvidence: "GitHub issue"},
				{RowID: "PC-AGENT-ROUTE", State: packet.StatePartial, Reason: "route is observation-only", ClosureEvidence: "complete route chain"},
				{RowID: "PC-REVIEW", State: packet.StatePartial, Reason: "review is external", ClosureEvidence: "GitHub review"},
				{RowID: "PC-AUTHORITY", State: packet.StateNotAssessed, Reason: "authority profile is outside fixture", ClosureEvidence: "authority profile"},
				{RowID: "PC-ATTESTATION", State: packet.StateNotAssessed, Reason: "signed trust is outside fixture", ClosureEvidence: "signed packet"},
			},
			DecisionOwners: []packet.DecisionOwner{
				{Decision: "merge", Owner: "maintainer", State: packet.StateNotAssessed, Reason: "packet is not approval"},
				{Decision: "release", Owner: "release owner", State: packet.StateNotAssessed, Reason: "packet is not release approval"},
				{Decision: "risk_acceptance", Owner: "risk owner", State: packet.StateNotAssessed, Reason: "packet is not risk acceptance"},
				{Decision: "security_review", Owner: "security owner", State: packet.StateNotAssessed, Reason: "packet is not security review"},
			},
			NonApproval: "This packet does not approve merge, release, compliance, production trust, semantic correctness, or signed external trust.",
		},
		Manifest: packet.BundleManifest{
			SchemaVersion: packet.BundleSchemaVersion,
			BundleID:      "bundle-cli-test",
			Entries:       entries,
		},
	}
	bundle.Manifest.PacketDigest = packet.PacketDigest(bundle.Packet)
	return bundle
}

func packetCLIRow(id, state, summary string, refs []string) packet.Row {
	row := packet.Row{ID: id, State: state, Summary: summary, EvidenceRefs: refs, Owner: "maintainer"}
	if state != packet.StatePass {
		row.Reason = "non-pass state is explicit"
	}
	return row
}

func rowStateForCLI(bundle packet.Bundle, id string) string {
	for _, row := range bundle.Packet.Rows {
		if row.ID == id {
			return row.State
		}
	}
	return ""
}

func packetEntryForCLI(bundle packet.Bundle, ref string) packet.BundleEntry {
	for _, entry := range bundle.Manifest.Entries {
		if entry.Ref == ref {
			return entry
		}
	}
	return packet.BundleEntry{}
}

func setPacketCLIRowState(bundle *packet.Bundle, id, state string, refs []string, reason string) {
	for i := range bundle.Packet.Rows {
		if bundle.Packet.Rows[i].ID == id {
			bundle.Packet.Rows[i].State = state
			bundle.Packet.Rows[i].EvidenceRefs = refs
			bundle.Packet.Rows[i].Reason = reason
			return
		}
	}
}

func setPacketCLIEntryResolver(bundle *packet.Bundle, ref, resolver string) {
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == ref {
			bundle.Manifest.Entries[i].Resolver = resolver
			return
		}
	}
}

func setPacketCLIEntryAccess(bundle *packet.Bundle, ref, access string) {
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == ref {
			bundle.Manifest.Entries[i].ArtifactAccess = access
			return
		}
	}
}

func setPacketCLIEntryExpiresAt(bundle *packet.Bundle, ref, expiresAt string) {
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == ref {
			bundle.Manifest.Entries[i].ExpiresAt = expiresAt
			return
		}
	}
}

func setPacketCLIEntryRouteObservation(bundle *packet.Bundle, ref, evidenceKind string, components []string) {
	for i := range bundle.Manifest.Entries {
		if bundle.Manifest.Entries[i].Ref == ref {
			bundle.Manifest.Entries[i].EvidenceKind = evidenceKind
			bundle.Manifest.Entries[i].ObservedComponents = components
			return
		}
	}
}

func withoutPacketCLIGaps(gaps []packet.ResidualGap, rowIDs ...string) []packet.ResidualGap {
	remove := map[string]bool{}
	for _, rowID := range rowIDs {
		remove[rowID] = true
	}
	filtered := []packet.ResidualGap{}
	for _, gap := range gaps {
		if !remove[gap.RowID] {
			filtered = append(filtered, gap)
		}
	}
	return filtered
}

func addPacketCLIGap(bundle *packet.Bundle, rowID, state, reason string) {
	if rowID == "PC-RESIDUAL-GAPS" {
		return
	}
	for _, gap := range bundle.Packet.ResidualGaps {
		if gap.RowID == rowID {
			return
		}
	}
	bundle.Packet.ResidualGaps = append(bundle.Packet.ResidualGaps, packet.ResidualGap{
		RowID:           rowID,
		State:           state,
		Reason:          reason,
		ClosureEvidence: "regression fixture closure evidence",
	})
}

func packetBuildPROptionsForTest(source, eventPath string) *flagSet {
	opts := &flagSet{name: "packet build-pr"}
	opts.setString("source", source)
	opts.setString("github-event", eventPath)
	opts.setString("checks-json", "")
	opts.setString("artifacts-json", "")
	opts.setString("route-manifest", "")
	opts.setString("github-api-url", "")
	opts.setString("out", "")
	return opts
}

func prFixtureEventForTest(t *testing.T) prFixtureEvent {
	t.Helper()
	event, err := loadPRFixtureEvent(writePRFixtureEventForTest(t, t.TempDir()))
	if err != nil {
		t.Fatalf("load fixture event: %v", err)
	}
	return event
}

func writePRFixtureEventForTest(t *testing.T, root string) string {
	t.Helper()
	eventPath := filepath.Join(root, "event.json")
	writeTestJSON(t, eventPath, map[string]any{
		"workflow_run_id": "1001",
		"pull_request": map[string]any{
			"number":   38,
			"html_url": "https://github.com/example/repo/pull/38",
			"title":    "Demo feature",
			"body_ref": "https://github.com/example/repo/pull/38",
			"diff_url": "https://github.com/example/repo/pull/38/files",
			"user":     map[string]any{"login": "developer"},
			"base":     map[string]any{"ref": "main", "sha": "base-sha"},
			"head":     map[string]any{"ref": "feature", "sha": "head-sha"},
		},
	})
	return eventPath
}
