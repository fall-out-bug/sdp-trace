package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/adaptercapture"
	"github.com/fall_out_bug/sdp-trace/internal/authority"
	"github.com/fall_out_bug/sdp-trace/internal/checkpoint"
	"github.com/fall_out_bug/sdp-trace/internal/ciartifact"
	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/forensic"
	"github.com/fall_out_bug/sdp-trace/internal/harnessobs"
	"github.com/fall_out_bug/sdp-trace/internal/interaction"
	"github.com/fall_out_bug/sdp-trace/internal/managed"
	"github.com/fall_out_bug/sdp-trace/internal/packet"
	"github.com/fall_out_bug/sdp-trace/internal/posture"
	"github.com/fall_out_bug/sdp-trace/internal/prreview"
	"github.com/fall_out_bug/sdp-trace/internal/query"
	"github.com/fall_out_bug/sdp-trace/internal/recorder"
	"github.com/fall_out_bug/sdp-trace/internal/releaseproof"
	"github.com/fall_out_bug/sdp-trace/internal/repoobserver"
	"github.com/fall_out_bug/sdp-trace/internal/telemetry"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"github.com/fall_out_bug/sdp-trace/internal/verifier"
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

const (
	exitFail         = 1
	exitUsage        = 2
	exitCannotVerify = 3
)

var cliStdin io.Reader = os.Stdin
var version = "dev"

type commandHandler func(context.Context, []string, io.Writer, io.Writer) int
type subcommandHandler func([]string, io.Writer, io.Writer) int

var commandHandlers = map[string]commandHandler{
	"version":           runVersion,
	"wrap":              runWrap,
	"run":               runWrappedCommand,
	"dry-run":           runDryRun,
	"preview":           runPreview,
	"doctor":            runDoctor,
	"install":           runInstall,
	"interaction":       runInteraction,
	"observe":           runObserveCommand,
	"harness":           runHarnessCommand,
	"envelope":          runEnvelope,
	"verify":            runVerify,
	"explain":           runExplain,
	"query":             runQuery,
	"query-pack":        runQueryPack,
	"export":            runExport,
	"report":            runReport,
	"gate":              runGate,
	"assess":            runAssess,
	"override":          runOverride,
	"checkpoint":        runCheckpoint,
	"witness":           runWitness,
	"validate-fixtures": runValidateFixtures,
	"release-proof":     runReleaseProof,
	"pr-review":         runPRReview,
	"packet":            runPacket,
}

func runVersion(_ context.Context, _ []string, stdout, _ io.Writer) int {
	fmt.Fprintf(stdout, "sdp-trace %s\n", version)
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if topLevelHelp(args) {
		printUsage(stdout)
		return 0
	}
	// The first token is the only command selector; everything after it stays
	// command-owned so subcommands can preserve their own evidence contract.
	return dispatchCommand(args[0], args[1:], stdout, stderr)
}

func topLevelHelp(args []string) bool {
	return len(args) == 0 || isHelpArg(args[0])
}

func dispatchCommand(cmd string, args []string, stdout, stderr io.Writer) int {
	handler, ok := commandHandlers[cmd]
	if !ok {
		// Unknown commands are usage defects, not verifier or gate verdicts.
		fmt.Fprintf(stderr, "unknown command: %s\n", cmd)
		printUsage(stderr)
		return 1
	}
	return handler(context.Background(), args, stdout, stderr)
}

func runObserveCommand(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runObserve(args, stdout, stderr)
}

func runHarnessCommand(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runHarness(args, stdout, stderr)
}

func runSubcommand(args []string, stdout, stderr io.Writer, label, usage string, handlers map[string]subcommandHandler) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, usage)
		return exitUsage
	}
	if isHelpArg(args[0]) {
		// Help for a command family is local CLI documentation, not an evidence
		// producing operation.
		fmt.Fprintf(stdout, "Usage: sdp-trace %s\n", label)
		return 0
	}
	return dispatchSubcommand(args[0], args[1:], stdout, stderr, label, handlers)
}

func dispatchSubcommand(cmd string, args []string, stdout, stderr io.Writer, label string, handlers map[string]subcommandHandler) int {
	handler, ok := handlers[cmd]
	if !ok {
		// Keep command-family names stable in diagnostics even when the usage label
		// contains argument suffixes.
		fmt.Fprintf(stderr, "unknown %s command: %s\n", subcommandName(label), cmd)
		return exitUsage
	}
	return handler(args, stdout, stderr)
}

func runOptionalSubcommand(args []string, stdout, stderr io.Writer, handlers map[string]subcommandHandler) (int, bool) {
	if len(args) == 0 {
		return 0, false
	}
	handler, ok := handlers[args[0]]
	if !ok {
		// Optional dispatch lets parent commands fall back to flag parsing.
		return 0, false
	}
	return handler(args[1:], stdout, stderr), true
}

func isHelpArg(arg string) bool {
	return arg == "help" || arg == "--help" || arg == "-h"
}

func subcommandName(label string) string {
	if before, _, ok := strings.Cut(label, " "); ok {
		// Help labels can include usage suffixes; dispatch diagnostics should
		// name only the stable subcommand token.
		return before
	}
	return label
}

func rejectRest(opts *flagSet, stderr io.Writer, message string) bool {
	if len(opts.rest()) == 0 {
		return false
	}
	// Positional arguments are rejected before required flags so diagnostics
	// cannot imply that an ignored payload was accepted.
	fmt.Fprintln(stderr, message)
	return true
}

func requireStringFlag(opts *flagSet, stderr io.Writer, flag, message string) bool {
	if strings.TrimSpace(opts.stringValue(flag)) != "" {
		return true
	}
	// Empty string flags are missing evidence inputs even if the flag appeared.
	fmt.Fprintln(stderr, message)
	return false
}

type requiredCLIFlag struct {
	name    string
	message string
}

func requireOnlyFlags(opts *flagSet, stderr io.Writer, restMessage string, required []requiredCLIFlag) bool {
	if rejectRest(opts, stderr, restMessage) {
		return false
	}
	// Required flag checks run only after the command proves it received no
	// positional payload.
	return requireRequiredFlags(opts, stderr, required)
}

func requireRequiredFlags(opts *flagSet, stderr io.Writer, required []requiredCLIFlag) bool {
	for _, flag := range required {
		// Preserve caller-specific messages so command docs and tests can pin the
		// exact missing input.
		if !requireStringFlag(opts, stderr, flag.name, flag.message) {
			return false
		}
	}
	return true
}

func writeJSONPayload(stdout, stderr io.Writer, value any, message string) bool {
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		// JSON encoding failure means the CLI cannot publish a structured result.
		fmt.Fprintf(stderr, "%s: %v\n", message, err)
		return false
	}
	fmt.Fprintf(stdout, "%s\n", payload)
	return true
}

func writeJSONPayloadUnchecked(stdout io.Writer, value any) {
	payload, _ := json.MarshalIndent(value, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
}

func requireNamedValues(values map[string]string, stderr io.Writer, messagePrefix string) bool {
	for flag, value := range values {
		if strings.TrimSpace(value) == "" {
			// The map caller owns flag order; this helper only centralizes the
			// trust-language for missing named inputs.
			fmt.Fprintf(stderr, "%s requires %s\n", messagePrefix, flag)
			return false
		}
	}
	return true
}

func stringExitCode(state string, codes map[string]int, fallback int) int {
	code, ok := codes[state]
	if !ok {
		// Unknown states are lowered by the caller-provided fallback.
		return fallback
	}
	return code
}

var harnessStateExitCodes = map[string]int{
	harnessobs.StatePass:         0,
	harnessobs.StateFail:         1,
	harnessobs.StateNotAssessed:  exitCannotVerify,
	harnessobs.StateCannotVerify: exitCannotVerify,
}

func harnessStateExitCode(state string) int {
	return stringExitCode(state, harnessStateExitCodes, exitCannotVerify)
}

var prReviewHandlers = map[string]subcommandHandler{
	"packet":     runPRReviewPacket,
	"run":        runPRReviewRun,
	"synthesize": runPRReviewSynthesize,
	"validate":   runPRReviewValidate,
	"summarize":  runPRReviewSummarize,
	"check":      runPRReviewCheck,
}

var packetHandlers = map[string]subcommandHandler{
	"build-pr":     runPacketBuildPR,
	"build-github": runPacketBuildGitHub,
	"validate":     runPacketValidate,
	"check-demo":   runPacketCheckDemo,
	"render":       runPacketRender,
}

var packetBuildPRRequiredFlags = []requiredCLIFlag{
	{"out", "packet build-pr requires --out"},
}

var packetBuildGitHubRequiredFlags = []requiredCLIFlag{
	{"github-input", "packet build-github requires --github-input"},
	{"out", "packet build-github requires --out"},
}

var packetValidateRequiredFlags = []requiredCLIFlag{
	{"bundle", "packet validate requires --bundle"},
}

var packetCheckDemoRequiredFlags = []requiredCLIFlag{
	{"bundle", "packet check-demo requires --bundle"},
}

var packetRenderRequiredFlags = []requiredCLIFlag{
	{"bundle", "packet render requires --bundle"},
	{"out", "packet render requires --out"},
}

var prReviewPacketRequiredFlags = []requiredCLIFlag{
	{"out", "pr-review packet requires --out"},
	{"repo-id", "pr-review packet requires --repo-id"},
	{"change-ref", "pr-review packet requires --change-ref"},
	{"base", "pr-review packet requires --base"},
	{"head", "pr-review packet requires --head"},
	{"diff", "pr-review packet requires --diff"},
}

var prReviewPacketStringFlags = []struct {
	name         string
	defaultValue string
}{
	{"out", ""},
	{"repo-id", ""},
	{"change-ref", ""},
	{"base", ""},
	{"head", ""},
	{"diff", ""},
	{"metadata", ""},
	{"context", ""},
	{"verification", ""},
	{"ci-state", prreview.StateNotAssessed},
	{"created-by", "sdp-trace-cli"},
}

func runPRReview(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "pr-review <packet|run|synthesize|validate|summarize|check> [flags]", "pr-review requires packet, run, synthesize, validate, summarize, or check", prReviewHandlers)
}

func runPacket(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runSubcommand(args, stdout, stderr, "packet <build-pr|build-github|validate|check-demo|render> [flags]", "packet requires build-pr, build-github, validate, check-demo, or render", packetHandlers)
}

func runPacketBuildPR(args []string, stdout, stderr io.Writer) int {
	// The CLI layer only parses command intent; packet trust decisions happen
	// after options are converted into a portable evidence input.
	opts, code, ok := parsePacketBuildPROptions(args, stderr)
	if !ok {
		return code
	}
	// build-pr first reconstructs a portable GitHub evidence input; no packet
	// artifact is written until that input and generated bundle validate.
	input, err := buildPRInputFromOptions(opts)
	if err != nil {
		// Input reconstruction failures are emitted as structured cannot_verify
		// output so automation can consume the failed packet attempt.
		result := packet.BuildPRResult{State: packet.StateCannotVerify, Errors: []string{err.Error()}}
		writeJSONPayloadUnchecked(stdout, result)
		return exitCannotVerify
	}
	result, bundle := buildPacketPRResult(input, opts.stringValue("out"))
	if result.State != packet.StatePass {
		// Failed packet gates are surfaced as JSON but not persisted as a
		// complete packet artifact set.
		writeJSONPayloadUnchecked(stdout, result)
		return exitCannotVerify
	}
	return writePacketPRArtifacts(opts.stringValue("out"), bundle, result, stdout, stderr)
}

func parsePacketBuildPROptions(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "packet build-pr"}
	// Defaults model the GitHub Actions runtime; explicit flags are used for
	// replay fixtures and local tests.
	opts.setString("source", "github-actions")
	opts.setString("github-event", "")
	opts.setString("checks-json", "")
	opts.setString("artifacts-json", "")
	opts.setString("route-manifest", "")
	opts.setString("github-api-url", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		// Parser errors happen before any review packet or run artifact exists.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Required flag validation keeps the command contract explicit before any
	// GitHub or filesystem evidence is loaded.
	if !requireOnlyFlags(opts, stderr, "packet build-pr accepts only flags", packetBuildPRRequiredFlags) {
		return nil, exitUsage, false
	}
	// The packet builder is flag-only so replay inputs can be reconstructed from
	// invocation text and artifacts.
	return opts, 0, true
}

func buildPacketPRResult(input packet.GitHubPREvidenceInput, outDir string) (packet.BuildPRResult, packet.Bundle) {
	// Build and validation times are generated together so packet rows and the
	// build result describe the same local replay.
	bundle := packet.BuildFromGitHubInput(input, time.Now().UTC())
	validation := packet.Validate(bundle, time.Now().UTC())
	liveGateErrors := packetBuildPRGateErrors(bundle)
	// Result paths are declared before writes so downstream tools can compare
	// the manifest against actual artifact publication.
	result := packet.BuildPRResult{
		State:      packet.StatePass,
		BundlePath: filepath.Join(outDir, "bundle.json"),
		PacketPath: filepath.Join(outDir, "change-evidence-packet.md"),
		ResultPath: filepath.Join(outDir, "build-pr-result.json"),
		Errors:     append(validation.Errors, liveGateErrors...),
	}
	if len(result.Errors) > 0 {
		// Validation and live-gate defects both lower the build to cannot_verify.
		result.State = packet.StateCannotVerify
	}
	return result, bundle
}

func writePacketPRArtifacts(outDir string, bundle packet.Bundle, result packet.BuildPRResult, stdout, stderr io.Writer) int {
	markdown, ok := renderPacketPRMarkdown(bundle, &result)
	if !ok {
		writeJSONPayloadUnchecked(stdout, result)
		return exitCannotVerify
	}
	// Durable artifacts are written only after both bundle validation and
	// markdown rendering have succeeded.
	if !writePacketPRFiles(outDir, bundle, result, markdown, stderr) {
		return exitCannotVerify
	}
	writeJSONPayloadUnchecked(stdout, result)
	return 0
}

func renderPacketPRMarkdown(bundle packet.Bundle, result *packet.BuildPRResult) (string, bool) {
	markdown, err := packet.RenderMarkdown(bundle)
	if err != nil {
		// Rendering failure downgrades the structured build result instead of
		// leaving callers with stderr-only state.
		result.State = packet.StateCannotVerify
		result.Errors = []string{err.Error()}
		return "", false
	}
	return markdown, true
}

func writePacketPRFiles(outDir string, bundle packet.Bundle, result packet.BuildPRResult, markdown string, stderr io.Writer) bool {
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		// Without the output directory, none of the packet artifacts are durable.
		fmt.Fprintf(stderr, "create packet output dir: %v\n", err)
		return false
	}
	// All packet output files share one directory so downstream PR review can
	// cite a single artifact root.
	return writePacketPRArtifactFiles(packetPRArtifactFiles(bundle, result, markdown), stderr)
}

func writePacketPRArtifactFiles(files []packetPRArtifactFile, stderr io.Writer) bool {
	for _, file := range files {
		// Stop at the first write failure to avoid publishing a partial packet set
		// as if it were complete.
		if !writePacketPRFile(file, stderr) {
			return false
		}
	}
	return true
}

type packetPRArtifactFile struct {
	label string
	write func() error
}

func packetPRArtifactFiles(bundle packet.Bundle, result packet.BuildPRResult, markdown string) []packetPRArtifactFile {
	// The result carries the paths that this file list materializes.
	return []packetPRArtifactFile{
		{label: "write packet bundle", write: func() error { return writeJSONFile(result.BundlePath, bundle) }},
		{label: "write packet markdown", write: func() error { return writeTextFileAtomic(result.PacketPath, markdown) }},
		{label: "write packet result", write: func() error { return writeJSONFile(result.ResultPath, result) }},
	}
}

func writePacketPRFile(file packetPRArtifactFile, stderr io.Writer) bool {
	if err := file.write(); err != nil {
		// Labels name the artifact role without exposing full write internals.
		fmt.Fprintf(stderr, "%s: %v\n", file.label, err)
		return false
	}
	return true
}

func packetBuildPRGateErrors(bundle packet.Bundle) []string {
	rows := map[string]packet.Row{}
	for _, row := range bundle.Packet.Rows {
		// Rows are keyed by packet id so the live gate checks can reference the
		// same identifiers as the rendered packet.
		rows[row.ID] = row
	}
	errors := []string{}
	// Route proof and CI verification are the live readiness rows for this CLI
	// build path.
	errors = append(errors, packetBuildPRRouteErrors(rows)...)
	errors = append(errors, packetBuildPRVerificationErrors(rows)...)
	return errors
}

func packetBuildPRRouteErrors(rows map[string]packet.Row) []string {
	route := rows["PC-AGENT-ROUTE"]
	if route.State == packet.StatePass || route.State == packet.StatePartial {
		return nil
	}
	// Missing route proof blocks PR packet readiness but remains cannot_verify.
	return []string{"PC-AGENT-ROUTE cannot verify live route proof: " + route.Reason}
}

func packetBuildPRVerificationErrors(rows map[string]packet.Row) []string {
	verification := rows["PC-VERIFICATION"]
	if verification.State == packet.StatePass {
		return nil
	}
	// CI evidence must be live enough for the packet row to pass before the build
	// command publishes artifacts.
	return []string{"PC-VERIFICATION cannot verify live CI evidence: " + verification.Reason}
}

func buildPRInputFromOptions(opts *flagSet) (packet.GitHubPREvidenceInput, error) {
	source, event, err := loadPRInputSourceEvent(opts)
	if err != nil {
		return packet.GitHubPREvidenceInput{}, err
	}
	// Event metadata seeds the input before optional local and live evidence is
	// layered in.
	input := githubPRInputFromEvent(event, source, os.Getenv)
	if err := completePRInputFromOptions(opts, source, &input); err != nil {
		// Optional evidence failures still invalidate the whole packet input,
		// because partial PR packets can overstate route or CI readiness.
		return packet.GitHubPREvidenceInput{}, err
	}
	return input, nil
}

func loadPRInputSourceEvent(opts *flagSet) (string, prFixtureEvent, error) {
	source := opts.stringValue("source")
	if !validPRInputSource(source) {
		// Unknown sources are rejected instead of silently falling back to a live
		// runner mode that could cite the wrong event.
		return "", prFixtureEvent{}, fmt.Errorf("unsupported packet build-pr source %q", source)
	}
	// The source mode decides whether the event path is explicit fixture data or
	// the GitHub Actions event file.
	eventPath := prEventPath(source, opts.stringValue("github-event"), os.Getenv)
	if eventPath == "" {
		// Without an event file there is no authoritative PR identity for the
		// packet rows.
		return "", prFixtureEvent{}, errors.New("missing GitHub event JSON")
	}
	event, err := loadPRFixtureEvent(eventPath)
	if err != nil {
		// Bad fixture/event JSON is a source failure, not a partially verified
		// packet.
		return "", prFixtureEvent{}, err
	}
	return source, event, nil
}

func completePRInputFromOptions(opts *flagSet, source string, input *packet.GitHubPREvidenceInput) error {
	if err := readOptionalPREvidence(opts, input); err != nil {
		return err
	}
	// Live GitHub hydration is skipped for fixture mode so local replay remains
	// hermetic.
	if err := hydrateGitHubActionsEvidence(source, opts.stringValue("github-api-url"), input, os.Getenv); err != nil {
		return err
	}
	route, err := readOptionalPRRoute(opts.stringValue("route-manifest"))
	if err != nil {
		return fmt.Errorf("read route manifest: %w", err)
	}
	// Route manifests are optional enrichment; an empty manifest leaves route
	// rows to validate as missing or cannot_verify.
	applyPRRoute(input, route)
	return nil
}

func readOptionalPREvidence(opts *flagSet, input *packet.GitHubPREvidenceInput) error {
	if err := readOptionalJSON(opts.stringValue("checks-json"), &input.Checks); err != nil {
		return fmt.Errorf("read checks json: %w", err)
	}
	// Artifacts can be provided by file or discovered from GitHub, but malformed
	// local artifact JSON is never ignored.
	if err := readOptionalJSON(opts.stringValue("artifacts-json"), &input.Artifacts); err != nil {
		return fmt.Errorf("read artifacts json: %w", err)
	}
	return nil
}

func validPRInputSource(source string) bool {
	return source == "github-actions" || source == "github-fixture"
}

func prEventPath(source string, eventPath string, getenv func(string) string) string {
	if source == "github-actions" && eventPath == "" {
		// GitHub Actions events default to the runner-provided event file when
		// the packet command is not using an explicit fixture path.
		return getenv("GITHUB_EVENT_PATH")
	}
	return eventPath
}

func githubPRInputFromEvent(event prFixtureEvent, source string, getenv func(string) string) packet.GitHubPREvidenceInput {
	// The GitHub event is the PR identity authority for packet construction.
	input := packet.GitHubPREvidenceInput{
		SchemaVersion:         "github-pr-evidence-input.v0",
		PR:                    githubPRFromEvent(event),
		CommitRange:           githubCommitRangeFromEvent(event),
		WorkflowRunID:         getenv("GITHUB_RUN_ID"),
		RequirePromptBoundary: true,
	}
	if source == "github-fixture" {
		// Fixtures carry a synthetic run id because the runner environment is not
		// available during local replay.
		input.WorkflowRunID = event.WorkflowRunID
	}
	// Prompt boundary evidence is required by default because route proof should
	// not be inferred from PR metadata alone.
	return input
}

func githubPRFromEvent(event prFixtureEvent) packet.GitHubPR {
	return packet.GitHubPR{
		// Human-facing PR fields come from the event payload that GitHub signed
		// into the runner context.
		Number:  event.PullRequest.Number,
		URL:     event.PullRequest.HTMLURL,
		Title:   event.PullRequest.Title,
		BodyRef: event.PullRequest.BodyRef,
		Author:  event.PullRequest.User.Login,
		// Branch identity stays separate from commit-range identity so packet
		// rows can explain both names and immutable SHAs.
		BaseRef: event.PullRequest.Base.Ref,
		HeadRef: event.PullRequest.Head.Ref,
		HeadSHA: event.PullRequest.Head.SHA,
	}
}

func githubCommitRangeFromEvent(event prFixtureEvent) packet.GitHubCommitRange {
	// The commit range is copied from the event so packet generation never
	// infers base/head state from the local checkout.
	return packet.GitHubCommitRange{
		Base:            event.PullRequest.Base.SHA,
		Head:            event.PullRequest.Head.SHA,
		ChangedFilesRef: event.PullRequest.DiffURL,
	}
}

func hydrateGitHubActionsEvidence(source string, apiURL string, input *packet.GitHubPREvidenceInput, getenv func(string) string) error {
	if source != "github-actions" {
		// Fixture mode must not make network calls.
		return nil
	}
	if err := hydrateGitHubActionArtifacts(apiURL, input, getenv); err != nil {
		return err
	}
	return nil
}

func hydrateGitHubActionArtifacts(apiURL string, input *packet.GitHubPREvidenceInput, getenv func(string) string) error {
	if len(input.Artifacts) != 0 {
		// Explicit artifact JSON wins over live discovery for replayability.
		return nil
	}
	artifacts, err := githubActionsArtifacts(apiURL, getenv)
	if err != nil {
		return err
	}
	input.Artifacts = artifacts
	return nil
}

func readOptionalPRRoute(path string) (packet.GitHubPREvidenceInput, error) {
	var route packet.GitHubPREvidenceInput
	return route, readOptionalJSON(path, &route)
}

func applyPRRoute(input *packet.GitHubPREvidenceInput, route packet.GitHubPREvidenceInput) {
	// Route manifests overwrite only route/review fields, leaving PR identity and
	// CI evidence anchored to the selected event source.
	input.AgentRouteRefs = route.AgentRouteRefs
	input.AgentRouteComponents = route.AgentRouteComponents
	input.AgentRouteDigest = route.AgentRouteDigest
	input.AgentRouteEvidenceKind = route.AgentRouteEvidenceKind
	input.PromptBoundary = route.PromptBoundary
	input.IntegrationActions = route.IntegrationActions
	input.Reviews = route.Reviews
}

func githubActionsArtifacts(apiURLFlag string, getenv func(string) string) ([]packet.GitHubArtifact, error) {
	ctx, err := newGitHubActionsArtifactContext(apiURLFlag, getenv)
	if err != nil {
		return nil, err
	}
	// Artifact discovery is live GitHub evidence and therefore requires a fully
	// validated request context before the network call.
	payload, err := fetchGitHubActionsArtifacts(ctx)
	if err != nil {
		return nil, err
	}
	artifacts := retainedGitHubArtifacts(payload, ctx)
	if len(artifacts) == 0 {
		// An empty retained set means there is no durable artifact evidence for
		// the packet to cite.
		return nil, errors.New("GitHub Actions artifact discovery returned no retained artifacts")
	}
	return artifacts, nil
}

type githubActionsArtifactContext struct {
	repo   string
	runID  string
	token  string
	apiURL string
}

type githubActionsArtifactPayload struct {
	Artifacts []githubActionsArtifact `json:"artifacts"`
}

type githubActionsArtifact struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Expired   bool   `json:"expired"`
	ExpiresAt string `json:"expires_at"`
	URL       string `json:"archive_download_url"`
}

func newGitHubActionsArtifactContext(apiURLFlag string, getenv func(string) string) (githubActionsArtifactContext, error) {
	apiURL, err := githubAPIURL(apiURLFlag, getenv)
	if err != nil {
		return githubActionsArtifactContext{}, err
	}
	// Context captures only the repository, run id, token, and validated API
	// endpoint needed for artifact listing.
	ctx := githubActionsArtifactContext{
		repo:   getenv("GITHUB_REPOSITORY"),
		runID:  getenv("GITHUB_RUN_ID"),
		token:  githubToken(getenv),
		apiURL: apiURL,
	}
	if err := validateGitHubActionsArtifactContext(ctx); err != nil {
		// Repository/run/token validation happens before URL construction to keep
		// failed discovery local and deterministic.
		return githubActionsArtifactContext{}, err
	}
	return ctx, nil
}

func validateGitHubActionsArtifactContext(ctx githubActionsArtifactContext) error {
	if missingGitHubArtifactIdentity(ctx) {
		// Repository and run id bind artifact discovery to the current PR run.
		return errors.New("missing GITHUB_REPOSITORY or GITHUB_RUN_ID for GitHub Actions artifact discovery")
	}
	if strings.TrimSpace(ctx.token) == "" {
		// Without a token, the command cannot prove artifact availability from the
		// configured run.
		return errors.New("missing GITHUB_TOKEN or GH_TOKEN for GitHub Actions artifact discovery")
	}
	return nil
}

func missingGitHubArtifactIdentity(ctx githubActionsArtifactContext) bool {
	return strings.TrimSpace(ctx.repo) == "" || strings.TrimSpace(ctx.runID) == ""
}

func githubToken(getenv func(string) string) string {
	token := getenv("GITHUB_TOKEN")
	if token == "" {
		// GH_TOKEN is a local CLI fallback only; callers still validate that a
		// token exists before making artifact API requests.
		return getenv("GH_TOKEN")
	}
	return token
}

func githubAPIURL(apiURLFlag string, getenv func(string) string) (string, error) {
	apiURL := strings.TrimSpace(apiURLFlag)
	if apiURL == "" {
		// GitHub Enterprise runners set GITHUB_API_URL; explicit flags remain
		// higher priority for replay fixtures and local verification.
		apiURL = getenv("GITHUB_API_URL")
	}
	if apiURL == "" {
		apiURL = "https://api.github.com"
	}
	if err := validateGitHubAPIURL(apiURL, getenv("GITHUB_SERVER_URL")); err != nil {
		return "", err
	}
	return strings.TrimRight(apiURL, "/"), nil
}

func validateGitHubAPIURL(apiURL, serverURL string) error {
	parsed, err := parseGitHubAPIURL(apiURL)
	if err != nil {
		return err
	}
	// Parse and target validation are split so error text can distinguish syntax
	// from trust-target failures.
	return validateParsedGitHubAPIURL(parsed, apiURL, serverURL)
}

func validateParsedGitHubAPIURL(parsed *url.URL, apiURL, serverURL string) error {
	if parsed.User != nil {
		// Credentials must travel through Authorization headers, never URLs that
		// can leak through logs, errors, or packet context.
		return errors.New("unsafe GitHub API URL: embedded credentials are not allowed")
	}
	return validateGitHubAPIURLTrustTarget(parsed, apiURL, serverURL)
}

func validateGitHubAPIURLTrustTarget(parsed *url.URL, apiURL, serverURL string) error {
	if localHTTPGitHubAPI(parsed) {
		// Loopback HTTP is allowed for hermetic tests; non-local API targets must
		// satisfy HTTPS and host binding before receiving credentials.
		return nil
	}
	if err := requireHTTPSGitHubAPI(parsed, apiURL); err != nil {
		return err
	}
	if githubAPIHostAllowed(strings.ToLower(parsed.Hostname()), serverURL) {
		// Allowed hosts are either public GitHub's API host or the configured
		// Enterprise host.
		return nil
	}
	return fmt.Errorf("unsafe GitHub API URL %q: host is not the configured GitHub host", apiURL)
}

func requireHTTPSGitHubAPI(parsed *url.URL, apiURL string) error {
	if parsed.Scheme == "https" {
		// HTTPS is the only scheme allowed for credential-bearing GitHub calls.
		return nil
	}
	return fmt.Errorf("unsafe GitHub API URL %q: HTTPS is required before sending GitHub credentials", apiURL)
}

func parseGitHubAPIURL(apiURL string) (*url.URL, error) {
	parsed, err := url.Parse(apiURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("unsafe GitHub API URL %q", apiURL)
	}
	// The parsed URL is still untrusted until scheme, credentials, and host are
	// validated by the caller.
	return parsed, nil
}

func localHTTPGitHubAPI(parsed *url.URL) bool {
	return parsed.Scheme == "http" && loopbackHost(strings.ToLower(parsed.Hostname()))
}

func githubAPIHostAllowed(host, serverURL string) bool {
	serverHost := githubServerHost(serverURL)
	if publicGitHubServerHost(serverHost) {
		// Public github.com maps to the public API host; Enterprise hosts must
		// bind exactly to the configured server hostname.
		return host == "api.github.com"
	}
	return host == serverHost
}

func publicGitHubServerHost(serverHost string) bool {
	return serverHost == "" || serverHost == "github.com"
}

func githubServerHost(serverURL string) string {
	if strings.TrimSpace(serverURL) == "" {
		// Empty server URL means public GitHub, handled by githubAPIHostAllowed.
		return ""
	}
	parsed, err := url.Parse(serverURL)
	if err != nil {
		return ""
	}
	return strings.ToLower(parsed.Hostname())
}

func loopbackHost(host string) bool {
	ip := net.ParseIP(host)
	return host == "localhost" || (ip != nil && ip.IsLoopback())
}

func fetchGitHubActionsArtifacts(ctx githubActionsArtifactContext) (githubActionsArtifactPayload, error) {
	req, err := githubActionsArtifactsRequest(ctx)
	if err != nil {
		// Request construction failures happen before credentials leave the
		// process.
		return githubActionsArtifactPayload{}, err
	}
	// This is the only live network fetch in packet build-pr artifact hydration.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return githubActionsArtifactPayload{}, fmt.Errorf("list GitHub Actions artifacts: %w", err)
	}
	defer resp.Body.Close()
	if !successfulHTTPStatus(resp.StatusCode) {
		// Non-2xx responses mean the retained artifact set cannot be verified.
		return githubActionsArtifactPayload{}, fmt.Errorf("list GitHub Actions artifacts: HTTP %d", resp.StatusCode)
	}
	// Decoding is the handoff from live GitHub response bytes to packet evidence
	// candidates; retained-artifact filtering happens after this step.
	return decodeGitHubActionsArtifacts(resp.Body)
}

func githubActionsArtifactsRequest(ctx githubActionsArtifactContext) (*http.Request, error) {
	url := strings.TrimRight(ctx.apiURL, "/") + "/repos/" + ctx.repo + "/actions/runs/" + ctx.runID + "/artifacts"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Keep the GitHub API media type explicit so artifact evidence fetches are
	// stable across server defaults.
	req.Header.Set("Accept", "application/vnd.github+json")
	// Never attach a GitHub token to the loopback HTTP test API path.
	if auth := githubActionsArtifactsAuthorization(ctx.apiURL, ctx.token); auth != "" {
		req.Header.Set("Authorization", auth)
	}
	return req, nil
}

func githubActionsArtifactsAuthorization(apiURL, token string) string {
	parsed, err := parseGitHubAPIURL(apiURL)
	// The caller validates apiURL earlier; this guard keeps token handling fail-closed.
	if err != nil || parsed.Scheme != "https" {
		return ""
	}
	return "Bearer " + token
}

func successfulHTTPStatus(statusCode int) bool {
	return statusCode >= 200 && statusCode <= 299
}

func decodeGitHubActionsArtifacts(reader io.Reader) (githubActionsArtifactPayload, error) {
	var payload githubActionsArtifactPayload
	if err := json.NewDecoder(reader).Decode(&payload); err != nil {
		return githubActionsArtifactPayload{}, fmt.Errorf("decode GitHub Actions artifacts: %w", err)
	}
	// Payload schema validation is minimal here; packet row validation decides
	// whether retained artifact refs are sufficient evidence.
	return payload, nil
}

func retainedGitHubArtifacts(payload githubActionsArtifactPayload, ctx githubActionsArtifactContext) []packet.GitHubArtifact {
	artifacts := []packet.GitHubArtifact{}
	for _, artifact := range payload.Artifacts {
		if artifact.Expired {
			// Expired artifacts are not retained evidence.
			continue
		}
		artifacts = append(artifacts, packet.GitHubArtifact{
			Name:         artifact.Name,
			Resolver:     githubArtifactResolver(artifact, ctx),
			RetainedForm: "external_ref",
			ExpiresAt:    artifact.ExpiresAt,
		})
	}
	// Retained artifacts are external references only; packet validation decides
	// whether they satisfy the required evidence rows.
	return artifacts
}

func githubArtifactResolver(artifact githubActionsArtifact, ctx githubActionsArtifactContext) string {
	if artifact.URL != "" {
		return artifact.URL
	}
	if artifact.ID == 0 {
		// Without a resolver URL or id, downstream packet rows cannot cite the
		// artifact.
		return ""
	}
	return strings.TrimRight(ctx.apiURL, "/") + "/repos/" + ctx.repo + "/actions/artifacts/" + fmt.Sprint(artifact.ID) + "/zip"
}

type prFixtureEvent struct {
	WorkflowRunID string `json:"workflow_run_id"`
	PullRequest   struct {
		Number  int    `json:"number"`
		HTMLURL string `json:"html_url"`
		Title   string `json:"title"`
		BodyRef string `json:"body_ref"`
		DiffURL string `json:"diff_url"`
		User    struct {
			Login string `json:"login"`
		} `json:"user"`
		Base struct {
			Ref string `json:"ref"`
			SHA string `json:"sha"`
		} `json:"base"`
		Head struct {
			Ref string `json:"ref"`
			SHA string `json:"sha"`
		} `json:"head"`
	} `json:"pull_request"`
}

func loadPRFixtureEvent(path string) (prFixtureEvent, error) {
	var event prFixtureEvent
	if err := readOptionalJSON(path, &event); err != nil {
		return event, err
	}
	if event.PullRequest.Number == 0 || strings.TrimSpace(event.PullRequest.HTMLURL) == "" {
		// Fixture events need enough PR identity to build packet rows.
		return event, errors.New("missing pull_request metadata in GitHub event")
	}
	// Remaining PR fields may be empty in focused fixtures; row validation will
	// lower missing evidence as needed.
	return event, nil
}

func readOptionalJSON(path string, target any) error {
	if strings.TrimSpace(path) == "" {
		// Empty optional path means "no supplemental evidence".
		return nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, target)
}

func packetValidationExit(result packet.Validation) int {
	if result.State == packet.StatePass {
		return 0
	}
	// Packet validation failures mean the packet cannot be trusted as evidence,
	// so the CLI reports cannot_verify rather than a generic runtime failure.
	return exitCannotVerify
}

func packetDemoGateExit(result packet.Validation) int {
	if result.State == packet.StatePass {
		return 0
	}
	// Demo gate validation is an expected fail/pass contract for fixtures.
	return exitFail
}

func runPRReviewPacket(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePRReviewPacketArgs(args, stderr)
	if !ok {
		return code
	}
	// Packet build errors mean the review input cannot be reconstructed as
	// evidence, so they lower trust instead of becoming a generic CLI failure.
	packet, err := buildPRReviewPacket(opts, args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	writeIndentedPayload(stdout, packet)
	return 0
}

func parsePRReviewPacketArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "pr-review packet"}
	registerPRReviewPacketFlags(opts)
	if err := opts.parse(args); err != nil {
		// Parser errors are command-shape failures before any packet evidence is
		// copied or hashed.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review packet accepts only flags") {
		// Positional arguments would be hidden packet inputs, so reject them.
		return nil, exitUsage, false
	}
	if err := requirePRReviewPacketInputs(opts); err != nil {
		// Missing packet anchors are usage errors because the packet cannot be
		// constructed at all.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func registerPRReviewPacketFlags(opts *flagSet) {
	// Packet metadata is fully flag-driven so generated review packets can be
	// replayed without hidden process context.
	for _, flag := range prReviewPacketStringFlags {
		opts.setString(flag.name, flag.defaultValue)
	}
}

func buildPRReviewPacket(opts *flagSet, args []string) (prreview.Packet, error) {
	// Repeated context and verification flags remain ordered packet inputs;
	// comma expansion is CLI sugar, not a separate evidence source.
	// BuildPacket owns validation and persistence of the packet directory.
	// CLI parsing only maps flags into the portable packet options shape.
	return prreview.BuildPacket(prReviewPacketOptions(opts, args, opts.stringValue("out")))
}

func prReviewPacketOptions(opts *flagSet, args []string, outDir string) prreview.PacketOptions {
	// The packet directory is caller-selected so decomposed and combined review
	// flows can publish different artifacts without changing packet identity.
	options := prreview.PacketOptions{OutDir: outDir}
	fillPRReviewPacketIdentity(&options, opts)
	fillPRReviewPacketEvidence(&options, opts, args)
	return options
}

func fillPRReviewPacketIdentity(options *prreview.PacketOptions, opts *flagSet) {
	// Repository and change refs are provenance anchors; they identify the
	// review subject without reading ambient git state.
	options.RepoID = opts.stringValue("repo-id")
	options.ChangeRef = opts.stringValue("change-ref")
	// Commit and diff inputs bind the packet to immutable source facts supplied
	// by the caller.
	options.BaseCommit = opts.stringValue("base")
	options.HeadCommit = opts.stringValue("head")
	options.DiffPath = opts.stringValue("diff")
	// Metadata remains optional context, not authority for the review verdict.
	options.MetadataPath = opts.stringValue("metadata")
}

func fillPRReviewPacketEvidence(options *prreview.PacketOptions, opts *flagSet, args []string) {
	// Repeated path flags are reconstructed from raw args so order survives the
	// simple parser's single-value storage.
	options.ContextPaths = repeatedFlagValues(args, "context", opts.stringValue("context"))
	options.VerificationPaths = repeatedFlagValues(args, "verification", opts.stringValue("verification"))
	// CI state and producer are declared packet metadata, not inferred local
	// evidence.
	options.CIState = opts.stringValue("ci-state")
	options.CreatedBy = opts.stringValue("created-by")
}

func runPRReviewRun(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePRReviewRunArgs(args, stderr)
	if !ok {
		return code
	}
	// Reviewer execution can only produce usable evidence when packet, profile,
	// runner allow-list, and work directory are all replayable.
	runs, preview, err := executePRReviewRun(opts, args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	writePRReviewRunOutput(stdout, runs, preview)
	return 0
}

func executePRReviewRun(opts *flagSet, args []string) (prreview.RunSet, *prreview.RunPreview, error) {
	packet, profile, err := readPRReviewPacketAndProfileValues(opts)
	if err != nil {
		return prreview.RunSet{}, nil, err
	}
	// The working directory is part of the runner boundary; nonexistent paths
	// would make external review evidence impossible to reproduce.
	if err := requireDirectory(opts.stringValue("work-dir")); err != nil {
		return prreview.RunSet{}, nil, err
	}
	// Runner allow-list values are reconstructed from raw args so repeated flags
	// cannot be collapsed by the flag parser.
	return prreview.RunReview(packet, profile, prreview.RunOptions{
		OutDir:         opts.stringValue("out"),
		AllowedRunners: allowedRunnerSet(repeatedFlagValues(args, "allow-external-runner", opts.stringValue("allow-external-runner"))),
		Preview:        opts.boolValue("preview"),
		WorkDir:        opts.stringValue("work-dir"),
	})
}

func writePRReviewRunOutput(stdout io.Writer, runs prreview.RunSet, preview *prreview.RunPreview) {
	if preview != nil {
		// Preview mode reports planned runner invocations without implying that
		// any review evidence has been produced.
		writeIndentedPayload(stdout, preview)
		return
	}
	writeIndentedPayload(stdout, runs)
}

func parsePRReviewRunArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "pr-review run"}
	// Runner selection defaults to no external allowance; callers must opt in to
	// every non-default runner family they want recorded.
	// Preview mode shares the same required inputs as an executing review run.
	// That keeps dry-run planning bound to the same packet/profile evidence.
	opts.setString("packet", "")
	opts.setString("profile", "")
	opts.setString("out", "")
	opts.setString("allow-external-runner", "")
	opts.setString("work-dir", ".")
	// Preview is the only boolean because it changes publication, not inputs.
	opts.setBool("preview", false)
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review run accepts only flags") {
		return nil, exitUsage, false
	}
	// Required packet/profile values are checked when the runner tries to load
	// them so read errors can carry file-specific context.
	return opts, 0, true
}

func runPRReviewSynthesize(args []string, stdout, stderr io.Writer) int {
	// Synthesize converts packet/run evidence into a ledger; it never executes
	// reviewers or upgrades review state on its own.
	opts, code, ok := parsePRReviewSynthesizeArgs(args, stderr)
	if !ok {
		return code
	}
	// Synthesis is evidence collation only; unreadable inputs keep the ledger
	// unverifiable rather than producing a partial trust record.
	inputs, err := readPRReviewSynthesisInputs(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// The ledger is built only after all requested inputs have been decoded, so
	// missing optional history cannot masquerade as a successful merge.
	ledger := prreview.SynthesizeLedger(inputs.packet, inputs.runs, inputs.existing)
	if err := prreview.WriteJSON(opts.stringValue("out"), ledger); err != nil {
		// A synthesized ledger that cannot be written is not durable evidence.
		fmt.Fprintln(stderr, err)
		return 1
	}
	// Stdout mirrors the durable artifact so users inspect the same ledger.
	writeIndentedPayload(stdout, ledger)
	return 0
}

func parsePRReviewSynthesizeArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// Synthesis accepts artifact paths only; there is no inline review payload
	// that could evade packet/run validation.
	opts := &flagSet{name: "pr-review synthesize"}
	// The synthesized ledger is a durable artifact, so the output path is
	// required instead of silently writing only to stdout.
	opts.setString("packet", "")
	opts.setString("runs", "")
	opts.setString("existing-ledger", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		// Bad flags fail before any review artifacts are read.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review synthesize accepts only flags") {
		// Synthesis inputs are explicit artifact paths only.
		return nil, exitUsage, false
	}
	// The ledger path is validated before artifact reads so a bad output target
	// cannot waste reviewer evidence processing.
	if err := requireOutputFile("pr-review synthesize", opts.stringValue("out")); err != nil {
		// Synthesis output is mandatory because stdout alone is not a stable
		// review ledger reference.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}

type prReviewSynthesisInputs struct {
	packet   prreview.Packet
	runs     prreview.RunSet
	existing *prreview.Ledger
}

func readPRReviewSynthesisInputs(opts *flagSet) (prReviewSynthesisInputs, error) {
	packet, err := prreview.ReadPacket(opts.stringValue("packet"))
	if err != nil {
		return prReviewSynthesisInputs{}, err
	}
	// Packet, run set, and optional prior ledger are read before synthesis so
	// the output is derived from complete local evidence.
	runs, err := prreview.ReadRunSet(opts.stringValue("runs"))
	if err != nil {
		return prReviewSynthesisInputs{}, err
	}
	existing, err := readOptionalPRReviewLedger(opts.stringValue("existing-ledger"))
	if err != nil {
		return prReviewSynthesisInputs{}, err
	}
	// Existing ledger, when supplied, is an input to synthesis rather than an
	// authority overriding fresh run outputs.
	return prReviewSynthesisInputs{packet: packet, runs: runs, existing: existing}, nil
}

func readOptionalPRReviewLedger(path string) (*prreview.Ledger, error) {
	if path == "" {
		// A missing optional ledger starts synthesis from an empty review record.
		return nil, nil
	}
	return readExistingPRReviewLedger(path)
}

func readExistingPRReviewLedger(path string) (*prreview.Ledger, error) {
	ledger, err := prreview.ReadLedger(path)
	if err != nil {
		return nil, err
	}
	// Return a pointer only after the ledger is fully decoded.
	return &ledger, nil
}

func runPRReviewValidate(args []string, stdout, stderr io.Writer) int {
	// Validate is the CLI gate for review evidence; it lowers unreadable inputs
	// to cannot_verify instead of treating them as a failed review.
	opts, code, ok := parsePRReviewValidateArgs(args, stderr)
	if !ok {
		return code
	}
	// Validation joins independent artifacts and does not trust a ledger unless
	// the packet, profile, and run set can all be loaded.
	inputs, err := readPRReviewValidationInputs(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// Package validation owns the verdict; the CLI only persists and maps it to
	// process status.
	validation := prreview.Validate(inputs.packet, inputs.profile, inputs.runs, inputs.ledger)
	if err := prreview.WriteJSON(opts.stringValue("out"), validation); err != nil {
		// A validation verdict is useful only after it is persisted.
		fmt.Fprintln(stderr, err)
		return 1
	}
	// The terminal payload is a copy of the persisted validation artifact.
	writeIndentedPayload(stdout, validation)
	return prReviewValidationCLIExitCode(validation)
}

func parsePRReviewValidateArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// The validation command requires explicit artifact paths so the resulting
	// JSON can be traced back to concrete packet, profile, run, and ledger files.
	opts := &flagSet{name: "pr-review validate"}
	// Validation output is required so the verdict can be cited by path instead
	// of relying on transient terminal text.
	// Packet/profile/runs/ledger remain independent paths so validation can
	// report malformed or missing evidence at the correct boundary.
	opts.setString("packet", "")
	opts.setString("profile", "")
	opts.setString("runs", "")
	opts.setString("ledger", "")
	// The output path is part of the command contract, not derived from inputs.
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		// Validation cannot begin until every artifact path is parsed.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Rejecting rest arguments before output validation keeps every accepted
	// input represented by a named field in the validation artifact.
	if rejectRest(opts, stderr, "pr-review validate accepts only flags") {
		// Extra positional data would not be represented in validation JSON.
		return nil, exitUsage, false
	}
	// Validate requires a durable output even when the verdict is cannot_verify,
	// because downstream PR gates cite the JSON artifact.
	if err := requireOutputFile("pr-review validate", opts.stringValue("out")); err != nil {
		// Validation output is a machine artifact and must not be implicit.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func prReviewValidationCLIExitCode(validation prreview.Validation) int {
	if reviewValidationExitCode(validation) != 0 {
		// Invalid review evidence cannot support a PR trust claim.
		return exitCannotVerify
	}
	// A zero exit here only means the review packet validated locally.
	return 0
}

type prReviewValidationInputs struct {
	packet  prreview.Packet
	profile prreview.ReviewProfile
	runs    prreview.RunSet
	ledger  prreview.Ledger
}

func readPRReviewValidationInputs(opts *flagSet) (prReviewValidationInputs, error) {
	packet, err := prreview.ReadPacket(opts.stringValue("packet"))
	if err != nil {
		return prReviewValidationInputs{}, err
	}
	// Packet read is separated from profile/run/ledger reads to keep the error
	// boundary precise for callers and tests.
	profile, runs, ledger, err := readPRReviewValidationArtifacts(opts)
	if err != nil {
		return prReviewValidationInputs{}, err
	}
	return prReviewValidationInputs{packet: packet, profile: profile, runs: runs, ledger: ledger}, nil
}

func readPRReviewValidationArtifacts(opts *flagSet) (prreview.ReviewProfile, prreview.RunSet, prreview.Ledger, error) {
	profile, err := prreview.ReadProfile(opts.stringValue("profile"))
	if err != nil {
		return prreview.ReviewProfile{}, prreview.RunSet{}, prreview.Ledger{}, err
	}
	// All validation artifacts are local files; missing or malformed rows keep
	// the final review state from being promoted.
	runs, err := prreview.ReadRunSet(opts.stringValue("runs"))
	if err != nil {
		return prreview.ReviewProfile{}, prreview.RunSet{}, prreview.Ledger{}, err
	}
	// Ledger is read last because it is derived from packet and runs.
	ledger, err := prreview.ReadLedger(opts.stringValue("ledger"))
	return profile, runs, ledger, err
}

func runPRReviewSummarize(args []string, stdout, stderr io.Writer) int {
	// Summaries are UX copies of validation and ledger state, not new proof.
	opts, code, ok := parsePRReviewSummarizeArgs(args, stderr)
	if !ok {
		return code
	}
	// Summary text is rendered from validation plus ledger evidence only; it is
	// not an independent approval source.
	validation, ledger, err := readPRReviewSummaryInputs(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	summary := prreview.Summarize(validation, ledger)
	if code, ok := writeOptionalPRReviewSummary(opts.stringValue("out"), summary, stderr); !ok {
		return code
	}
	// Stdout mirrors the summary even when a durable summary file is requested.
	// Human-readable summary text does not replace validation or ledger JSON.
	fmt.Fprint(stdout, summary)
	return 0
}

func readPRReviewSummaryInputs(opts *flagSet) (prreview.Validation, prreview.Ledger, error) {
	validation, err := prreview.ReadValidation(opts.stringValue("validation"))
	if err != nil {
		return prreview.Validation{}, prreview.Ledger{}, err
	}
	// Ledger loading happens after validation so summary failures identify the
	// missing artifact class without masking the validation path error.
	ledger, err := prreview.ReadLedger(opts.stringValue("ledger"))
	return validation, ledger, err
}

func writeOptionalPRReviewSummary(path, summary string, stderr io.Writer) (int, bool) {
	if path == "" {
		// Summary output is optional; no path means there is no write attempt and
		// no additional failure state.
		return 0, true
	}
	return writePRReviewSummaryFile(path, summary, stderr)
}

func writePRReviewSummaryFile(path, summary string, stderr io.Writer) (int, bool) {
	if err := refuseExistingFile(path); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage, false
	}
	// Summaries are write-once CLI artifacts; refusing overwrite protects the
	// cited review text from accidental replacement.
	if err := os.WriteFile(path, []byte(summary), 0o644); err != nil {
		fmt.Fprintln(stderr, err)
		return 1, false
	}
	return 0, true
}

func parsePRReviewSummarizeArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "pr-review summarize"}
	// Summaries accept only evidence paths and an optional output path; any extra
	// positional text would be unaudited report content.
	opts.setString("validation", "")
	opts.setString("ledger", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review summarize accepts only flags") {
		return nil, exitUsage, false
	}
	// Required inputs are loaded by readPRReviewSummaryInputs so bad paths remain
	// cannot_verify rather than usage-only failures.
	return opts, 0, true
}

func runPRReviewCheck(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parsePRReviewCheckArgs(args, stderr)
	if !ok {
		return code
	}
	// The combined check preserves the packet -> runs -> ledger -> validation
	// sequence so later artifacts can cite earlier ones.
	packet, profile, code, ok := preparePRReviewCheck(opts, args, stderr)
	if !ok {
		return code
	}
	runs, preview, code, ok := executePRReviewCheck(packet, profile, opts, args, stderr)
	if !ok {
		return code
	}
	// Finish handles preview and non-preview publication from the same prepared
	// packet/profile pair.
	return finishPRReviewCheck(opts.stringValue("out"), packet, profile, runs, preview, stdout, stderr)
}

func parsePRReviewCheckArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "pr-review check"}
	registerPRReviewCheckFlags(opts)
	if err := opts.parse(args); err != nil {
		// The combined command still has a flag-only contract.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "pr-review check accepts only flags") {
		// Positional payload would bypass the packet's declared provenance.
		return nil, exitUsage, false
	}
	// Required review anchors are checked after parsing so diagnostics reflect
	// the declared command shape.
	if err := requirePRReviewCheckInputs(opts); err != nil {
		// Missing anchors are caught before any reviewer process can run.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Successful parsing only validates command shape; execution still has to
	// build packet, profile, and run evidence.
	return opts, 0, true
}

func registerPRReviewCheckFlags(opts *flagSet) {
	// Check mode intentionally mirrors packet and run flags so the one-shot path
	// records the same provenance as the decomposed commands.
	registerPRReviewPacketFlags(opts)
	// Profile, runner policy, and work-dir describe the review boundary; preview
	// selects a dry publication path without changing parsed evidence inputs.
	opts.setString("profile", "")
	opts.setString("allow-external-runner", "")
	opts.setString("work-dir", ".")
	// Preview changes publication only; it does not add evidence inputs.
	opts.setBool("preview", false)
}

func requirePRReviewCheckInputs(opts *flagSet) error {
	outDir := opts.stringValue("out")
	if strings.TrimSpace(outDir) == "" {
		// A combined review check needs a directory because it writes multiple
		// artifacts whose paths become later evidence refs.
		return errors.New("pr-review check requires --out")
	}
	return requirePRReviewPacketInputs(opts)
}

func preparePRReviewCheck(opts *flagSet, args []string, stderr io.Writer) (prreview.Packet, prreview.ReviewProfile, int, bool) {
	outDir := opts.stringValue("out")
	// Packet construction is first because every later review artifact is bound
	// to the packet identity and change metadata.
	packet, err := prreview.BuildPacket(prReviewPacketOptions(opts, args, filepath.Join(outDir, "packet")))
	if err != nil {
		// Packet construction failure means no downstream review artifact should
		// be produced.
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, exitCannotVerify, false
	}
	// The profile controls review planes and runner policy; missing profile
	// evidence prevents the combined check from claiming review coverage.
	profile, err := prreview.ReadProfile(opts.stringValue("profile"))
	if err != nil {
		// A malformed profile means required planes and runners cannot be known.
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, exitCannotVerify, false
	}
	if err := requireDirectory(opts.stringValue("work-dir")); err != nil {
		// Review runners operate relative to a concrete working directory.
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, exitCannotVerify, false
	}
	// Profile and directory checks are part of preparation, not a review verdict.
	// A prepared check has all inputs needed for the runner boundary and no
	// persisted artifact has been written yet.
	return packet, profile, 0, true
}

func executePRReviewCheck(packet prreview.Packet, profile prreview.ReviewProfile, opts *flagSet, args []string, stderr io.Writer) (prreview.RunSet, *prreview.RunPreview, int, bool) {
	outDir := opts.stringValue("out")
	// Combined execution writes runs under a stable subdirectory so ledger and
	// validation artifacts can refer to the same run set.
	runs, preview, err := prreview.RunReview(packet, profile, prreview.RunOptions{
		OutDir:         filepath.Join(outDir, "runs"),
		AllowedRunners: allowedRunnerSet(repeatedFlagValues(args, "allow-external-runner", opts.stringValue("allow-external-runner"))),
		Preview:        opts.boolValue("preview"),
		WorkDir:        opts.stringValue("work-dir"),
	})
	if err != nil {
		// Runner failures are recorded as cannot_verify because review evidence is
		// incomplete.
		fmt.Fprintln(stderr, err)
		return prreview.RunSet{}, nil, exitCannotVerify, false
	}
	return runs, preview, 0, true
}

func finishPRReviewCheck(outDir string, packet prreview.Packet, profile prreview.ReviewProfile, runs prreview.RunSet, preview *prreview.RunPreview, stdout, stderr io.Writer) int {
	if writePRReviewCheckPreview(stdout, preview) {
		return 0
	}
	// Non-preview mode persists artifacts before printing the human summary, so
	// the summary never outruns the machine-readable evidence.
	ledger, validation, code, ok := writePRReviewCheckArtifacts(outDir, packet, profile, runs, stderr)
	if !ok {
		return code
	}
	fmt.Fprint(stdout, prreview.Summarize(validation, ledger))
	return reviewValidationExit(validation)
}

func writePRReviewCheckPreview(stdout io.Writer, preview *prreview.RunPreview) bool {
	if preview == nil {
		return false
	}
	// Preview output is intentionally terminal-only planning data, not persisted
	// review evidence.
	writeIndentedPayload(stdout, preview)
	return true
}

func writePRReviewCheckArtifacts(outDir string, packet prreview.Packet, profile prreview.ReviewProfile, runs prreview.RunSet, stderr io.Writer) (prreview.Ledger, prreview.Validation, int, bool) {
	if !writePRReviewJSON(filepath.Join(outDir, "runs", "results.json"), runs, stderr) {
		// Run-set persistence is the first durable artifact in non-preview mode.
		return prreview.Ledger{}, prreview.Validation{}, 1, false
	}
	// Ledger and validation are derived after run persistence so artifact paths
	// and in-memory validation stay in the same review cycle.
	ledger := prreview.SynthesizeLedger(packet, runs, nil)
	validation := prreview.Validate(packet, profile, runs, ledger)
	if !writePRReviewJSON(filepath.Join(outDir, "ledger.json"), ledger, stderr) {
		return prreview.Ledger{}, prreview.Validation{}, 1, false
	}
	// Validation is written after ledger so readers never see a validation file
	// whose cited ledger is missing.
	if !writePRReviewJSON(filepath.Join(outDir, "validation.json"), validation, stderr) {
		return prreview.Ledger{}, prreview.Validation{}, 1, false
	}
	return ledger, validation, 0, true
}

func writePRReviewJSON(path string, value any, stderr io.Writer) bool {
	if err := prreview.WriteJSON(path, value); err != nil {
		// Artifact write failure means the review evidence cannot be cited later.
		fmt.Fprintln(stderr, err)
		return false
	}
	return true
}

func reviewValidationExit(validation prreview.Validation) int {
	if reviewValidationExitCode(validation) != 0 {
		// Synthesis validation failures are trust gaps, not usage errors.
		return exitCannotVerify
	}
	return 0
}

func writeIndentedPayload(stdout io.Writer, value any) {
	payload, _ := json.MarshalIndent(value, "", "  ")
	// Callers only pass values that have already been constructed for output.
	fmt.Fprintf(stdout, "%s\n", payload)
}

func requirePRReviewPacketInputs(opts *flagSet) error {
	for _, flag := range prReviewPacketRequiredFlags {
		if strings.TrimSpace(opts.stringValue(flag.name)) == "" {
			// Required packet fields are provenance anchors, not cosmetic labels.
			return errors.New(flag.message)
		}
	}
	return nil
}

func requireOutputFile(command, path string) error {
	if strings.TrimSpace(path) == "" {
		// Commands that produce artifacts require an explicit destination to
		// avoid pretending stdout-only output is persisted evidence.
		return fmt.Errorf("%s requires --out", command)
	}
	return refuseExistingFile(path)
}

func refuseExistingFile(path string) error {
	info, err := os.Stat(path)
	if err == nil {
		if info.IsDir() {
			// Directories are never valid write-once text artifact targets.
			return fmt.Errorf("output path is a directory: %s", path)
		}
		return fmt.Errorf("output file exists: %s", path)
	}
	if errors.Is(err, os.ErrNotExist) {
		// Missing path is the only acceptable state for write-once artifacts.
		return nil
	}
	return err
}

func requireDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("work-dir: %w", err)
	}
	if !info.IsDir() {
		// Runner working directory must be a directory, not a file path.
		return fmt.Errorf("work-dir is not a directory: %s", path)
	}
	return nil
}

func reviewValidationExitCode(validation prreview.Validation) int {
	return stringExitCode(validation.ReviewCoverageState, reviewValidationExitCodes, 0)
}

var reviewValidationExitCodes = map[string]int{
	prreview.CoverageCannotVerify: exitCannotVerify,
	prreview.CoverageUnresolved:   exitCannotVerify,
}

func readPRReviewPacketAndProfile(opts *flagSet, stderr io.Writer) (prreview.Packet, prreview.ReviewProfile, bool) {
	packet, profile, err := readPRReviewPacketAndProfileValues(opts)
	if err != nil {
		// Packet/profile load failures mean review evidence cannot be replayed.
		fmt.Fprintln(stderr, err)
		return prreview.Packet{}, prreview.ReviewProfile{}, false
	}
	return packet, profile, true
}

func readPRReviewPacketAndProfileValues(opts *flagSet) (prreview.Packet, prreview.ReviewProfile, error) {
	packet, err := prreview.ReadPacket(opts.stringValue("packet"))
	if err != nil {
		return prreview.Packet{}, prreview.ReviewProfile{}, err
	}
	// Profile is loaded after packet so packet path failures stay first.
	profile, err := prreview.ReadProfile(opts.stringValue("profile"))
	if err != nil {
		return prreview.Packet{}, prreview.ReviewProfile{}, err
	}
	return packet, profile, nil
}

func repeatedFlagValues(args []string, key, parsedFallback string) []string {
	prefix := "--" + key + "="
	values := []string{}
	for i := 0; i < len(args); i++ {
		// Raw args preserve repeated flag order that the simple parser collapses.
		values, i = appendRepeatedFlagValue(values, args, i, key, prefix)
	}
	if len(values) == 0 && strings.TrimSpace(parsedFallback) != "" {
		// The parsed fallback covers the single-value case.
		values = append(values, parsedFallback)
	}
	return values
}

func appendRepeatedFlagValue(values []string, args []string, i int, key, prefix string) ([]string, int) {
	arg := args[i]
	if strings.HasPrefix(arg, prefix) {
		// --key=value contributes exactly one ordered value.
		return append(values, strings.TrimPrefix(arg, prefix)), i
	}
	if arg == "--"+key && i+1 < len(args) {
		// --key value consumes the following argument as an ordered value.
		return append(values, args[i+1]), i + 1
	}
	return values, i
}

func allowedRunnerSet(values []string) map[string]bool {
	allowed := map[string]bool{}
	for _, value := range values {
		addAllowedRunnerItems(allowed, value)
	}
	// Empty input intentionally means no local external runners are allowed.
	return allowed
}

func addAllowedRunnerItems(allowed map[string]bool, value string) {
	for _, item := range strings.Split(value, ",") {
		// Runner allow-lists accept comma-separated flags while preserving the
		// normalized set used by review validation.
		addAllowedRunnerItem(allowed, item)
	}
}

func addAllowedRunnerItem(allowed map[string]bool, item string) {
	item = strings.TrimSpace(item)
	if item != "" {
		// Empty allow-list entries are ignored so accidental commas do not create
		// wildcard-like runner names.
		allowed[item] = true
	}
}

func runReleaseProof(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseReleaseProofArgs(args, stderr)
	if !ok {
		return code
	}
	// Release proof is source-bound: the CLI must evaluate the current repo root
	// before writing any proof JSON that downstream gates might cite.
	result, code, ok := evaluateAndWriteReleaseProof(opts, stderr)
	if !ok {
		return code
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return releaseProofExitCode(result.ReleaseVerificationState)
}

func evaluateAndWriteReleaseProof(opts *flagSet, stderr io.Writer) (releaseproof.Verification, int, bool) {
	repoRoot, err := releaseproof.RepoRoot(".")
	if err != nil {
		// Without a repository root, the manifest cannot be tied to an immutable
		// source boundary.
		fmt.Fprintln(stderr, err)
		return releaseproof.Verification{}, exitCannotVerify, false
	}
	return writeReleaseProofResult(repoRoot, opts, stderr)
}

func writeReleaseProofResult(repoRoot string, opts *flagSet, stderr io.Writer) (releaseproof.Verification, int, bool) {
	result, err := releaseproof.Evaluate(repoRoot, opts.stringValue("manifest"), time.Now())
	if err != nil {
		// Evaluation errors leave release proof unverifiable; they are not a
		// successful proof with warnings.
		fmt.Fprintln(stderr, err)
		return releaseproof.Verification{}, exitCannotVerify, false
	}
	if err := releaseproof.Write(opts.stringValue("out"), result); err != nil {
		// A proof that cannot be persisted cannot be referenced by later gates.
		fmt.Fprintln(stderr, err)
		return releaseproof.Verification{}, 1, false
	}
	return result, 0, true
}

func parseReleaseProofArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "release-proof"}
	// The default manifest is an example contract; callers still need an output
	// path so the generated proof is a durable artifact.
	opts.setString("manifest", "examples/contract-foundation/contract-manifest.example.json")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if !requireOnlyFlags(opts, stderr, "release-proof accepts only flags", releaseProofRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}

var releaseProofRequiredFlags = []requiredCLIFlag{
	{"out", "release-proof requires --out"},
}

var releaseProofExitCodes = map[string]int{
	releaseproof.StatePass: 0,
	releaseproof.StateFail: 1,
}

func releaseProofExitCode(state string) int {
	return stringExitCode(state, releaseProofExitCodes, exitCannotVerify)
}

func runInteraction(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	// Interaction commands all materialize or inspect trace artifacts; the
	// router decides only which artifact contract applies.
	if len(args) == 0 {
		// Without a verb there is no interaction evidence contract to apply.
		fmt.Fprintln(stderr, "interaction requires relay, import-transcript, or summarize")
		return exitUsage
	}
	// Interaction commands intentionally share one router so transcript imports,
	// relays, and summaries use the same trace vocabulary.
	handlers := map[string]func([]string, io.Writer, io.Writer) int{
		// Relay needs context because it may execute a forwarded command.
		"relay": func(args []string, stdout, stderr io.Writer) int {
			return runInteractionRelay(ctx, args, stdout, stderr)
		},
		// Import and summarize are pure artifact transforms from the CLI layer.
		"import-transcript": runInteractionImportTranscript,
		"summarize":         runInteractionSummarize,
	}
	handler, ok := handlers[args[0]]
	if !ok {
		// Unknown interaction verbs are command-shape errors, not trace states.
		fmt.Fprintf(stderr, "unknown interaction command: %s\n", args[0])
		return exitUsage
	}
	// The selected handler owns parsing for its artifact shape.
	return handler(args[1:], stdout, stderr)
}

func runInteractionRelay(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	// Relay records stdin plus command metadata before invoking the downstream
	// command, so the feedback event and command boundary stay coupled.
	opts, code, ok := parseInteractionRelayArgs(args, stderr)
	if !ok {
		return code
	}
	// Relay records the interaction before forwarding to the wrapped command, so
	// corrective feedback is not lost when the downstream command fails.
	_, exitCode, err := interaction.Relay(ctx, interactionRelayOptions(opts), cliStdin, stdout, stderr)
	if err != nil {
		// Relay package errors mean the interaction trace could not be recorded
		// with sufficient provenance.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return exitCode
}

func interactionRelayOptions(opts *flagSet) interaction.RelayOptions {
	return interaction.RelayOptions{
		// Identity fields bind the feedback event to a task and actor.
		TaskID:    opts.stringValue("task-id"),
		ActorType: opts.stringValue("actor-type"),
		ActorID:   opts.stringValue("actor-id"),
		// Target and event fields describe the interaction without trusting the
		// forwarded command to provide trace metadata.
		Target:      opts.stringValue("target"),
		EventType:   opts.stringValue("event-type"),
		OperationID: opts.stringValue("operation-id"),
		StageID:     opts.stringValue("stage-id"),
		// Out and Command define the durable trace location and replay boundary.
		Out:     opts.stringValue("out"),
		Command: opts.rest(),
	}
}

func parseInteractionRelayArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// Relay parsing keeps trace identity flags separate from the command after
	// `--`, which is forwarded and recorded exactly.
	opts := newInteractionRelayFlagSet()
	if err := opts.parse(args); err != nil {
		// Flag parse errors happen before the command boundary can be trusted.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// The task and output path are the minimum durable trace coordinates.
	if !requireRequiredFlags(opts, stderr, interactionRelayRequiredFlags) {
		return nil, exitUsage, false
	}
	// Rest arguments are mandatory because relay exists to bind feedback to a
	// concrete forwarded command.
	if !requireRest(opts, stderr, "interaction relay requires forward command after --") {
		return nil, exitUsage, false
	}
	// The remaining args after `--` are forwarded exactly by Relay and become
	// part of the recorded command boundary.
	return opts, 0, true
}

func newInteractionRelayFlagSet() *flagSet {
	opts := &flagSet{name: "interaction relay"}
	// Relay defaults encode a human-to-agent corrective-feedback event; callers
	// override them only when the trace source is more specific.
	for _, flag := range interactionRelayStringFlags {
		opts.setString(flag.name, flag.defaultValue)
	}
	return opts
}

var interactionRelayStringFlags = []struct {
	name         string
	defaultValue string
}{
	{"task-id", ""},
	{"actor-type", "human_user"},
	{"actor-id", ""},
	{"target", "agent"},
	{"event-type", "corrective_feedback"},
	{"operation-id", ""},
	{"stage-id", ""},
	{"out", ""},
}

func requireRest(opts *flagSet, stderr io.Writer, message string) bool {
	if len(opts.rest()) != 0 {
		return true
	}
	// Commands after `--` are part of the replay boundary; missing rest args
	// would record feedback without a target command.
	fmt.Fprintln(stderr, message)
	return false
}

func requireOnlyFlagsCode(opts *flagSet, stderr io.Writer, restMessage string, required []requiredCLIFlag) (int, bool) {
	if !requireOnlyFlags(opts, stderr, restMessage, required) {
		// Keep parser helpers returning CLI usage codes instead of package
		// verifier states.
		return exitUsage, false
	}
	return 0, true
}

var interactionRelayRequiredFlags = []requiredCLIFlag{
	{"task-id", "interaction relay requires --task-id"},
	{"out", "interaction relay requires --out"},
}

func runInteractionImportTranscript(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseInteractionImportTranscriptArgs(args, stderr)
	if !ok {
		return code
	}
	// Transcript import normalizes external rows into trace events before the CLI
	// emits JSON, keeping summary commands independent of source file shape.
	trace, err := importTranscriptFromOptions(opts)
	return writeImportedTranscript(trace, err, stdout, stderr)
}

func writeImportedTranscript(trace interaction.Trace, err error, stdout, stderr io.Writer) int {
	if err != nil {
		// Import failures mean the transcript source cannot be trusted as trace
		// evidence.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	writeJSONPayloadUnchecked(stdout, trace)
	return 0
}

func importTranscriptFromOptions(opts *flagSet) (interaction.Trace, error) {
	// Source identity and source ref are preserved so imported transcript events
	// remain attributable after normalization.
	return interaction.ImportTranscript(interaction.ImportOptions{
		TaskID:      opts.stringValue("task-id"),
		Source:      opts.stringValue("source"),
		SourceRef:   opts.stringValue("source-ref"),
		EventsJSONL: opts.stringValue("events-jsonl"),
		Out:         opts.stringValue("out"),
	})
}

func parseInteractionImportTranscriptArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// Import parsing requires the external transcript path and the target trace
	// output path up front.
	opts := &flagSet{name: "interaction import-transcript"}
	// Imported transcripts require both input and output paths to avoid
	// terminal-only trace evidence.
	opts.setString("source", "")
	opts.setString("source-ref", "")
	opts.setString("task-id", "")
	opts.setString("events-jsonl", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		// Malformed flags stop before transcript rows are read.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if !requireOnlyFlags(opts, stderr, "interaction import-transcript accepts only flags", interactionImportTranscriptRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func runInteractionSummarize(args []string, stdout, stderr io.Writer) int {
	// Summarize reads an existing trace artifact and emits a derived report; it
	// never imports or mutates trace events.
	opts, code, ok := parseInteractionSummarizeArgs(args, stderr)
	if !ok {
		return code
	}
	// Summaries are derived views of trace evidence; unreadable traces keep the
	// command from producing an overconfident report.
	trace, err := interaction.ReadTrace(opts.stringValue("trace"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// Summary generation is pure over the decoded trace.
	summary := interaction.SummarizeTrace(trace)
	if err := writeOptionalJSONFile(opts.stringValue("out"), summary); err != nil {
		// Optional output write failures are publication failures after a valid
		// trace summary was derived.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, summary)
	return 0
}

func parseInteractionSummarizeArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "interaction summarize"}
	// The summary command accepts a trace artifact and optional output only; it
	// does not accept ad hoc report text.
	opts.setString("trace", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if !requireOnlyFlags(opts, stderr, "interaction summarize accepts only flags", interactionSummarizeRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func runEnvelope(_ context.Context, args []string, stdout, stderr io.Writer) int {
	// Envelope summarize is a read-only inspection command over one persisted
	// interaction envelope.
	opts, code, ok := parseEnvelopeSummarizeArgs(args, stderr)
	if !ok {
		return code
	}
	// Envelope summaries are inspection artifacts; unreadable envelopes remain
	// cannot_verify rather than being summarized from partial data.
	envelope, err := interaction.ReadEnvelope(opts.stringValue("envelope"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	// Summary generation is pure over the decoded envelope.
	summary := interaction.SummarizeEnvelope(envelope)
	if err := writeOptionalJSONFile(opts.stringValue("out"), summary); err != nil {
		// Envelope summary output is a derived artifact; failed persistence is
		// distinct from envelope readability.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, summary)
	return 0
}

func writeOptionalJSONFile(path string, value any) error {
	if strings.TrimSpace(path) == "" {
		// Optional JSON outputs are side effects; an omitted path must not change
		// the command verdict.
		return nil
	}
	return writeJSONFile(path, value)
}

func parseEnvelopeSummarizeArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	if len(args) == 0 || args[0] != "summarize" {
		// The envelope namespace currently has one explicit verb, keeping room for
		// future envelope operations without ambiguous positional parsing.
		fmt.Fprintln(stderr, "envelope requires summarize")
		return nil, exitUsage, false
	}
	opts := &flagSet{name: "envelope summarize"}
	// Envelope summaries require a concrete envelope path; optional output is a
	// second copy of the same derived view.
	opts.setString("envelope", "")
	opts.setString("out", "")
	if err := opts.parse(args[1:]); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	code, ok := requireOnlyFlagsCode(opts, stderr, "envelope summarize accepts only flags", envelopeSummarizeRequiredFlags)
	return opts, code, ok
}

var interactionImportTranscriptRequiredFlags = []requiredCLIFlag{
	{"task-id", "interaction import-transcript requires --task-id"},
	{"events-jsonl", "interaction import-transcript requires --events-jsonl"},
	{"out", "interaction import-transcript requires --out"},
}

var interactionSummarizeRequiredFlags = []requiredCLIFlag{
	{"trace", "interaction summarize requires --trace"},
}

var envelopeSummarizeRequiredFlags = []requiredCLIFlag{
	{"envelope", "envelope summarize requires --envelope"},
}

func runAssess(_ context.Context, args []string, stdout, stderr io.Writer) int {
	// Assess accepts either a documented subcommand or a profile-bound verdict
	// run; no positional assessment payloads are accepted.
	if code, ok := runAssessSubcommand(args, stdout, stderr); ok {
		return code
	}
	opts, ok := parseAssessOptions(args, stderr)
	if !ok {
		return exitUsage
	}
	// The profile flag is the assessment boundary: each profile has a distinct
	// evidence shape and exit-code policy.
	handler, ok := assessHandlers()[opts.stringValue("profile")]
	if !ok {
		fmt.Fprintln(stderr, "assess requires --profile adapter-capture, managed-harness, forensic-retention, ci-artifact-observation, or authority-envelope")
		return exitUsage
	}
	return handler(opts, stdout, stderr)
}

func runAssessSubcommand(args []string, stdout, stderr io.Writer) (int, bool) {
	return runOptionalSubcommand(args, stdout, stderr, assessSubcommandHandlers)
}

var assessSubcommandHandlers = map[string]subcommandHandler{
	"preview": runAssessPreview,
	"explain": runAssessExplain,
}

func parseAssessOptions(args []string, stderr io.Writer) (*flagSet, bool) {
	opts := newStringFlagSet("assess", assessStringFlags)
	if err := opts.parse(args); err != nil {
		// Parse failures happen before profile-specific evidence loading.
		fmt.Fprintln(stderr, err)
		return nil, false
	}
	if len(opts.rest()) != 0 {
		// Assessments are entirely flag-addressed so verdict artifacts can be
		// replayed from named evidence inputs.
		fmt.Fprintln(stderr, "assess accepts only flags")
		return nil, false
	}
	return opts, true
}

func newStringFlagSet(name string, flags []string) *flagSet {
	opts := &flagSet{name: name}
	// Shared assess flags keep the command surface stable while each selected
	// profile validates only the inputs it can actually use.
	for _, flag := range flags {
		opts.setString(flag, "")
	}
	return opts
}

var assessStringFlags = []string{
	"profile",
	"out",
	"contract",
	"run",
	"adapter-registry",
	"managed-policy",
	"managed-witness",
	"redaction-policy",
	"artifact-manifest",
	"authority-package",
}

type assessHandler func(*flagSet, io.Writer, io.Writer) int

func assessHandlers() map[string]assessHandler {
	// Handler keys are product profile names, not implementation package names,
	// so CLI output remains aligned with documented assessment profiles.
	return map[string]assessHandler{
		"adapter-capture":         runAdapterCaptureAssess,
		"managed-harness":         runManagedAssess,
		"forensic-retention":      runForensicAssess,
		"ci-artifact-observation": runCIArtifactAssess,
		"authority-envelope":      runAuthorityAssess,
	}
}

func runAdapterCaptureAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireAdapterCaptureAssessInputs(opts, stderr) {
		// Missing durable input/output flags are usage failures before
		// adaptercapture can evaluate run evidence.
		return exitUsage
	}
	// Adapter-capture assessment is run-bound only; missing run evidence is a
	// usage error before any verdict artifact exists.
	input, err := loadAdapterCaptureInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := adaptercapture.Evaluate(input)
	return writeAssessmentArtifact(opts.stringValue("out"), result, stdout, stderr, adapterCaptureExitCode)
}

func runManagedAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireManagedAssessInputs(opts, stderr) {
		// Managed assessment has no implicit defaults for registry, policy, or
		// witness authority.
		return exitUsage
	}
	// Managed-harness assessment joins contract, policy, registry, run, and
	// witness evidence before deriving a trust state.
	input, err := loadManagedInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := managed.Evaluate(input)
	return writeAssessmentArtifact(opts.stringValue("out"), result, stdout, stderr, managedExitCode)
}

func runForensicAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireForensicAssessInputs(opts, stderr) {
		// A forensic verdict without a redaction policy would overclaim retention
		// coverage.
		return exitUsage
	}
	// Forensic retention assessment requires both policy and run evidence so
	// missing redaction rules cannot be treated as passing defaults.
	input, err := loadForensicInput(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := forensic.Evaluate(input)
	return writeAssessmentArtifact(opts.stringValue("out"), result, stdout, stderr, forensicExitCode)
}

func requireAdapterCaptureAssessInputs(opts *flagSet, stderr io.Writer) bool {
	// Adapter capture has only a run source and a durable result path; both are
	// required before evaluation can produce citeable assessment JSON.
	return requireNamedValues(map[string]string{
		"--out": opts.stringValue("out"),
		"--run": opts.stringValue("run"),
	}, stderr, "adapter capture assess")
}

func requireManagedAssessInputs(opts *flagSet, stderr io.Writer) bool {
	// Managed assessment combines five independent evidence inputs. Keeping the
	// named flag list here preserves explicit provenance for each missing input.
	return requireNamedValues(map[string]string{
		"--out":              opts.stringValue("out"),
		"--run":              opts.stringValue("run"),
		"--adapter-registry": opts.stringValue("adapter-registry"),
		"--managed-policy":   opts.stringValue("managed-policy"),
		"--managed-witness":  opts.stringValue("managed-witness"),
	}, stderr, "managed assess")
}

func requireForensicAssessInputs(opts *flagSet, stderr io.Writer) bool {
	// Forensic retention needs a run plus the redaction policy that defines what
	// safe retained evidence means for this assessment.
	return requireNamedValues(map[string]string{
		"--out":              opts.stringValue("out"),
		"--run":              opts.stringValue("run"),
		"--redaction-policy": opts.stringValue("redaction-policy"),
	}, stderr, "forensic assess")
}

func writeAssessmentArtifact[T any](path string, result T, stdout, stderr io.Writer, exitCode func(T) int) int {
	if err := writeJSONFile(path, result); err != nil {
		// Persisted JSON is the assessment authority for downstream gates.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, result)
	return exitCode(result)
}

func runCIArtifactAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireNamedValues(map[string]string{
		"--out":               opts.stringValue("out"),
		"--artifact-manifest": opts.stringValue("artifact-manifest"),
	}, stderr, "ci artifact observation assess") {
		// CI artifact assessment starts from a named manifest, not repository
		// discovery.
		return exitUsage
	}
	var manifest ciartifact.Manifest
	// CI artifact observation starts from a manifest snapshot; unreadable
	// manifests are input failures, not not_assessed observations.
	if err := readJSONFile(opts.stringValue("artifact-manifest"), &manifest); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	result := ciartifact.Evaluate(manifest)
	return writeCIArtifactAssessment(opts, result, stdout, stderr)
}

func writeCIArtifactAssessment(opts *flagSet, result ciartifact.ObservationResult, stdout, stderr io.Writer) int {
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		// Observation JSON is the durable artifact used by later review gates.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, result)
	return ciArtifactExitCode(result)
}

func loadManagedInput(opts *flagSet) (managed.Input, error) {
	contract, err := trace.LoadContract(opts.stringValue("contract"))
	if err != nil {
		return managed.Input{}, err
	}
	// Contract required events are copied so managed evaluation cannot mutate the
	// loaded trace contract.
	return loadManagedInputWithContract(opts, contract)
}

func loadManagedInputWithContract(opts *flagSet, contract trace.Contract) (managed.Input, error) {
	// Managed assessment combines a trace contract with managed-harness artifacts;
	// each JSON input is decoded before the final package input is assembled.
	policy, registry, runEvidence, witness, err := loadManagedJSONInputs(opts)
	if err != nil {
		return managed.Input{}, err
	}
	// Required events are copied into the managed package shape so downstream
	// evaluation cannot mutate the trace contract.
	return managed.Input{
		Contract: managed.Contract{RequiredEventTypes: append([]string(nil), contract.RequiredEvents...)},
		Policy:   policy,
		Registry: registry,
		Run:      runEvidence,
		Witness:  witness,
	}, nil
}

func loadManagedJSONInputs(opts *flagSet) (managed.Policy, managed.Registry, managed.RunEvidence, managed.Witness, error) {
	var policy managed.Policy
	if err := readJSONFile(opts.stringValue("managed-policy"), &policy); err != nil {
		return managed.Policy{}, managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	// Policy is read first because it defines how registry, run, and witness rows
	// will be interpreted.
	registry, runEvidence, witness, err := readManagedJSONInputs(opts)
	if err != nil {
		return managed.Policy{}, managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	return policy, registry, runEvidence, witness, nil
}

func readManagedJSONInputs(opts *flagSet) (managed.Registry, managed.RunEvidence, managed.Witness, error) {
	var registry managed.Registry
	if err := readJSONFile(opts.stringValue("adapter-registry"), &registry); err != nil {
		// Registry read failures stop before run or witness evidence is combined.
		return managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	return readManagedRunAndWitness(opts, registry)
}

func readManagedRunAndWitness(opts *flagSet, registry managed.Registry) (managed.Registry, managed.RunEvidence, managed.Witness, error) {
	var runEvidence managed.RunEvidence
	// Run evidence lives under the run directory contract; callers pass the run
	// root rather than a free-form JSON path.
	if err := readJSONFile(filepath.Join(opts.stringValue("run"), "run.json"), &runEvidence); err != nil {
		return managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	var witness managed.Witness
	if err := readJSONFile(opts.stringValue("managed-witness"), &witness); err != nil {
		return managed.Registry{}, managed.RunEvidence{}, managed.Witness{}, err
	}
	return registry, runEvidence, witness, nil
}

func loadForensicInput(opts *flagSet) (forensic.Input, error) {
	var policy forensic.Policy
	if err := readJSONFile(opts.stringValue("redaction-policy"), &policy); err != nil {
		return forensic.Input{}, err
	}
	var runEvidence forensic.RunEvidence
	// Forensic assessment reuses the run directory contract so raw and sanitized
	// evidence stay tied to the captured run.
	if err := readJSONFile(filepath.Join(opts.stringValue("run"), "run.json"), &runEvidence); err != nil {
		return forensic.Input{}, err
	}
	return forensic.Input{Policy: policy, Run: runEvidence}, nil
}

func loadAdapterCaptureInput(opts *flagSet) (adaptercapture.Input, error) {
	var runEvidence adaptercapture.RunEvidence
	// Adapter capture currently consumes the normalized run evidence file only.
	if err := readJSONFile(filepath.Join(opts.stringValue("run"), "run.json"), &runEvidence); err != nil {
		return adaptercapture.Input{}, err
	}
	return adaptercapture.Input{Run: runEvidence}, nil
}

type managedPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

func runAssessPreview(args []string, stdout, stderr io.Writer) int {
	opts, ok := parseAssessPreviewOptions(args, stderr)
	if !ok {
		return exitUsage
	}
	// Preview is advisory setup output. It must not imply that assessment
	// evidence has been evaluated.
	handler, ok := assessPreviewHandlers()[opts.stringValue("profile")]
	if !ok {
		fmt.Fprintln(stderr, "assess preview requires --profile adapter-capture, managed-harness, forensic-retention, ci-artifact-observation, or authority-envelope")
		return exitUsage
	}
	return handler(opts, stdout)
}

func parseAssessPreviewOptions(args []string, stderr io.Writer) (*flagSet, bool) {
	opts := newStringFlagSet("assess preview", assessPreviewStringFlags)
	if err := opts.parse(args); err != nil {
		// Preview parse failures are usage errors, not assessment verdicts.
		fmt.Fprintln(stderr, err)
		return nil, false
	}
	if len(opts.rest()) != 0 {
		// Preview output is generated from named local paths only; positional
		// prose is not evidence.
		fmt.Fprintln(stderr, "assess preview accepts only flags")
		return nil, false
	}
	return opts, true
}

var assessPreviewStringFlags = []string{
	"profile",
	"out",
	"run",
	"adapter-registry",
	"managed-policy",
	"managed-witness",
	"redaction-policy",
	"artifact-manifest",
	"authority-package",
}

type assessPreviewHandler func(*flagSet, io.Writer) int

func assessPreviewHandlers() map[string]assessPreviewHandler {
	// Preview handlers mirror assess profile names for command-surface parity.
	return map[string]assessPreviewHandler{
		"adapter-capture":         runAdapterCaptureAssessPreview,
		"managed-harness":         runManagedAssessPreview,
		"forensic-retention":      runForensicAssessPreview,
		"ci-artifact-observation": runCIArtifactAssessPreview,
		"authority-envelope":      runAuthorityAssessPreview,
	}
}

func runAuthorityAssess(opts *flagSet, stdout, stderr io.Writer) int {
	if !requireNamedValues(map[string]string{
		"--out":               opts.stringValue("out"),
		"--authority-package": opts.stringValue("authority-package"),
	}, stderr, "authority-envelope assess") {
		// Authority assessment requires an explicit package because raw prompts
		// or model outputs are outside this CLI's accepted evidence shape.
		return exitUsage
	}
	// Authority evaluation is package-bound; unreadable packages are
	// cannot_verify because the caller supplied an authority artifact path.
	pkg, err := authority.ReadPackage(opts.stringValue("authority-package"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	result := authority.Evaluate(pkg)
	return writeAuthorityAssessment(opts, result, stdout, stderr)
}

func writeAuthorityAssessment(opts *flagSet, result authority.Result, stdout, stderr io.Writer) int {
	if err := authority.Write(opts.stringValue("out"), result); err != nil {
		// Authority results need the package-specific writer to preserve schema
		// shape and deterministic formatting.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeJSONPayloadUnchecked(stdout, result)
	return authorityExitCode(result)
}

type adapterCapturePreviewReport struct {
	Command          string            `json:"command"`
	SelectedProfile  string            `json:"selected_profile"`
	Inputs           map[string]string `json:"inputs"`
	ExpectedEvidence map[string]string `json:"expected_evidence"`
	Safety           map[string]string `json:"safety"`
	NextActions      []string          `json:"next_actions"`
	Claim            string            `json:"claim"`
}

func runAdapterCaptureAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Preview records only whether the run path is present, not whether the run
	// satisfies adapter-capture evidence rules.
	inputs := map[string]string{
		"run": managedInputStatus(opts.stringValue("run")),
	}
	writeJSONPayloadUnchecked(stdout, newAdapterCapturePreviewReport(inputs))
	return previewInputExitCode(inputs)
}

func newAdapterCapturePreviewReport(inputs map[string]string) adapterCapturePreviewReport {
	// Adapter preview lists expected vocabulary and safety constraints only; it
	// does not inspect raw adapter payloads or issue a verdict.
	return adapterCapturePreviewReport{
		Command:          "assess preview",
		SelectedProfile:  adaptercapture.ProfileAdapterCapture,
		Inputs:           inputs,
		ExpectedEvidence: adapterCapturePreviewExpectedEvidence,
		Safety:           adapterCapturePreviewSafety,
		NextActions:      adapterCapturePreviewActions(inputs),
		Claim:            "preview is read-only and does not emit an adapter capture verdict",
	}
}

var adapterCapturePreviewExpectedEvidence = map[string]string{
	"binding_modes":        "same_chain,adapter_bundle",
	"test_provenance":      "ci_executed,wrapper_executed,harness_observed,agent_reported,cannot_verify",
	"capture_depth_states": "captured,missing_telemetry,not_integrated,unsupported,retention_limited,not_assessed,cannot_verify",
}

var adapterCapturePreviewSafety = map[string]string{
	"raw_payloads":    "not_rendered",
	"adapter_secrets": "not_rendered",
	"gateway_refs":    "token_free_refs_only",
	"model_payloads":  "digest_or_block18_reference_only",
}

func runManagedAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Managed preview reports local artifact readiness without evaluating
	// capability, registry, run, or witness state.
	inputs := map[string]string{
		"run":              managedInputStatus(opts.stringValue("run")),
		"adapter_registry": managedInputStatus(opts.stringValue("adapter-registry")),
		"managed_policy":   managedInputStatus(opts.stringValue("managed-policy")),
		"managed_witness":  managedInputStatus(opts.stringValue("managed-witness")),
	}
	writeJSONPayloadUnchecked(stdout, newManagedPreviewReport(inputs))
	return previewInputExitCode(inputs)
}

func newManagedPreviewReport(inputs map[string]string) managedPreviewReport {
	// Managed preview exposes which local artifacts are ready to be assessed,
	// while keeping capability and witness verdicts uncomputed.
	return managedPreviewReport{
		Command:         "assess preview",
		SelectedProfile: managed.ProfileManagedHarness,
		Inputs:          inputs,
		NextActions:     managedPreviewActions(inputs),
		Claim:           "preview is read-only and does not emit a managed verdict",
	}
}

type forensicPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	PolicyEffects   map[string]string `json:"policy_effects"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

func runForensicAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Forensic preview reports only run/policy presence and leaves redaction
	// evaluation to the real assessment command.
	inputs := map[string]string{
		"run":              managedInputStatus(opts.stringValue("run")),
		"redaction_policy": managedInputStatus(opts.stringValue("redaction-policy")),
	}
	writeJSONPayloadUnchecked(stdout, newForensicPreviewReport(inputs))
	return previewInputExitCode(inputs)
}

func newForensicPreviewReport(inputs map[string]string) forensicPreviewReport {
	// Forensic preview documents retention policy effects without executing
	// redaction or exposing matched sensitive values.
	return forensicPreviewReport{
		Command:         "assess preview",
		SelectedProfile: forensic.ProfileForensicRetention,
		Inputs:          inputs,
		PolicyEffects:   forensicPreviewPolicyEffects,
		NextActions:     forensicPreviewActions(inputs),
		Claim:           "preview is read-only and does not emit a forensic verdict",
	}
}

var forensicPreviewPolicyEffects = map[string]string{
	"redaction_engine": "not_executed_in_preview",
	"matched_values":   "not_rendered",
	"rule_refs":        "shown_when_present_in_policy_or_run_metadata",
	"retention_modes":  "digest_only,sanitized_excerpt,encrypted_raw_ref,external_artifact_ref,not_assessed",
}

type ciArtifactPreviewReport struct {
	Command          string            `json:"command"`
	SelectedProfile  string            `json:"selected_profile"`
	Inputs           map[string]string `json:"inputs"`
	ObservedFamilies []string          `json:"observed_families"`
	StateModel       map[string]string `json:"state_model"`
	Safety           map[string]string `json:"safety"`
	NextActions      []string          `json:"next_actions"`
	Claim            string            `json:"claim"`
}

type authorityPreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	Inputs          map[string]string `json:"inputs"`
	StateModel      map[string]string `json:"state_model"`
	Safety          map[string]string `json:"safety"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

func runAuthorityAssessPreview(opts *flagSet, stdout io.Writer) int {
	// Authority preview treats the package path as an input pointer only; it does
	// not evaluate authority rules.
	inputs := map[string]string{
		"authority_package": managedInputStatus(opts.stringValue("authority-package")),
	}
	writeJSONPayloadUnchecked(stdout, newAuthorityPreviewReport(inputs))
	return previewInputExitCode(inputs)
}

func newAuthorityPreviewReport(inputs map[string]string) authorityPreviewReport {
	// Authority preview reports package readiness and state vocabulary without
	// evaluating policy effects or accepting raw prompt/model material.
	return authorityPreviewReport{
		Command:         "assess preview",
		SelectedProfile: authority.Profile,
		Inputs:          inputs,
		StateModel:      authorityPreviewStateModel,
		Safety:          authorityPreviewSafety,
		NextActions:     authorityPreviewActions(inputs),
		Claim:           "preview is read-only and does not emit an authority or policy verdict",
	}
}

var authorityPreviewStateModel = map[string]string{
	"authority":   "within_authority,outside_authority,not_assessed,cannot_verify",
	"attribution": "verified,not_assessed,cannot_verify",
	"binding":     "verified,not_assessed,cannot_verify",
}

var authorityPreviewSafety = map[string]string{
	"raw_prompts":       "not_accepted",
	"raw_model_outputs": "not_accepted",
	"credential_refs":   "rejected_as_malformed",
	"policy_effects":    "not_emitted",
}

func runCIArtifactAssessPreview(opts *flagSet, stdout io.Writer) int {
	// CI artifact preview never reads artifact content or calls remote systems.
	inputs := map[string]string{
		"artifact_manifest": managedInputStatus(opts.stringValue("artifact-manifest")),
	}
	writeJSONPayloadUnchecked(stdout, newCIArtifactPreviewReport(inputs))
	return previewInputExitCode(inputs)
}

func newCIArtifactPreviewReport(inputs map[string]string) ciArtifactPreviewReport {
	// CI artifact preview names required families and safety posture without
	// fetching network artifacts or reading raw artifact content.
	return ciArtifactPreviewReport{
		Command:          "assess preview",
		SelectedProfile:  ciartifact.ProfileCIArtifactObservation,
		Inputs:           inputs,
		ObservedFamilies: ciArtifactPreviewObservedFamilies,
		StateModel:       ciArtifactPreviewStateModel,
		Safety:           ciArtifactPreviewSafety,
		NextActions:      ciArtifactPreviewActions(inputs),
		Claim:            "preview is read-only and does not emit a CI artifact observation verdict",
	}
}

var ciArtifactPreviewObservedFamilies = []string{
	"run", "report", "witness", "provenance", "evidence",
	"trace", "artifact_index", "redaction_scan", "review", "change_ci",
}

var ciArtifactPreviewStateModel = map[string]string{
	"top_level": "pass,fail,cannot_verify,not_assessed",
	"producer":  "ci_uploaded,checked_in,local_generated,agent_reported,harness_observed,external_artifact_ref,not_assessed",
	"access":    "present,absent,partial,expired,inaccessible,malformed,unsafe,not_assessed,cannot_verify",
}

var ciArtifactPreviewSafety = map[string]string{
	"raw_artifact_content": "not_rendered",
	"reason_payloads":      "safe_reason_codes_only",
	"network_fetch":        "not_performed",
}

func previewInputExitCode(inputs map[string]string) int {
	for _, state := range inputs {
		if previewInputCannotVerify(state) {
			// Bad preview inputs block setup confidence without emitting a profile
			// assessment verdict.
			return exitCannotVerify
		}
	}
	return 0
}

func previewInputCannotVerify(state string) bool {
	return state == "present_unreadable" || state == "present_malformed"
}

func runAssessExplain(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "assess explain"}
	opts.setString("assessment-result", "")
	// Explanation renders an existing assessment artifact only; it never
	// re-evaluates evidence or upgrades a verdict.
	path, err := parseAssessExplainArgs(opts, args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	return explainAssessmentResult(path, stdout, stderr)
}

func parseAssessExplainArgs(opts *flagSet, args []string) (string, error) {
	if err := opts.parse(args); err != nil {
		return "", err
	}
	if len(opts.rest()) != 0 {
		return "", errors.New("assess explain accepts only flags")
	}
	path := opts.stringValue("assessment-result")
	if path == "" {
		// A concrete artifact path is required so operators cannot ask the
		// explainer to interpret ad hoc text as evidence.
		return "", errors.New("assess explain requires --assessment-result <file>")
	}
	return path, nil
}

func explainAssessmentResult(path string, stdout, stderr io.Writer) int {
	var envelope struct {
		SchemaVersion   string `json:"schema_version"`
		SelectedProfile string `json:"selected_profile"`
	}
	// Read the minimal schema envelope first; selected_profile is descriptive
	// and must not choose the parser.
	if err := readJSONFile(path, &envelope); err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return dispatchAssessmentExplanation(path, envelope.SchemaVersion, stdout, stderr)
}

func dispatchAssessmentExplanation(path, schemaVersion string, stdout, stderr io.Writer) int {
	handler, ok := assessmentExplainHandlers[schemaVersion]
	if !ok {
		// Unknown schemas remain cannot_verify instead of falling back to a
		// best-effort renderer that could hide a profile mismatch.
		fmt.Fprintf(stderr, "unsupported assessment-result schema_version: %s\n", schemaVersion)
		return exitCannotVerify
	}
	return handler(path, stdout, stderr)
}

type assessmentExplainHandler func(string, io.Writer, io.Writer) int

var assessmentExplainHandlers = map[string]assessmentExplainHandler{
	// Schema versions select typed explainers so spoofed profile names cannot
	// redirect artifact interpretation.
	adaptercapture.SchemaVersion:  explainTypedAssessment[adaptercapture.AssessmentResult](explainAdapterCaptureAssessment),
	managed.SchemaVersion:         explainTypedAssessment[managed.AssessmentResult](explainManagedAssessment),
	forensic.SchemaVersion:        explainTypedAssessment[forensic.AssessmentResult](explainForensicAssessment),
	ciartifact.SchemaVersion:      explainTypedAssessment[ciartifact.ObservationResult](explainCIArtifactObservation),
	authority.ResultSchemaVersion: explainTypedAssessment[authority.Result](explainAuthorityEvaluation),
}

func explainTypedAssessment[T any](explain func(T, io.Writer) int) assessmentExplainHandler {
	return func(path string, stdout, stderr io.Writer) int {
		var result T
		// The typed load is the trust boundary for explanation; renderers only
		// restate fields from that decoded artifact.
		if err := readJSONFile(path, &result); err != nil {
			fmt.Fprintln(stderr, err)
			return exitCannotVerify
		}
		return explain(result, stdout)
	}
}

func explainAdapterCaptureAssessment(result adaptercapture.AssessmentResult, stdout io.Writer) int {
	// Adapter explanations are derived views of stored conditions, reasons,
	// and next actions; the assessment verdict is not recomputed here.
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Adapter capture assessment: %s\n", result.AdapterCaptureAssessment)
	fmt.Fprintf(stdout, "Trust scope: %s\n", result.TrustScope)
	for _, condition := range result.AdapterCaptureConditions {
		explainAdapterCaptureCondition(condition, stdout)
	}
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
	return 0
}

func explainAdapterCaptureCondition(condition adaptercapture.Condition, stdout io.Writer) {
	fmt.Fprintf(stdout, "Adapter condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	if condition.CappedToRetentionMode != "" {
		// Retention caps are printed beside the condition that caused them so
		// the explanation preserves the evidence chain.
		fmt.Fprintf(stdout, "Capped to retention mode: %s\n", condition.CappedToRetentionMode)
	}
}

func explainManagedAssessment(result managed.AssessmentResult, stdout io.Writer) int {
	// Managed explanations keep condition state and remediation text visible
	// without treating setup readability as an independent pass claim.
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Managed harness assessment: %s\n", result.ManagedHarnessAssessment)
	fmt.Fprintf(stdout, "Trust scope: %s\n", result.TrustScope)
	for _, condition := range result.ManagedConditions {
		fmt.Fprintf(stdout, "Managed condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	}
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
	return 0
}

func explainForensicAssessment(result forensic.AssessmentResult, stdout io.Writer) int {
	// Forensic explanations restate retention findings from the artifact so
	// missing raw references stay represented as recorded assessment state.
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Forensic retention assessment: %s\n", result.ForensicRetentionAssessment)
	fmt.Fprintf(stdout, "Trust scope: %s\n", result.TrustScope)
	for _, condition := range result.ForensicConditions {
		explainForensicCondition(condition, stdout)
	}
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
	return 0
}

func explainForensicCondition(condition forensic.Condition, stdout io.Writer) {
	fmt.Fprintf(stdout, "Forensic condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	if condition.CappedToRetentionMode != "" {
		// The cap belongs to this condition, not to the entire artifact.
		fmt.Fprintf(stdout, "Capped to retention mode: %s\n", condition.CappedToRetentionMode)
	}
}

func explainCIArtifactObservation(result ciartifact.ObservationResult, stdout io.Writer) int {
	// CI artifact explanations keep producer, access, binding, and output
	// safety states separate so one gap cannot mask another.
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "CI artifact observation: %s\n", result.ArtifactObservationState)
	// Top-level scopes summarize the overall observation; family rows below
	// retain their own finer-grained state.
	fmt.Fprintf(stdout, "Authority scope: %s\n", result.AuthorityScope)
	fmt.Fprintf(stdout, "Producer scope: %s\n", result.ProducerScope)
	fmt.Fprintf(stdout, "Artifact access state: %s\n", result.ArtifactAccessState)
	explainCIArtifactFamilies(result.ArtifactFamilies, stdout)
	// Index and output-safety summaries are printed after family rows so they
	// cannot be confused with per-family evidence.
	fmt.Fprintf(stdout, "Artifact index: %s (%s)\n", result.ArtifactIndex.Result, result.ArtifactIndex.ReasonCode)
	fmt.Fprintf(stdout, "Output safety: %s (%s)\n", result.OutputSafety.State, result.OutputSafety.ReasonCode)
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
	return 0
}

func explainCIArtifactFamilies(families []ciartifact.FamilyObservation, stdout io.Writer) {
	for _, family := range families {
		// Each family line preserves producer, access, and binding state as
		// separate evidence dimensions.
		fmt.Fprintf(stdout, "Artifact family %s: %s (%s)\n", family.Family, family.FamilyState, family.ReasonCode)
		fmt.Fprintf(stdout, "  Producer scope: %s\n", family.ProducerScope)
		fmt.Fprintf(stdout, "  Artifact access: %s\n", family.ArtifactAccessState)
		fmt.Fprintf(stdout, "  Binding: %s\n", family.BindingState)
	}
}

func explainAuthorityEvaluation(result authority.Result, stdout io.Writer) int {
	// Authority explanations separate action policy results from evidence
	// binding results so provenance failures stay visible.
	fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	fmt.Fprintf(stdout, "Authority evaluation: %s\n", result.AuthorityEvaluationState)
	fmt.Fprintf(stdout, "Selected policy: %s\n", result.SelectedPolicyID)
	explainAuthorityActionEvaluations(result.Evaluations, stdout)
	explainAuthorityBindingEvaluations(result.BindingEvaluations, stdout)
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
	return 0
}

func explainAuthorityActionEvaluations(evaluations []authority.AuthorityEvaluation, stdout io.Writer) {
	for _, eval := range evaluations {
		// Attribution fields are printed independently; a matched action rule
		// does not imply actor, tool, or model provenance is complete.
		fmt.Fprintf(stdout, "Observed action %s: %s (%s)\n", eval.EventID, eval.State, eval.ReasonCode)
		fmt.Fprintf(stdout, "  Actor attribution: %s\n", eval.ActorAttribution)
		fmt.Fprintf(stdout, "  Tool attribution: %s\n", eval.ToolAttribution)
		fmt.Fprintf(stdout, "  Model attribution: %s\n", eval.ModelAttribution)
		if eval.MatchedRuleRef != "" {
			fmt.Fprintf(stdout, "  Matched rule: %s\n", eval.MatchedRuleRef)
		}
	}
}

func explainAuthorityBindingEvaluations(bindings []authority.EvidenceBinding, stdout io.Writer) {
	for _, binding := range bindings {
		// Binding lines keep evidence-binding state separate from action state
		// so provenance gaps remain visible in CLI output.
		fmt.Fprintf(stdout, "Binding %s: %s (%s)\n", binding.BindingID, binding.BindingState, binding.ReasonCode)
	}
}

func managedInputStatus(path string) string {
	if strings.TrimSpace(path) == "" {
		return "absent"
	}
	info, err := os.Stat(path)
	if err != nil {
		// Preview status intentionally avoids leaking filesystem error details.
		return "present_unreadable"
	}
	if info.IsDir() {
		// Run-directory inputs are assessed through their normalized run.json.
		return jsonReadableStatus(filepath.Join(path, "run.json"))
	}
	return jsonReadableStatus(path)
}

func jsonReadableStatus(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		// Readability is enough for preview; the real assessment reports the
		// concrete parse/load error.
		return "present_unreadable"
	}
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		// Malformed JSON blocks setup without interpreting partial contents.
		return "present_malformed"
	}
	return "present_readable"
}

func managedPreviewActions(inputs map[string]string) []string {
	// Managed preview actions name setup gaps for each required evidence class
	// without deriving a managed-harness verdict.
	return previewActionsForInputs(
		inputs,
		[]string{"run", "adapter_registry", "managed_policy", "managed_witness"},
		"Supply %s before managed assessment.",
		"Fix %s so it is readable JSON or a run directory.",
	)
}

func forensicPreviewActions(inputs map[string]string) []string {
	// Forensic preview actions stay limited to setup remediation; the retention
	// assessment itself still requires a full policy/run evaluation.
	return previewActionsForInputs(
		inputs,
		[]string{"run", "redaction_policy"},
		"Supply %s before forensic retention assessment.",
		"Fix %s so it is readable JSON or a run directory.",
	)
}

func previewActionsForInputs(inputs map[string]string, order []string, missingMessage, invalidMessage string) []string {
	var actions []string
	for _, key := range order {
		// Fixed key order keeps preview remediation stable for docs/tests.
		switch inputs[key] {
		case "absent":
			// Missing setup evidence is actionable before any assessment runs.
			actions = append(actions, fmt.Sprintf(missingMessage, key))
		case "present_unreadable", "present_malformed":
			// Unreadable or malformed setup evidence must be replaced, not
			// interpreted as a negative assessment condition.
			actions = append(actions, fmt.Sprintf(invalidMessage, key))
		}
	}
	return actions
}

func adapterCapturePreviewActions(inputs map[string]string) []string {
	return previewActionForInputState(inputs, "run", "Supply run before adapter capture assessment.", "Fix run so it is a readable JSON run directory.")
}

func ciArtifactPreviewActions(inputs map[string]string) []string {
	return previewActionForInputState(inputs, "artifact_manifest", "Supply artifact manifest before CI artifact observation assessment.", "Fix artifact manifest so it is readable JSON.")
}

func authorityPreviewActions(inputs map[string]string) []string {
	return previewActionForInputState(inputs, "authority_package", "Supply authority package before authority envelope assessment.", "Fix authority package so it is readable JSON.")
}

func previewActionForInputState(inputs map[string]string, key, absentAction, malformedAction string) []string {
	// Preview action helpers translate setup state into remediation text only;
	// they do not evaluate the underlying assessment payload.
	switch inputs[key] {
	case "absent":
		return []string{absentAction}
	case "present_unreadable", "present_malformed":
		// Unreadable and malformed inputs share the same repair path because
		// both block assessment before verdict logic can run.
		return []string{malformedAction}
	default:
		return nil
	}
}

func adapterCaptureExitCode(result adaptercapture.AssessmentResult) int {
	return stringExitCode(result.AdapterCaptureAssessment, adapterCaptureExitCodes, exitCannotVerify)
}

var adapterCaptureExitCodes = map[string]int{
	adaptercapture.StatePass: 0,
	adaptercapture.StateFail: 1,
}

func managedExitCode(result managed.AssessmentResult) int {
	return stringExitCode(result.ManagedHarnessAssessment, managedExitCodes, exitCannotVerify)
}

var managedExitCodes = map[string]int{
	managed.StatePass: 0,
	managed.StateFail: 1,
}

func forensicExitCode(result forensic.AssessmentResult) int {
	return stringExitCode(result.ForensicRetentionAssessment, forensicExitCodes, exitCannotVerify)
}

var forensicExitCodes = map[string]int{
	forensic.StatePass: 0,
	forensic.StateFail: 1,
}

func ciArtifactExitCode(result ciartifact.ObservationResult) int {
	return stringExitCode(result.ArtifactObservationState, ciArtifactExitCodes, exitCannotVerify)
}

var ciArtifactExitCodes = map[string]int{
	ciartifact.StatePass: 0,
	ciartifact.StateFail: 1,
}

var authorityExitCodes = map[string]int{
	authority.StateWithinAuthority:  0,
	authority.StateOutsideAuthority: 1,
	authority.StateNotAssessed:      exitCannotVerify,
}

func authorityExitCode(result authority.Result) int {
	return stringExitCode(result.AuthorityEvaluationState, authorityExitCodes, exitCannotVerify)
}

func runCheckpoint(_ context.Context, args []string, stdout, stderr io.Writer) int {
	// Checkpoint subcommands are explicit because create signs local run state
	// while verify only replays an existing signed checkpoint.
	command, rest, ok := checkpointCommand(args, stderr)
	if !ok {
		return exitUsage
	}
	return checkpointCommandHandlers[command](rest, stdout, stderr)
}

var checkpointCommandHandlers = map[string]subcommandHandler{
	"create": runCheckpointCreate,
	"verify": runCheckpointVerify,
}

func checkpointCommand(args []string, stderr io.Writer) (string, []string, bool) {
	if len(args) == 0 {
		// A checkpoint command without a verb cannot decide whether it is
		// creating signing evidence or replaying it.
		fmt.Fprintln(stderr, "checkpoint requires create or verify")
		return "", nil, false
	}
	// Keep the remainder untouched so each subcommand owns its flag contract.
	return knownCheckpointCommand(args[0], args[1:], stderr)
}

func knownCheckpointCommand(command string, rest []string, stderr io.Writer) (string, []string, bool) {
	if command != "create" && command != "verify" {
		// Unknown checkpoint verbs are usage errors, not verifier states.
		fmt.Fprintln(stderr, "checkpoint requires create or verify")
		return "", nil, false
	}
	return command, rest, true
}

func runCheckpointCreate(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseCheckpointCreateArgs(args, stderr)
	if !ok {
		return code
	}
	// Creation persists the signed artifact before printing its id; stdout is
	// only a convenience pointer, not the checkpoint proof.
	created, code, ok := createCheckpointFromOptions(opts, stderr)
	if !ok {
		return code
	}
	fmt.Fprintf(stdout, "checkpoint: %s\n", created.CheckpointID)
	return 0
}

func createCheckpointFromOptions(opts *flagSet, stderr io.Writer) (checkpoint.SignedCheckpoint, int, bool) {
	var key checkpoint.KeyPair
	// The private key file is local signing material, not an authority proof by
	// itself; policy binding is checked later during verification.
	if err := readJSONFile(opts.stringValue("private-key"), &key); err != nil {
		fmt.Fprintln(stderr, err)
		return checkpoint.SignedCheckpoint{}, 1, false
	}
	return createAndWriteCheckpoint(opts, key, stderr)
}

func createAndWriteCheckpoint(opts *flagSet, key checkpoint.KeyPair, stderr io.Writer) (checkpoint.SignedCheckpoint, int, bool) {
	// checkpoint.Create binds the requested run directory and signer identity
	// into the signed payload before this CLI writes the JSON artifact.
	// The CLI does not infer signer authority here; policy binding is replayed
	// by checkpoint verification.
	created, err := checkpoint.Create(opts.stringValue("run"), checkpoint.CreateOptions{
		CheckpointID: opts.stringValue("id"),
		SignerID:     opts.stringValue("signer-id"),
		Key:          key,
	})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return checkpoint.SignedCheckpoint{}, 1, false
	}
	if err := writeJSONFile(opts.stringValue("out"), created); err != nil {
		// A checkpoint that cannot be persisted is not usable evidence.
		fmt.Fprintln(stderr, err)
		return checkpoint.SignedCheckpoint{}, 1, false
	}
	// Return the created artifact only after it exists at the requested path.
	return created, 0, true
}

func parseCheckpointCreateArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := newCheckpointCreateFlagSet()
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Required create flags provide all source, sink, and signer inputs for a
	// replayable checkpoint artifact.
	if !requireOnlyFlags(opts, stderr, "checkpoint create accepts only flags", checkpointCreateRequiredFlags) {
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func newCheckpointCreateFlagSet() *flagSet {
	opts := &flagSet{name: "checkpoint create"}
	for _, flag := range checkpointCreateStringFlags {
		// Registration order is fixed so help/tests observe the same command
		// contract while defaults stay beside their flag names.
		opts.setString(flag.name, flag.defaultValue)
	}
	return opts
}

var checkpointCreateStringFlags = []struct {
	name         string
	defaultValue string
}{
	{"run", ""},
	{"out", ""},
	{"private-key", ""},
	{"signer-id", ""},
	{"id", "checkpoint-001"},
}

var checkpointCreateRequiredFlags = []requiredCLIFlag{
	{"run", "checkpoint create requires --run"},
	{"out", "checkpoint create requires --out"},
	{"private-key", "checkpoint create requires --private-key"},
	{"signer-id", "checkpoint create requires --signer-id"},
}

func runCheckpointVerify(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseCheckpointVerifyArgs(args, stderr)
	if !ok {
		return code
	}
	// Verification reads immutable inputs first, then delegates all replay and
	// signer-policy decisions to the checkpoint package.
	// The CLI only renders the resulting verification envelope and exit state.
	signed, policy, code, ok := loadCheckpointVerifyInputs(opts, stderr)
	if !ok {
		return code
	}
	result := checkpoint.Verify(opts.stringValue("run"), signed, policy)
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return checkpointVerifyExitCode(result.Result)
}

func parseCheckpointVerifyArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := newCheckpointVerifyFlagSet()
	if err := opts.parse(args); err != nil {
		// Parse errors happen before any signed checkpoint or policy is loaded.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if len(opts.rest()) != 0 {
		// Positional arguments would make the verification source ambiguous.
		fmt.Fprintln(stderr, "checkpoint verify accepts only flags")
		return nil, exitUsage, false
	}
	if err := requireCheckpointVerifyInputs(opts); err != nil {
		// Required verification inputs identify the run replay source and signed
		// checkpoint artifact before optional policy can affect trust scope.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func newCheckpointVerifyFlagSet() *flagSet {
	opts := &flagSet{name: "checkpoint verify"}
	// Verification names the run, signed checkpoint, and optional trust policy
	// as separate inputs so each can fail independently.
	for _, flag := range checkpointVerifyStringFlags {
		opts.setString(flag, "")
	}
	return opts
}

var checkpointVerifyStringFlags = []string{"run", "checkpoint", "policy"}

func loadCheckpointVerifyInputs(opts *flagSet, stderr io.Writer) (checkpoint.SignedCheckpoint, *checkpoint.TrustedCheckpointPolicy, int, bool) {
	var signed checkpoint.SignedCheckpoint
	// Decode the signed checkpoint before policy so malformed proof artifacts
	// fail before optional trust policy is considered.
	if err := readJSONFile(opts.stringValue("checkpoint"), &signed); err != nil {
		fmt.Fprintln(stderr, err)
		return checkpoint.SignedCheckpoint{}, nil, 1, false
	}
	policy, err := readCheckpointPolicy(opts.stringValue("policy"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return checkpoint.SignedCheckpoint{}, nil, 1, false
	}
	return signed, policy, 0, true
}

func readCheckpointPolicy(path string) (*checkpoint.TrustedCheckpointPolicy, error) {
	if path == "" {
		// An absent policy leaves signer trust to the checkpoint verifier's
		// default local semantics; it is not treated as a green CI policy.
		return nil, nil
	}
	var loaded checkpoint.TrustedCheckpointPolicy
	if err := readJSONFile(path, &loaded); err != nil {
		return nil, err
	}
	return &loaded, nil
}

func requireCheckpointVerifyInputs(opts *flagSet) error {
	if strings.TrimSpace(opts.stringValue("run")) == "" {
		// The run directory is the source replay target for the signed payload.
		return fmt.Errorf("checkpoint verify requires --run")
	}
	if strings.TrimSpace(opts.stringValue("checkpoint")) == "" {
		// The signed checkpoint artifact is mandatory verification input.
		return fmt.Errorf("checkpoint verify requires --checkpoint")
	}
	return nil
}

func checkpointVerifyExitCode(state string) int {
	if state == checkpoint.StatePass {
		// Only an explicit checkpoint pass maps to shell success.
		return 0
	}
	if state == checkpoint.StateCannotVerify {
		return exitCannotVerify
	}
	return 1
}

func runReport(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts, target, code, ok := parseReportArgs(args, stderr)
	if !ok {
		return code
	}
	// Reports regenerate summary artifacts from run evidence; stdout mirrors
	// the summary but the output directory is the durable review surface.
	artifacts, err := demo.WriteReport(target, opts.stringValue("out"), opts.stringValue("contract"))
	if err != nil {
		// Report generation failure means no durable summary can be trusted.
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(artifacts.Summary, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return reportExitCode(artifacts.Summary)
}

func reportExitCode(summary demo.Summary) int {
	if summary.CannotVerifyCount > 0 {
		// Cannot-verify rows take precedence because missing evidence is not a
		// successful report even when no explicit failures were counted.
		return exitCannotVerify
	}
	if summary.FailedCount > 0 {
		return 1
	}
	return 0
}

func parseReportArgs(args []string, stderr io.Writer) (*flagSet, string, int, bool) {
	opts := &flagSet{name: "report"}
	// Reports have one durable output root and an optional contract override;
	// the observed run target remains positional for command readability.
	opts.setString("out", "")
	opts.setString("contract", "")
	// Parse flag names before inspecting the target so malformed options are
	// not collapsed into an evidence-root error.
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, "", exitUsage, false
	}
	// A report is useful only when one source and one durable output sink are
	// both explicit in the command line.
	target, ok := reportTargetArg(opts, stderr)
	if !ok {
		return nil, "", exitUsage, false
	}
	if !requireStringFlag(opts, stderr, "out", "report requires --out <dir>") {
		return nil, "", exitUsage, false
	}
	return opts, target, 0, true
}

func reportTargetArg(opts *flagSet, stderr io.Writer) (string, bool) {
	targets := opts.rest()
	if len(targets) != 1 {
		// A single target keeps report provenance bound to one run root.
		fmt.Fprintln(stderr, "report requires <runs-root-or-run-dir>")
		return "", false
	}
	return targets[0], true
}

func runGate(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if code, ok := runGateSubcommand(args, stdout, stderr); ok {
		// Gate subcommands are read-only/explanatory paths and do not evaluate a
		// new gate result.
		return code
	}
	opts, target, outPath, code, ok := parseGateArgs(args, stderr)
	if !ok {
		return code
	}
	if opts.stringValue("profile") == demo.GateProfileProtected {
		// Protected mode requires checkpoint, policy, and witness inputs in
		// addition to local run rows.
		return runProtectedGate(target, outPath, opts, stdout, stderr)
	}
	return runStandardGate(target, outPath, opts, stdout, stderr)
}

func runStandardGate(target, outPath string, opts *flagSet, stdout, stderr io.Writer) int {
	// Standard gates derive local/CI/audit verdicts from demo run evidence and
	// optional witness data; they do not grant protected-profile trust.
	result, err := demo.WriteGate(target, outPath, opts.stringValue("contract"), opts.stringValue("witness"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return gateExitCode(result)
}

func runGateSubcommand(args []string, stdout, stderr io.Writer) (int, bool) {
	return runOptionalSubcommand(args, stdout, stderr, gateSubcommandHandlers)
}

var gateSubcommandHandlers = map[string]subcommandHandler{
	"explain": runGateExplain,
	"preview": runGatePreview,
}

func parseGateArgs(args []string, stderr io.Writer) (*flagSet, string, string, int, bool) {
	// Protected-profile inputs are parsed alongside standard gate inputs so one
	// command surface can expose both local and protected modes.
	opts := newStringFlagSet("gate", gateStringFlags)
	if err := opts.parse(args); err != nil {
		// Parse failures happen before any report rows or protected inputs are
		// read.
		fmt.Fprintln(stderr, err)
		return nil, "", "", exitUsage, false
	}
	// Target and output validation stay separate so diagnostics distinguish the
	// evidence source from the artifact destination.
	target, ok := gateTargetArg(opts, stderr)
	if !ok {
		return nil, "", "", exitUsage, false
	}
	outPath, ok := gateOutputPath(opts, stderr)
	if !ok {
		return nil, "", "", exitUsage, false
	}
	return opts, target, outPath, 0, true
}

func gateTargetArg(opts *flagSet, stderr io.Writer) (string, bool) {
	targets := opts.rest()
	if len(targets) == 1 {
		return targets[0], true
	}
	// Gate evaluation is bound to exactly one run root or run directory.
	fmt.Fprintln(stderr, "gate requires <runs-root-or-run-dir>")
	return "", false
}

func gateOutputPath(opts *flagSet, stderr io.Writer) (string, bool) {
	// The output path is validated after target arity so diagnostics first
	// establish which evidence source the gate would evaluate.
	outPath := opts.stringValue("out")
	if outPath != "" {
		return outPath, true
	}
	// Persisted gate JSON is the artifact later explain/preview commands can
	// inspect; stdout is only a rendered copy.
	fmt.Fprintln(stderr, "gate requires --out <file>")
	return "", false
}

var gateStringFlags = []string{"out", "contract", "witness", "profile", "checkpoint", "checkpoint-policy"}

func runProtectedGate(target, outPath string, opts *flagSet, stdout, stderr io.Writer) int {
	// Protected gate resolution is separated from writing so input failures do
	// not create a partial gate artifact.
	result, code := resolveProtectedGate(target, opts, stderr)
	if code != 0 {
		return code
	}
	return writeProtectedGateResult(outPath, result, stdout, stderr)
}

func resolveProtectedGate(target string, opts *flagSet, stderr io.Writer) (demo.GateResult, int) {
	// Read external trust inputs before loading rows so missing checkpoint,
	// policy, or witness evidence is reported as a gate setup error.
	inputs, code, ok := readProtectedGateInputs(opts, stderr)
	if !ok {
		return demo.GateResult{}, code
	}
	// Rows and contract are loaded after setup evidence so protected gate
	// failures cannot hide missing external authority inputs.
	contract, rows, runDir, code, ok := loadProtectedGateRows(target, opts.stringValue("contract"), stderr)
	if !ok {
		return demo.GateResult{}, code
	}
	// The expected witness is derived from the protected run itself; supplied
	// witness files must match it rather than define their own expectation.
	expected, code, ok := loadProtectedWitnessExpectation(target, stderr)
	if !ok {
		return demo.GateResult{}, code
	}
	checkpointResult := verifiedProtectedCheckpoint(runDir, inputs, expected)
	return evaluateProtectedGate(rows, contract, checkpointResult, inputs.witness, expected), 0
}

func evaluateProtectedGate(rows []demo.RunRow, contract trace.Contract, checkpointResult checkpoint.VerificationResult, witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) demo.GateResult {
	// Protected evaluation receives already-bound checkpoint and witness facts;
	// the CLI does not override package verdicts.
	return demo.EvaluateProtectedGate(rows, contract, demo.ProtectedGateInput{
		Checkpoint:         checkpointResult,
		PolicyProvided:     true,
		Witness:            &witnessSummary,
		WitnessExpectation: expected,
		Now:                time.Now().UTC(),
	})
}

func verifiedProtectedCheckpoint(runDir string, inputs protectedGateInputs, expected demo.WitnessExpectation) checkpoint.VerificationResult {
	// Checkpoint replay runs before protected aggregation so signature/policy
	// facts stay separate from gate conditions.
	return protectedCheckpointVerification(
		checkpoint.Verify(runDir, inputs.signed, &inputs.policy),
		inputs.signed,
		inputs.policy,
		inputs.witness,
		expected,
	)
}

func writeProtectedGateResult(path string, result demo.GateResult, stdout, stderr io.Writer) int {
	if err := writeJSONFile(path, result); err != nil {
		// A protected verdict that cannot be written is not reviewable evidence.
		fmt.Fprintln(stderr, err)
		return 1
	}
	writeIndentedPayload(stdout, result)
	return gateExitCode(result)
}

type protectedGateInputs struct {
	signed  checkpoint.SignedCheckpoint
	policy  checkpoint.TrustedCheckpointPolicy
	witness demo.WitnessSummary
}

func readProtectedGateInputs(opts *flagSet, stderr io.Writer) (protectedGateInputs, int, bool) {
	var inputs protectedGateInputs
	// Signed checkpoint data is the replay proof for the selected run.
	if code, ok := readRequiredProtectedInput("--checkpoint", opts.stringValue("checkpoint"), &inputs.signed, stderr); !ok {
		return protectedGateInputs{}, code, false
	}
	// Checkpoint policy pins the accepted signer authority.
	if code, ok := readRequiredProtectedInput("--checkpoint-policy", opts.stringValue("checkpoint-policy"), &inputs.policy, stderr); !ok {
		return protectedGateInputs{}, code, false
	}
	// Witness summary binds the protected run to observed CI evidence.
	if code, ok := readRequiredProtectedInput("--witness", opts.stringValue("witness"), &inputs.witness, stderr); !ok {
		return protectedGateInputs{}, code, false
	}
	return inputs, 0, true
}

func readRequiredProtectedInput(flag, path string, value any, stderr io.Writer) (int, bool) {
	if strings.TrimSpace(path) == "" {
		// Protected mode has no implicit local defaults for external trust
		// inputs.
		fmt.Fprintf(stderr, "protected gate requires %s\n", flag)
		return exitUsage, false
	}
	// All protected inputs are decoded as JSON artifacts before evaluation so
	// the gate never accepts unchecked path strings as trust evidence.
	if err := readJSONFile(path, value); err != nil {
		// Malformed trust inputs are usage/setup failures, not a green local gate
		// with omitted protected evidence.
		fmt.Fprintln(stderr, err)
		return exitUsage, false
	}
	return 0, true
}

func loadProtectedGateRows(target, contractPath string, stderr io.Writer) (trace.Contract, []demo.RunRow, string, int, bool) {
	// Protected rows are always checked against the requested contract; this
	// path has no silent default upgrade after the profile is selected.
	contract, err := trace.LoadContract(contractPath)
	if err != nil {
		// Protected mode cannot fall back to an implicit contract after the user
		// supplied one.
		fmt.Fprintln(stderr, err)
		return trace.Contract{}, nil, "", 1, false
	}
	rows, err := demo.VerifiedRows(target, contract)
	if err != nil {
		// Row replay failures block protected evaluation before checkpoint facts
		// can be joined.
		fmt.Fprintln(stderr, err)
		return trace.Contract{}, nil, "", 1, false
	}
	// Protected checkpoints replay a single concrete run directory, not only a
	// report-level collection of rows.
	runDir, err := protectedRunDir(target)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return trace.Contract{}, nil, "", exitCannotVerify, false
	}
	return contract, rows, runDir, 0, true
}

func loadProtectedWitnessExpectation(target string, stderr io.Writer) (demo.WitnessExpectation, int, bool) {
	// The expectation loader derives the run id and artifact digests that the
	// supplied witness summary must bind to.
	expected, err := demoWitnessExpectation(target)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return demo.WitnessExpectation{}, exitCannotVerify, false
	}
	return expected, 0, true
}

func runGateExplain(args []string, stdout, stderr io.Writer) int {
	path, code, ok := parseGateExplainArgs(args, stderr)
	if !ok {
		return code
	}
	result, code, ok := readGateExplainResult(path, stderr)
	if !ok {
		return code
	}
	// Explanation is read-only: it restates a persisted gate result without
	// re-running gate evaluation.
	explainGateResult(result, stdout)
	return 0
}

func parseGateExplainArgs(args []string, stderr io.Writer) (string, int, bool) {
	opts := &flagSet{name: "gate explain"}
	// Explanation is keyed by one existing gate-result artifact.
	opts.setString("gate-result", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return "", exitUsage, false
	}
	if len(opts.rest()) != 0 {
		// Gate explanations accept only an artifact path to avoid mixing
		// positional prose with result evidence.
		// Positional text would not be replayable as a source-bound result.
		fmt.Fprintln(stderr, "gate explain accepts only flags")
		return "", exitUsage, false
	}
	path := opts.stringValue("gate-result")
	if path == "" {
		// A persisted gate result is required because explanation does not
		// synthesize verdicts from loose fields.
		fmt.Fprintln(stderr, "gate explain requires --gate-result <file>")
		return "", exitUsage, false
	}
	// The caller reads and validates the artifact before rendering.
	return path, 0, true
}

func readGateExplainResult(path string, stderr io.Writer) (demo.GateResult, int, bool) {
	var result demo.GateResult
	if err := readJSONFile(path, &result); err != nil {
		// Missing gate artifacts are cannot_verify for explanation, not usage.
		fmt.Fprintln(stderr, err)
		return demo.GateResult{}, exitCannotVerify, false
	}
	if result.SchemaVersion != demo.GateSchemaVersion && result.SchemaVersion != demo.GateSchemaVersionBlock16 {
		// Unsupported result schemas remain cannot_verify instead of being
		// rendered with stale field assumptions.
		fmt.Fprintf(stderr, "unsupported gate-result schema_version: %s\n", result.SchemaVersion)
		return demo.GateResult{}, exitCannotVerify, false
	}
	return result, 0, true
}

func explainGateResult(result demo.GateResult, stdout io.Writer) {
	if result.SchemaVersion == demo.GateSchemaVersion {
		// Legacy gate results have no protected-profile fields; state that
		// absence explicitly rather than implying not_assessed conditions.
		fmt.Fprintln(stdout, "Protected profile fields: absent")
	}
	explainGateSummary(result, stdout)
	explainProtectedGateDetails(result, stdout)
	explainGateCollections(result, stdout)
}

func explainGateSummary(result demo.GateResult, stdout io.Writer) {
	// Summary lines keep the layered gate states distinct so local, CI witness,
	// audit, and protected outcomes are not collapsed into one score.
	fmt.Fprintf(stdout, "Gate mode: %s\n", result.GateMode)
	fmt.Fprintf(stdout, "Trust cap: %s\n", result.TrustCap)
	if result.SelectedProfile != "" {
		fmt.Fprintf(stdout, "Selected profile: %s\n", result.SelectedProfile)
	}
	fmt.Fprintf(stdout, "Local gate: %s\n", result.LocalGate)
	fmt.Fprintf(stdout, "CI witness gate: %s\n", result.CIWitnessGate)
	fmt.Fprintf(stdout, "Audit-grade gate: %s\n", result.AuditGradeGate)
	if result.ProtectedGate != "" {
		fmt.Fprintf(stdout, "Protected gate: %s\n", result.ProtectedGate)
	}
}

func explainProtectedGateDetails(result demo.GateResult, stdout io.Writer) {
	if result.CheckpointVerification != nil {
		// Checkpoint verification is shown separately from protected conditions
		// because replay/signature failure and policy failure are different
		// evidence gaps.
		fmt.Fprintf(stdout, "Checkpoint result: %s\n", result.CheckpointVerification.Result)
		fmt.Fprintf(stdout, "Checkpoint trust scope: %s\n", result.CheckpointVerification.TrustScope)
	}
	for _, condition := range result.ProtectedConditions {
		// Protected conditions remain individual rows for auditability.
		fmt.Fprintf(stdout, "Protected condition %s: %s (%s)\n", condition.ID, condition.State, condition.ReasonCode)
	}
}

func explainGateCollections(result demo.GateResult, stdout io.Writer) {
	// Collection explainers preserve the original evidence categories and
	// remediation lists from the persisted gate result.
	explainRequiredRuns(result.RequiredRuns, stdout)
	explainWitnessBindings(result.WitnessBindings, stdout)
	explainMissingAuditEvidence(result.MissingAuditEvidence, stdout)
	explainOverrideRequests(result.OverrideRequests, stdout)
	explainReasons(result.Reasons, stdout)
	explainNextActions(result.NextActions, stdout)
}

func explainRequiredRuns(requiredRuns []demo.RequiredRunResult, stdout io.Writer) {
	for _, requiredRun := range requiredRuns {
		// One stable line per required run keeps the human explanation auditable
		// without inventing a separate summary verdict.
		fmt.Fprintf(stdout, "Required run %s: %s\n", requiredRun.ID, requiredRun.State)
	}
}

func explainWitnessBindings(bindings []demo.WitnessBinding, stdout io.Writer) {
	for _, binding := range bindings {
		// Binding lines expose the witness-to-run link directly instead of
		// hiding provenance under a combined health score.
		fmt.Fprintf(stdout, "Witness binding %s: %s\n", binding.ID, binding.State)
	}
}

func explainMissingAuditEvidence(missingEvidence []string, stdout io.Writer) {
	for _, missing := range missingEvidence {
		// Missing audit evidence stays visible as a concrete gap; explanation
		// output must not collapse it into a green summary.
		fmt.Fprintf(stdout, "Missing audit evidence: %s\n", missing)
	}
}

func explainOverrideRequests(overrides []demo.OverrideRequest, stdout io.Writer) {
	for _, override := range overrides {
		// Override requests remain separate records because each one needs its
		// own evidence-backed state.
		fmt.Fprintf(stdout, "Override %s: %s\n", override.OverrideID, override.State)
	}
}

func explainReasons(reasons []string, stdout io.Writer) {
	for _, reason := range reasons {
		// Reasons are emitted verbatim as traceable verdict support, not as
		// prose-only interpretation.
		fmt.Fprintf(stdout, "Reason: %s\n", reason)
	}
}

func explainNextActions(actions []string, stdout io.Writer) {
	for _, action := range actions {
		// Next actions are remediation hints and do not upgrade the current
		// assessed state.
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
}

type gatePreviewReport struct {
	Command            string   `json:"command"`
	GateMode           string   `json:"gate_mode"`
	TrustCap           string   `json:"trust_cap"`
	RequiredRuns       []string `json:"required_runs"`
	RequiredEvidence   []string `json:"required_evidence"`
	WitnessInspectable bool     `json:"witness_inspectable"`
	WitnessMismatches  []string `json:"witness_mismatches,omitempty"`
	Claim              string   `json:"claim"`
}

type protectedGatePreviewReport struct {
	Command         string            `json:"command"`
	SelectedProfile string            `json:"selected_profile"`
	TrustCap        string            `json:"trust_cap"`
	Inputs          map[string]string `json:"inputs"`
	NextActions     []string          `json:"next_actions"`
	Claim           string            `json:"claim"`
}

func runGatePreview(args []string, stdout, stderr io.Writer) int {
	opts, targets, code, ok := parseGatePreviewArgs(args, stderr)
	if !ok {
		return code
	}
	if opts.stringValue("profile") == demo.GateProfileProtected {
		// Protected preview reports input readiness only; it never evaluates a
		// protected verdict.
		return runProtectedGatePreview(opts, stdout)
	}
	// Standard preview reads the contract to describe required evidence, not to
	// inspect run artifacts or predict pass/fail.
	contract, err := trace.LoadContract(opts.stringValue("contract"))
	if err != nil {
		// A preview with an unreadable contract cannot describe the expected
		// evidence surface.
		fmt.Fprintln(stderr, err)
		return 1
	}
	// Standard preview emits required evidence only; it never inspects run rows.
	report := buildGatePreviewReport(contract, opts.stringValue("witness"), targets[0])
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func parseGatePreviewArgs(args []string, stderr io.Writer) (*flagSet, []string, int, bool) {
	// Preview accepts both standard and protected flags because it reports setup
	// readiness without committing to a verdict mode.
	opts := newStringFlagSet("gate preview", gatePreviewStringFlags)
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, nil, exitUsage, false
	}
	targets := opts.rest()
	if len(targets) != 1 {
		// Preview is still target-scoped so witness-binding checks, when
		// requested, compare against one run root.
		fmt.Fprintln(stderr, "gate preview requires <runs-root-or-run-dir>")
		return nil, nil, exitUsage, false
	}
	return opts, targets, 0, true
}

var gatePreviewStringFlags = []string{"contract", "witness", "profile", "checkpoint", "checkpoint-policy"}

func buildGatePreviewReport(contract trace.Contract, witnessPath, target string) gatePreviewReport {
	// The preview report is a planning artifact: it names required runs and
	// evidence IDs without claiming the gate will pass.
	report := gatePreviewReport{
		Command:          "gate preview",
		GateMode:         previewGateMode(contract),
		TrustCap:         string(trace.TrustScopeLocalObserved),
		RequiredRuns:     requiredRunIDs(contract),
		RequiredEvidence: requiredEvidenceIDsForCLI(contract),
		Claim:            "preview is read-only and does not claim the gate will pass",
	}
	if witnessPath != "" {
		// Optional witness inspection checks binding shape only; it does not
		// produce a CI-witness gate verdict.
		report.WitnessInspectable, report.WitnessMismatches = demo.PreviewWitnessBinding(witnessPath, target)
	}
	return report
}

func runProtectedGatePreview(opts *flagSet, stdout io.Writer) int {
	// Protected preview classifies required input files before any checkpoint
	// replay or witness trust evaluation.
	inputs := map[string]string{
		"checkpoint":        protectedInputStatus(opts.stringValue("checkpoint")),
		"checkpoint_policy": protectedInputStatus(opts.stringValue("checkpoint-policy")),
		"witness":           protectedInputStatus(opts.stringValue("witness")),
	}
	report := newProtectedGatePreviewReport(inputs)
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	for _, state := range inputs {
		if state == "present_unreadable" || state == "present_malformed" {
			// Bad setup inputs keep preview cannot_verify; they are not lowered
			// into protected gate failures because no verdict was evaluated.
			return exitCannotVerify
		}
	}
	return 0
}

func newProtectedGatePreviewReport(inputs map[string]string) protectedGatePreviewReport {
	// The report deliberately mirrors the gate command vocabulary while keeping
	// the claim explicit that no protected verdict was produced.
	return protectedGatePreviewReport{
		Command:         "gate preview",
		SelectedProfile: demo.GateProfileProtected,
		TrustCap:        string(trace.TrustScopeLocalObserved),
		Inputs:          inputs,
		NextActions:     protectedPreviewActions(inputs),
		Claim:           "preview is read-only and does not emit a protected verdict",
	}
}

func protectedRunDir(target string) (string, error) {
	runDirs, err := demo.DiscoverRunDirs(target)
	if err != nil {
		return "", err
	}
	if len(runDirs) != 1 {
		// Protected replay requires exactly one run so checkpoint payloads,
		// witness expectation, and observed rows all bind to the same source.
		return "", fmt.Errorf("protected gate requires one selected run, got %d", len(runDirs))
	}
	return runDirs[0], nil
}

func protectedCheckpointVerification(result checkpoint.VerificationResult, signed checkpoint.SignedCheckpoint, policy checkpoint.TrustedCheckpointPolicy, witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) checkpoint.VerificationResult {
	if !canGrantProtectedCheckpointTrust(result, signed, policy, witnessSummary, expected) {
		return result
	}
	// A protected upgrade is allowed only after checkpoint replay, signer
	// policy, and witness binding all agree on the same run.
	result.SignerAuthorityState = checkpoint.StatePass
	result.TrustScope = checkpoint.TrustScopeCISigned
	result.Result = checkpoint.StatePass
	return result
}

func canGrantProtectedCheckpointTrust(result checkpoint.VerificationResult, signed checkpoint.SignedCheckpoint, policy checkpoint.TrustedCheckpointPolicy, witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	// Do not upgrade explicit checkpoint failures; protected trust can only
	// tighten cannot-verify/local states when all external bindings match.
	return result.Result != checkpoint.StateFail &&
		signed.Signer.Authority == checkpoint.AuthorityCIIsolatedJob &&
		policyAllowsSigner(policy, signed) &&
		witnessMatchesProtectedInput(witnessSummary, expected)
}

func policyAllowsSigner(policy checkpoint.TrustedCheckpointPolicy, signed checkpoint.SignedCheckpoint) bool {
	for _, signer := range policy.AllowedSigners {
		// Policy must bind signer id, authority class, and public key to the
		// checkpoint signature.
		if signerMatchesCheckpoint(signer, signed) {
			return true
		}
	}
	return false
}

func signerMatchesCheckpoint(signer checkpoint.TrustedSigner, signed checkpoint.SignedCheckpoint) bool {
	return signer.SignerID == signed.Signer.SignerID && signer.Authority == signed.Signer.Authority && signer.PublicKey == signed.Signature.PublicKey
}

func witnessMatchesProtectedInput(witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	if !witnessHasProtectedTrust(witnessSummary) || !witnessSourceMatches(witnessSummary, expected) {
		// Witness status/source mismatch blocks protected trust before artifact
		// digest comparison.
		return false
	}
	return witnessArtifactsMatch(witnessSummary.RunArtifacts, expected.RunArtifacts)
}

func witnessArtifactsMatch(runArtifacts, expectedRunArtifacts []demo.WitnessArtifactDigest) bool {
	// Artifact counts must match exactly so extra or missing witness artifacts
	// cannot hide under a partial digest match.
	expectedArtifacts := expectedArtifactDigests(expectedRunArtifacts)
	if len(runArtifacts) != len(expectedArtifacts) {
		return false
	}
	for _, artifact := range runArtifacts {
		// Every reported witness artifact must match the expected path and
		// digest derived from the observed run.
		if !witnessArtifactMatchesExpectation(artifact, expectedArtifacts) {
			return false
		}
	}
	return true
}

func expectedArtifactDigests(expectedRunArtifacts []demo.WitnessArtifactDigest) map[string]string {
	expectedArtifacts := map[string]string{}
	for _, artifact := range expectedRunArtifacts {
		// The map is keyed by normalized artifact path to make digest matching
		// deterministic and independent of input ordering.
		expectedArtifacts[artifact.Path] = artifact.SHA256
	}
	return expectedArtifacts
}

func witnessArtifactMatchesExpectation(artifact demo.WitnessArtifactDigest, expectedArtifacts map[string]string) bool {
	expectedSHA, ok := expectedArtifacts[artifact.Path]
	return ok && expectedSHA == artifact.SHA256
}

func witnessHasProtectedTrust(witnessSummary demo.WitnessSummary) bool {
	return witnessSummary.Kind == "github-actions" && witnessSummary.Status == demo.GatePass && witnessSummary.TrustScope == "ci_witnessed"
}

func witnessSourceMatches(witnessSummary demo.WitnessSummary, expected demo.WitnessExpectation) bool {
	// Empty expected source fields are intentionally wildcards for portable
	// examples; non-empty fields must match the witness identity exactly.
	return optionalStringMatches(expected.Repository, witnessSummary.Source.Repository) &&
		optionalStringMatches(expected.Ref, witnessSummary.Source.Ref) &&
		optionalStringMatches(expected.CommitSHA, witnessSummary.Source.CommitSHA) &&
		optionalStringMatches(expected.RunID, witnessSummary.CIIdentity.RunID)
}

func optionalStringMatches(expected, actual string) bool {
	return expected == "" || actual == expected
}

func demoWitnessExpectation(target string) (demo.WitnessExpectation, error) {
	// Expectations are derived from observed run artifacts, not from the
	// supplied witness summary.
	runDirs, err := demo.DiscoverRunDirs(target)
	if err != nil {
		return demo.WitnessExpectation{}, err
	}
	runID, artifacts, err := demoWitnessArtifacts(runDirs)
	if err != nil {
		return demo.WitnessExpectation{}, err
	}
	return demo.WitnessExpectation{RunID: runID, RunArtifacts: artifacts}, nil
}

func demoWitnessArtifacts(runDirs []string) (string, []demo.WitnessArtifactDigest, error) {
	artifacts := make([]demo.WitnessArtifactDigest, 0, len(runDirs))
	runID := ""
	for _, runDir := range runDirs {
		// Each discovered run contributes the digest for its retained run.json
		// artifact.
		// The artifact helper keeps path reading and hash calculation together.
		artifactRunID, digest, err := demoWitnessArtifact(runDir)
		if err != nil {
			return "", nil, err
		}
		if runID == "" {
			// The first run artifact anchors the demo witness expectation; later
			// artifacts contribute digests without changing the expected run ID.
			runID = artifactRunID
		}
		artifacts = append(artifacts, demo.WitnessArtifactDigest{
			Path:   filepath.ToSlash(filepath.Join(filepath.Base(runDir), "run.json")),
			SHA256: digest,
		})
	}
	// The completed digest list remains local evidence until witness evaluation.
	return runID, artifacts, nil
}

func demoWitnessArtifact(runDir string) (string, string, error) {
	// The retained run artifact supplies the run id while the file bytes supply
	// the digest that witness summaries must match.
	artifact, err := trace.OpenRunArtifact(runDir)
	if err != nil {
		return "", "", err
	}
	digest, err := sha256File(runDir, "run.json")
	if err != nil {
		return "", "", err
	}
	return artifact.Manifest.RunID, digest, nil
}

func sha256File(dir, name string) (string, error) {
	// Digest calculation reads the artifact exactly as retained on disk.
	data, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func protectedInputStatus(path string) string {
	if strings.TrimSpace(path) == "" {
		// Missing protected inputs are reported as setup gaps, not as verdicts.
		return "absent"
	}
	var value any
	if err := readJSONFile(path, &value); err != nil {
		return protectedInputErrorStatus(err)
	}
	return "present_readable"
}

func protectedInputErrorStatus(err error) string {
	if os.IsNotExist(err) || errors.Is(err, os.ErrPermission) {
		// Protected preview distinguishes unavailable inputs from malformed JSON
		// so users know whether to fix access or content.
		return "present_unreadable"
	}
	return "present_malformed"
}

func protectedPreviewActions(inputs map[string]string) []string {
	names := []string{"checkpoint", "checkpoint_policy", "witness"}
	actions := make([]string, 0)
	for _, name := range names {
		// Fixed input order keeps preview remediation stable for docs/tests.
		switch inputs[name] {
		case "absent":
			actions = append(actions, fmt.Sprintf("Supply %s input before running protected gate.", name))
		case "present_unreadable", "present_malformed":
			// Both unreadable and malformed artifacts require replacement before
			// protected replay can make an evidence-backed claim.
			actions = append(actions, fmt.Sprintf("Replace %s input with readable JSON.", name))
		}
	}
	return actions
}

func runOverride(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseOverrideRequestArgs(args, stderr)
	if !ok {
		return code
	}
	// Appending is the only state-changing step; parsing alone never creates an
	// override artifact.
	event, err := appendOverrideRequestEvent(opts)
	if err != nil {
		// Append failure means no override request was durably recorded.
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "override_event: %s\n", event.EventID)
	return 0
}

func appendOverrideRequestEvent(opts *flagSet) (trace.Event, error) {
	return trace.AppendRunEvent(opts.stringValue("out"), trace.EventPolicyOverrideRequested, overrideRequestPayload(opts), "sdp-trace-cli")
}

func overrideRequestPayload(opts *flagSet) map[string]any {
	// Override requests are appended as trace events; they request policy
	// review and never upgrade an existing gate verdict by themselves.
	// The event payload mirrors trace fields so the appended event is the
	// authoritative override request artifact.
	payload := map[string]any{
		"override_id":  opts.stringValue("id"),
		"producer":     "sdp-trace-cli",
		"origin":       "native_cli",
		"requested_by": opts.stringValue("by"),
		"reason":       opts.stringValue("reason"),
		"source_ref":   opts.stringValue("source-ref"),
		"scope":        opts.stringValue("scope"),
		"created_at":   time.Now().UTC().Format(time.RFC3339Nano),
	}
	if external := opts.stringValue("external-reference"); external != "" {
		// External references remain metadata until another verifier resolves
		// them; the CLI does not treat them as approval proof.
		payload["external_reference"] = external
	}
	return payload
}

func parseOverrideRequestArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	if !isOverrideRequest(args) {
		// The override namespace currently has one explicit write action.
		fmt.Fprintln(stderr, "override requires request")
		return nil, exitUsage, false
	}
	return parseOverrideRequestFlags(args[1:], stderr)
}

func isOverrideRequest(args []string) bool {
	return len(args) != 0 && args[0] == "request"
}

func parseOverrideRequestFlags(args []string, stderr io.Writer) (*flagSet, int, bool) {
	// Override requests are write operations, so parsing must establish a fully
	// named payload before any trace event can be appended.
	opts := &flagSet{name: "override request"}
	// Each flag maps directly to a persisted trace payload key, keeping the
	// event reviewable without positional inference.
	// Requiredness is checked after parsing so diagnostics can distinguish
	// unknown flags from missing trace fields.
	opts.setString("out", "")
	opts.setString("id", "")
	opts.setString("by", "")
	opts.setString("reason", "")
	opts.setString("source-ref", "")
	opts.setString("scope", "")
	opts.setString("external-reference", "")
	if err := opts.parse(args); err != nil {
		// Parser errors stop before the command can create partial override
		// evidence.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// No positional text is accepted because all persisted override request
	// fields must have stable JSON keys.
	if len(opts.rest()) != 0 {
		// Free-form positional text would make the override reason ambiguous.
		fmt.Fprintln(stderr, "override request accepts only flags")
		return nil, exitUsage, false
	}
	return requireOverrideRequestFlags(opts, stderr)
}

func requireOverrideRequestFlags(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	// Required field validation happens before the run directory is opened for
	// append, preventing partial override events.
	if !requireRequiredFlags(opts, stderr, overrideRequestRequiredFlags) {
		// Required fields identify who asked, what scope is affected, and which
		// source evidence the override references.
		return nil, exitUsage, false
	}
	return opts, 0, true
}

var overrideRequestRequiredFlags = []requiredCLIFlag{
	{"out", "override request requires --out"},
	{"id", "override request requires --id"},
	{"by", "override request requires --by"},
	{"reason", "override request requires --reason"},
	{"source-ref", "override request requires --source-ref"},
	{"scope", "override request requires --scope"},
}

func readJSONFile(path string, dst any) error {
	// Shared JSON reads are local artifact loads; callers decide whether a
	// failure is usage, cannot_verify, or ordinary command failure.
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func writeJSONFile(path string, value any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Pretty JSON keeps generated evidence reviewable and stable in fixtures.
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func writeTextFileAtomic(path, value string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	// Text artifacts are written through a sibling temp file so readers never
	// observe a partially rendered Markdown/report file.
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(path)+".*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	return finishAtomicTextWrite(tmp, tmpName, path, value)
}

func finishAtomicTextWrite(tmp *os.File, tmpName, path, value string) error {
	if err := writeAndCloseTempText(tmp, value); err != nil {
		return err
	}
	// Permissions are normalized before rename so the final artifact has the
	// same readable mode as other generated evidence files.
	if err := os.Chmod(tmpName, 0o644); err != nil {
		return err
	}
	// Rename is the publication step for the completed text artifact.
	return os.Rename(tmpName, path)
}

func writeAndCloseTempText(tmp *os.File, value string) error {
	if _, err := tmp.WriteString(value); err != nil {
		// Close on write failure so the temp file handle is not leaked before
		// caller cleanup removes it.
		_ = tmp.Close()
		return err
	}
	return tmp.Close()
}

func previewGateMode(contract trace.Contract) string {
	mode := demo.GateModeObservation
	for _, required := range contract.RequiredRuns {
		switch required.Profile {
		case demo.GateModeProtectedFuture:
			// Protected-future requirements dominate advisory CI requirements in
			// preview because they imply a stricter future gate path.
			return demo.GateModeProtectedFuture
		case demo.GateModeAdvisoryCI:
			// Advisory CI is retained unless a protected-future requirement is
			// found later.
			mode = demo.GateModeAdvisoryCI
		}
	}
	return mode
}

func requiredRunIDs(contract trace.Contract) []string {
	ids := make([]string, 0, len(contract.RequiredRuns))
	for _, required := range contract.RequiredRuns {
		if required.ID != "" {
			// Empty IDs are omitted from CLI preview output rather than rendered
			// as ambiguous evidence handles.
			ids = append(ids, required.ID)
		}
	}
	return ids
}

func requiredEvidenceIDsForCLI(contract trace.Contract) []string {
	ids := make([]string, 0, len(contract.RequiredEvidence))
	for _, requirement := range contract.RequiredEvidence {
		if requirement.ID != "" {
			// Preview exposes stable evidence identifiers only.
			ids = append(ids, requirement.ID)
		}
	}
	return ids
}

func gateExitCode(result demo.GateResult) int {
	if code, ok := protectedGateExitCode(result); ok {
		// Protected profiles own the process exit because local-only gate state
		// can otherwise overstate release readiness.
		return code
	}
	return gateStateExitCode(gateExitStates(result))
}

func protectedGateExitCode(result demo.GateResult) (int, bool) {
	if result.SelectedProfile != demo.GateProfileProtected {
		// Non-protected gates fall back to layered local/CI/audit state.
		return 0, false
	}
	code, ok := protectedGateExitCodes[result.ProtectedGate]
	if !ok {
		return 0, false
	}
	return code, true
}

var protectedGateExitCodes = map[string]int{
	demo.GatePass:         0,
	demo.GateFail:         1,
	demo.GateCannotVerify: exitCannotVerify,
	demo.GateNotAssessed:  exitCannotVerify,
}

func gateExitStates(result demo.GateResult) []string {
	states := []string{result.LocalGate, result.CIWitnessGate, result.AuditGradeGate}
	for _, requiredRun := range result.RequiredRuns {
		// Required run states participate in the process exit because missing
		// required evidence should fail even if aggregate fields are stale.
		states = append(states, requiredRun.State)
	}
	return states
}

func gateStateExitCode(states []string) int {
	if hasGateState(states, demo.GateFail, demo.GateMissingTelemetry) {
		// Explicit failure and missing telemetry are shell failures.
		return 1
	}
	if hasGateState(states, demo.GateCannotVerify) {
		// Cannot-verify remains distinct from ordinary failure for automation.
		return exitCannotVerify
	}
	return 0
}

func hasGateState(states []string, targets ...string) bool {
	for _, state := range states {
		// Match against the closed state vocabulary selected by the caller.
		if slices.Contains(targets, state) {
			return true
		}
	}
	return false
}

func runWitness(_ context.Context, args []string, stdout, stderr io.Writer) int {
	options, ok := parseWitnessOptions(args, stderr)
	if !ok {
		return exitUsage
	}
	// Witness output is generated from explicit CLI inputs; missing trust
	// material is rejected before record construction.
	record, err := buildWitnessRecord(options)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return writeWitnessRecordOutput(stdout, record)
}

type witnessOptions struct {
	kind                      string
	target                    string
	out                       string
	reportDir                 string
	witnessEnvelope           string
	customerPKIAuthorityPath  string
	customerPKIPublicCertPath string
	customerPKIPublicKeyPath  string
	customerPKIPayloadDigest  string
	customerPKIFreshnessPath  string
}

func parseWitnessOptions(args []string, stderr io.Writer) (witnessOptions, bool) {
	opts, ok := parseWitnessFlagSet(args, stderr)
	if !ok {
		return witnessOptions{}, false
	}
	// Validation is separated from flag parsing so missing required values can
	// return product-specific messages.
	options, message, ok := witnessOptionsFromFlags(opts)
	if !ok {
		fmt.Fprintln(stderr, message)
		return witnessOptions{}, false
	}
	return options, true
}

func parseWitnessFlagSet(args []string, stderr io.Writer) (*flagSet, bool) {
	opts := &flagSet{name: "witness"}
	// Witness flags cover both generic CI witnesses and customer-PKI material;
	// semantic validation happens after parsing.
	opts.setString("kind", "")
	opts.setString("out", "")
	opts.setString("report-dir", "")
	opts.setString("witness-envelope", "")
	opts.setString("customer-pki-authority-policy", "")
	opts.setString("customer-pki-public-cert", "")
	opts.setString("customer-pki-public-key", "")
	opts.setString("customer-pki-payload-digest", "")
	opts.setString("customer-pki-freshness-evidence", "")
	if err := opts.parse(args); err != nil {
		// Malformed witness flags stop before any CI or Customer PKI material is
		// read from disk.
		fmt.Fprintln(stderr, err)
		return nil, false
	}
	// The witness command has no positional target; target comes from flags so
	// generated records have explicit provenance fields.
	return opts, true
}

func witnessOptionsFromFlags(opts *flagSet) (witnessOptions, string, bool) {
	// Required fields are normalized before optional witness-specific material
	// is copied into the final options struct.
	fields, message, ok := witnessRequiredFieldsFromFlags(opts)
	if !ok {
		return witnessOptions{}, message, false
	}
	return witnessOptionsFromRequiredFields(fields, opts), "", true
}

type witnessRequiredFields struct {
	target string
	kind   string
	out    string
}

func witnessRequiredFieldsFromFlags(opts *flagSet) (witnessRequiredFields, string, bool) {
	// Target is validated before kind-specific flags because every witness
	// record must bind to one observed run root.
	target, message, ok := witnessTargetFromFlags(opts)
	if !ok {
		return witnessRequiredFields{}, message, false
	}
	return witnessKindOutFromFlags(opts, target)
}

func witnessKindOutFromFlags(opts *flagSet, target string) (witnessRequiredFields, string, bool) {
	kind, message, ok := witnessKindFromFlags(opts)
	if !ok {
		return witnessRequiredFields{}, message, false
	}
	// Output path is required for every kind so the witness record is persisted
	// before stdout renders it.
	out, message, ok := witnessOutFromFlags(opts)
	if !ok {
		return witnessRequiredFields{}, message, false
	}
	// Kind-specific validation prevents customer-PKI witnesses from being
	// created without custody/freshness evidence.
	if message, ok := validateWitnessKindFlags(kind, opts); !ok {
		return witnessRequiredFields{}, message, false
	}
	return witnessRequiredFields{target: target, kind: kind, out: out}, "", true
}

func witnessOptionsFromRequiredFields(fields witnessRequiredFields, opts *flagSet) witnessOptions {
	// Optional fields are copied verbatim; witness package validation decides
	// whether a specific profile can use or must reject them.
	// Required fields stay first so the resulting record always has its core
	// target, kind, and output provenance.
	return witnessOptions{
		kind:                      fields.kind,
		target:                    fields.target,
		out:                       fields.out,
		reportDir:                 opts.stringValue("report-dir"),
		witnessEnvelope:           opts.stringValue("witness-envelope"),
		customerPKIAuthorityPath:  opts.stringValue("customer-pki-authority-policy"),
		customerPKIPublicCertPath: opts.stringValue("customer-pki-public-cert"),
		customerPKIPublicKeyPath:  opts.stringValue("customer-pki-public-key"),
		customerPKIPayloadDigest:  opts.stringValue("customer-pki-payload-digest"),
		customerPKIFreshnessPath:  opts.stringValue("customer-pki-freshness-evidence"),
	}
}

func witnessTargetFromFlags(opts *flagSet) (string, string, bool) {
	targets := opts.rest()
	if len(targets) != 1 {
		// One target keeps witness provenance tied to a single run root.
		return "", "witness requires <runs-root-or-run-dir>", false
	}
	return targets[0], "", true
}

func witnessKindFromFlags(opts *flagSet) (string, string, bool) {
	kind := opts.stringValue("kind")
	if !allowedWitnessKind(kind) {
		// The allowed kind list is closed so CLI output maps to known witness
		// schema semantics.
		return "", "witness requires --kind github-actions, gitlab-ci, buildkite, or customer-pki", false
	}
	return kind, "", true
}

func witnessOutFromFlags(opts *flagSet) (string, string, bool) {
	out := opts.stringValue("out")
	if out == "" {
		// Persisted witness JSON is the authority; stdout is only a rendered
		// copy for the caller.
		return "", "witness requires --out <file>", false
	}
	return out, "", true
}

func validateWitnessKindFlags(kind string, opts *flagSet) (string, bool) {
	// Missing kind-specific material is a usage failure, not a generated
	// not_assessed witness record.
	missing := missingWitnessKindFlags(kind, opts)
	if len(missing) > 0 {
		return fmt.Sprintf("customer-pki witness requires %s", strings.Join(missing, ", ")), false
	}
	return "", true
}

func missingWitnessKindFlags(kind string, opts *flagSet) []string {
	if kind != witness.KindCustomerPKI {
		// Non-customer-PKI witnesses do not require customer key custody inputs.
		return nil
	}
	return missingCustomerPKIFlags(opts)
}

func buildWitnessRecord(opts witnessOptions) (witness.Record, error) {
	builder, ok := witnessRecordBuilders()[opts.kind]
	if !ok {
		// This should be unreachable after option validation; keep the error so
		// direct helper misuse cannot silently produce a generic witness.
		return witness.Record{}, fmt.Errorf("unsupported witness kind %q", opts.kind)
	}
	return builder(opts)
}

type witnessRecordBuilder func(witnessOptions) (witness.Record, error)

func witnessRecordBuilders() map[string]witnessRecordBuilder {
	// Builders map closed witness kinds to the package function that owns their
	// evidence interpretation.
	return map[string]witnessRecordBuilder{
		witness.KindGitHubActions: buildGitHubActionsWitness,
		witness.KindGitLabCI:      buildEnvelopeWitness,
		witness.KindBuildkite:     buildEnvelopeWitness,
		witness.KindCustomerPKI:   buildCustomerPKIWitness,
	}
}

func buildGitHubActionsWitness(opts witnessOptions) (witness.Record, error) {
	return witness.WriteGitHubActions(opts.out, opts.target, opts.reportDir, witness.EnvironmentFromOS())
}

func buildEnvelopeWitness(opts witnessOptions) (witness.Record, error) {
	// Envelope-backed profiles all use the supplied envelope as their portable
	// provenance input.
	return witness.WriteProfile(opts.kind, opts.out, opts.target, opts.reportDir, witness.ProfileOptions{
		EnvelopePath: opts.witnessEnvelope,
	})
}

func buildCustomerPKIWitness(opts witnessOptions) (witness.Record, error) {
	// Customer-PKI witnesses require explicit authority, public credential,
	// payload digest, and freshness evidence paths.
	return witness.WriteProfile(witness.KindCustomerPKI, opts.out, opts.target, opts.reportDir, witness.ProfileOptions{
		CustomerPKIAuthorityPolicy: opts.customerPKIAuthorityPath,
		CustomerPKIPublicCert:      opts.customerPKIPublicCertPath,
		CustomerPKIPublicKey:       opts.customerPKIPublicKeyPath,
		CustomerPKIPayloadDigest:   opts.customerPKIPayloadDigest,
		CustomerPKIFreshness:       opts.customerPKIFreshnessPath,
	})
}

func writeWitnessRecordOutput(stdout io.Writer, record witness.Record) int {
	// The witness package has already written the record; stdout repeats the
	// generated artifact in a human-reviewable form.
	payload, _ := json.MarshalIndent(record, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	switch record.Status {
	case witness.StatusCannotVerify, witness.StatusNotAssessed:
		// Automation needs cannot_verify/not_assessed separated from ordinary
		// failed witness checks.
		return exitCannotVerify
	case witness.StatusFail:
		return 1
	default:
		return 0
	}
}

func missingCustomerPKIFlags(opts *flagSet) []string {
	missing := []string{}
	// These fields establish authority, payload binding, and freshness; all are
	// required for customer-PKI witness construction.
	missing = appendMissingStringFlags(missing, opts, map[string]string{
		"customer-pki-authority-policy":   "--customer-pki-authority-policy",
		"customer-pki-payload-digest":     "--customer-pki-payload-digest",
		"customer-pki-freshness-evidence": "--customer-pki-freshness-evidence",
	})
	if missingCustomerPKIPublicCredential(opts) {
		// Either a certificate or raw public key is enough for the public
		// credential anchor.
		missing = append(missing, "--customer-pki-public-cert or --customer-pki-public-key")
	}
	// Sorted output keeps remediation deterministic.
	sort.Strings(missing)
	return missing
}

func appendMissingStringFlags(missing []string, opts *flagSet, required map[string]string) []string {
	for name, flag := range required {
		if strings.TrimSpace(opts.stringValue(name)) == "" {
			// Preserve the user-facing flag spelling in remediation output.
			missing = append(missing, flag)
		}
	}
	return missing
}

func missingCustomerPKIPublicCredential(opts *flagSet) bool {
	return strings.TrimSpace(opts.stringValue("customer-pki-public-cert")) == "" && strings.TrimSpace(opts.stringValue("customer-pki-public-key")) == ""
}

func allowedWitnessKind(kind string) bool {
	switch kind {
	case witness.KindGitHubActions, witness.KindGitLabCI, witness.KindBuildkite, witness.KindCustomerPKI:
		// Each allowed kind has an explicit builder and schema contract.
		return true
	default:
		return false
	}
}

func runWrap(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	// Legacy wrap keeps parsing separate from recorder execution so malformed
	// wrapper metadata cannot create partial run artifacts.
	opts, command, code, ok := parseWrapArgs(args, stderr)
	if !ok {
		return code
	}
	return runLegacyWrapRecorder(ctx, opts, command, stdout, stderr)
}

func runLegacyWrapRecorder(ctx context.Context, opts *flagSet, command []string, stdout, stderr io.Writer) int {
	// recorder.Run owns artifact creation and event sequencing for wrapped
	// commands.
	res, err := recorder.Run(ctx, recorder.RecorderOptions{
		WrapperName:        opts.stringValue("name"),
		ContractPath:       opts.stringValue("contract"),
		UseDefaultContract: true,
		OutputDir:          opts.stringValue("output-dir"),
		Command:            command,
	})
	return writeRunResult(res, err, stdout, stderr)
}

func parseWrapArgs(args []string, stderr io.Writer) (*flagSet, []string, int, bool) {
	// Legacy wrap records a command with the default contract unless the caller
	// supplies a contract path.
	opts := &flagSet{name: "wrap"}
	opts.setString("name", "")
	opts.setString("contract", "")
	opts.setString("output-dir", "")
	if err := opts.parse(args); err != nil {
		return nil, nil, exitUsage, false
	}
	// The legacy wrapper still requires an explicit child command; flags only
	// describe recorder metadata.
	command, ok := wrapCommand(opts, stderr)
	if !ok {
		return nil, nil, exitUsage, false
	}
	return opts, command, 0, true
}

func wrapCommand(opts *flagSet, stderr io.Writer) ([]string, bool) {
	command := opts.rest()
	if len(command) == 0 {
		// A wrap without a child command would create an empty provenance shell.
		fmt.Fprintln(stderr, "wrap requires a command")
		return nil, false
	}
	return command, true
}

func writeRunResult(res recorder.RecorderResult, err error, stdout, stderr io.Writer) int {
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	// The run directory is the durable artifact root; the child exit code is
	// preserved for shell automation.
	fmt.Fprintf(stdout, "run_dir: %s\n", res.RunDir)
	return res.ExitCode
}

func runWrappedCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts, command, code, ok := parseWrappedCommandArgs(args, stderr)
	if !ok {
		return code
	}
	// The modern run command requires an explicit task and contract choice
	// before recorder execution.
	return runTaskRecorder(ctx, opts, command, stdout, stderr)
}

func runTaskRecorder(ctx context.Context, opts *flagSet, command []string, stdout, stderr io.Writer) int {
	useDefault := opts.boolValue("use-default-contract")
	// The recorder is the only layer that writes run manifests and trace events.
	// CLI flags are passed through as explicit recorder options so the retained
	// manifest can explain the task, wrapper, contract, and command sources.
	res, err := recorder.Run(ctx, recorder.RecorderOptions{
		Task:               opts.stringValue("task"),
		WrapperName:        opts.stringValue("name"),
		ContractPath:       opts.stringValue("contract"),
		UseDefaultContract: useDefault,
		OutputDir:          opts.stringValue("output-dir"),
		Command:            command,
	})
	return writeRunResult(res, err, stdout, stderr)
}

func parseWrappedCommandArgs(args []string, stderr io.Writer) (*flagSet, []string, int, bool) {
	opts := &flagSet{name: "run"}
	// The modern run command keeps task identity and contract selection explicit
	// before the observed command payload.
	opts.setString("task", "")
	opts.setString("contract", "")
	opts.setBool("use-default-contract", false)
	opts.setString("name", "")
	opts.setString("output-dir", "")
	if err := opts.parse(args); err != nil {
		// Run flag parsing fails before command execution, so no partial trace is
		// created for malformed recorder options.
		return nil, nil, exitUsage, false
	}
	command := opts.rest()
	// The remaining argv is the command under observation; flags after this
	// point belong to the child process.
	if !requireWrappedCommandArgs(opts, command, stderr) {
		return nil, nil, exitUsage, false
	}
	return opts, command, 0, true
}

func requireWrappedCommandArgs(opts *flagSet, command []string, stderr io.Writer) bool {
	if len(command) == 0 {
		// Recorder runs require an observed command to produce meaningful trace
		// evidence.
		fmt.Fprintln(stderr, "run requires a command")
		return false
	}
	if opts.stringValue("task") == "" {
		// Task id binds the run to the SpecKit task vocabulary.
		fmt.Fprintln(stderr, "run requires --task")
		return false
	}
	if missingRequiredContract(opts) {
		// Contract choice must be explicit unless the caller opts into the
		// built-in default.
		fmt.Fprintln(stderr, "run requires --contract unless --use-default-contract is set")
		return false
	}
	return true
}

func missingRequiredContract(opts *flagSet) bool {
	return opts.stringValue("contract") == "" && !opts.boolValue("use-default-contract")
}

func runDryRun(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runPreviewCommand("dry-run", "simulation", args, stdout, stderr)
}

func runPreview(_ context.Context, args []string, stdout, stderr io.Writer) int {
	return runPreviewCommand("preview", "preview", args, stdout, stderr)
}

func runPreviewCommand(commandName, mode string, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	// Preview identifies the hypothetical command and contract before rendering
	// a non-evidence plan.
	opts, command, code, ok := parsePreviewCommandArgs(commandName, args, stderr)
	if !ok {
		return code
	}
	// Preview loads the contract but deliberately avoids recorder.Run so no run
	// artifacts or trace events are written.
	// Contract parsing is the only validation preview performs.
	contract, code, ok := loadPreviewContract(commandName, opts, stderr)
	if !ok {
		return code
	}
	payload := previewCommandPayload(mode, command, contract)
	writePreviewCommandPayload(stdout, payload)
	return 0
}

func writePreviewCommandPayload(stdout io.Writer, payload map[string]any) {
	// Preview output is a declarative plan, not evidence that the command ran.
	data, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintf(stdout, "%s\n", data)
}

func previewCommandPayload(mode string, command []string, contract trace.Contract) map[string]any {
	return map[string]any{
		// The descriptor is the only child-process representation; preview does
		// not execute or retain stdout/stderr.
		"mode":                 mode,
		"command_descriptor":   trace.NewCommandDescriptor(command),
		"contract":             contract,
		"boundaries":           previewBoundaries(),
		"offline_implications": previewOfflineImplications(),
		"writes_artifacts":     false,
		"safe_retention_modes": safeRetentionModes(),
		"warning":              "no run artifacts were written",
	}
}

func parsePreviewCommandArgs(commandName string, args []string, stderr io.Writer) (*flagSet, []string, int, bool) {
	opts := &flagSet{name: commandName}
	// Preview defaults to the built-in contract unless the caller supplies an
	// override, but the output still records which contract was selected.
	opts.setString("contract", "")
	opts.setBool("use-default-contract", true)
	opts.setString("name", "")
	if err := opts.parse(args); err != nil {
		return nil, nil, exitUsage, false
	}
	// Preview command payload is required because the output describes the
	// command that would be observed.
	command := opts.rest()
	if len(command) == 0 {
		// Preview still needs a command descriptor even though it will not run.
		fmt.Fprintf(stderr, "%s requires a command\n", commandName)
		return nil, nil, exitUsage, false
	}
	if missingRequiredContract(opts) {
		// Default contract use is explicit in preview output so dry-run reports
		// remain replayable.
		fmt.Fprintf(stderr, "%s requires --contract unless --use-default-contract is set\n", commandName)
		return nil, nil, exitUsage, false
	}
	// Successful parsing only produces a plan; recorder execution is not
	// reachable from this command path.
	// The command slice is retained only as a descriptor input for the preview
	// payload.
	return opts, command, 0, true
}

func loadPreviewContract(commandName string, opts *flagSet, stderr io.Writer) (trace.Contract, int, bool) {
	contractPath := opts.stringValue("contract")
	contract := trace.DefaultContract
	if contractPath != "" {
		// A malformed preview contract is cannot_verify because the preview
		// cannot describe valid evidence requirements.
		loaded, err := trace.LoadContract(contractPath)
		if err != nil {
			fmt.Fprintf(stderr, "failed to load contract: %v\n", err)
			return trace.Contract{}, exitCannotVerify, false
		}
		contract = loaded
	}
	return contract, 0, true
}

func runDoctor(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts, code, ok := parseDoctorArgs(args, stderr)
	if !ok {
		return code
	}
	if strings.TrimSpace(opts.stringValue("profile")) != "" {
		// Profile mode changes the evidence surface from local defaults to
		// repository installation/proof diagnostics.
		// Profile-scoped doctor delegates to repoobserver because those checks
		// inspect repository installation/proof state, not local run defaults.
		return runRepoObserverDoctor(opts, stdout, stderr)
	}
	return runLocalDoctor(opts, stdout, stderr)
}

func parseDoctorArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "doctor"}
	// Doctor defaults inspect local run/report readiness; a profile switches to
	// repository-observer diagnostics.
	opts.setString("contract", "")
	opts.setString("output-dir", defaultRunRoot)
	opts.setString("report-dir", defaultReportDir)
	opts.setString("profile", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		// Doctor parsing is read-only; malformed flags never trigger filesystem
		// diagnostics or repo-observer checks.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if len(opts.rest()) != 0 {
		// Doctor diagnostics are selected only by flags so output remains
		// deterministic and profile-scoped.
		fmt.Fprintln(stderr, "doctor does not accept positional arguments")
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func runRepoObserverDoctor(opts *flagSet, stdout, stderr io.Writer) int {
	if opts.stringValue("profile") != repoobserver.ProfileGithubActionsGitHooksV1 {
		// The CLI exposes only the portable GitHub Actions/git-hooks profile.
		fmt.Fprintf(stderr, "doctor --profile requires %s\n", repoobserver.ProfileGithubActionsGitHooksV1)
		return exitUsage
	}
	// Doctor is read-only; it reports install/proof state without modifying the
	// repository.
	status, err := repoobserver.Doctor(repoobserver.Options{Profile: opts.stringValue("profile")})
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	return writeRepoObserverDoctor(opts, status, stdout, stderr)
}

func writeRepoObserverDoctor(opts *flagSet, status repoobserver.Status, stdout, stderr io.Writer) int {
	if err := repoobserver.WriteJSON(opts.stringValue("out"), status); err != nil {
		// Persisted doctor JSON is the machine-readable diagnostic artifact.
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprint(stdout, repoobserver.HumanTable(status))
	return repoObserverExitCode(status)
}

func runLocalDoctor(opts *flagSet, stdout, stderr io.Writer) int {
	// Local doctor checks CLI defaults and environment-derived witness context;
	// it does not inspect repository install hooks.
	report, exitCode := buildDoctorReport(doctorOptions{
		ContractPath: opts.stringValue("contract"),
		OutputDir:    opts.stringValue("output-dir"),
		ReportDir:    opts.stringValue("report-dir"),
		Env:          witness.EnvironmentFromOS(),
	})
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		// A doctor report that cannot be serialized cannot be trusted by
		// automation even if the underlying checks ran.
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "%s\n", data)
	return exitCode
}

func runInstall(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseInstallRepoObserverArgs(args, stdout, stderr)
	if !ok {
		return code
	}
	// Install returns status for both preview and write modes; JSON is written
	// before error handling so failed attempts are still inspectable.
	// The repoobserver package owns preview/write semantics; the CLI guarantees
	// a status artifact is attempted first.
	status, err := repoobserver.Install(repoObserverOptionsFromFlags(opts))
	if writeErr := repoobserver.WriteJSON(opts.stringValue("out"), status); writeErr != nil {
		// Install status JSON is the durable diagnostic surface for both preview
		// and write modes.
		fmt.Fprintln(stderr, writeErr)
		return 1
	}
	if code, handled := handleRepoObserverInstallError(status, err, stdout, stderr); handled {
		return code
	}
	fmt.Fprint(stdout, repoobserver.HumanTable(status))
	return repoObserverInstallExitCode(opts.boolValue("write"), status)
}

func handleRepoObserverInstallError(status repoobserver.Status, err error, stdout, stderr io.Writer) (int, bool) {
	if err == nil {
		return 0, false
	}
	if status.SchemaVersion != "" {
		// Partial status with a schema is still useful human evidence.
		fmt.Fprint(stdout, repoobserver.HumanTable(status))
	}
	fmt.Fprintln(stderr, err)
	return exitCannotVerify, true
}

func repoObserverInstallExitCode(write bool, status repoobserver.Status) int {
	if !write {
		// Preview mode reports planned changes but does not fail on an uninstalled
		// repository state.
		return 0
	}
	return repoObserverExitCode(status)
}

func parseInstallRepoObserverArgs(args []string, stdout, stderr io.Writer) (*flagSet, int, bool) {
	if isHelp(args) {
		printUsage(stdout)
		return nil, 0, false
	}
	if !hasInstallRepoObserverSubcommand(args) {
		// The installer namespace is intentionally closed until another portable
		// installer contract is added.
		// Keep install scoped to repo-observer so future installers do not share
		// ambiguous flag contracts.
		fmt.Fprintln(stderr, "install requires repo-observer")
		return nil, exitUsage, false
	}
	opts := installRepoObserverFlagSet()
	if err := opts.parse(args[1:]); err != nil {
		// Parse only the arguments after the required repo-observer verb.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return requireInstallRepoObserverFlags(opts, stderr)
}

func requireInstallRepoObserverFlags(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	if len(opts.rest()) != 0 {
		// Repo-observer install is fully flag-driven; no positional repository
		// path is interpreted here.
		fmt.Fprintln(stderr, "install repo-observer accepts only flags")
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func hasInstallRepoObserverSubcommand(args []string) bool {
	return len(args) != 0 && args[0] == "repo-observer"
}

func installRepoObserverFlagSet() *flagSet {
	opts := &flagSet{name: "install repo-observer"}
	// Default to the only supported portable repo-observer profile.
	opts.setString("profile", repoobserver.ProfileGithubActionsGitHooksV1)
	opts.setString("repository-id", "")
	opts.setString("out", "")
	opts.setBool("write", false)
	opts.setBool("force", false)
	return opts
}

func repoObserverOptionsFromFlags(opts *flagSet) repoobserver.Options {
	// Convert flags directly into repoobserver options so install/doctor keep a
	// single source of truth for profile semantics.
	return repoobserver.Options{
		Profile:      opts.stringValue("profile"),
		RepositoryID: opts.stringValue("repository-id"),
		Write:        opts.boolValue("write"),
		Force:        opts.boolValue("force"),
	}
}

func repoObserverExitCode(status repoobserver.Status) int {
	if status.InstallState == repoobserver.StateCannotVerify || status.ProofState == repoobserver.StateCannotVerify {
		// Cannot-verify install/proof state stays distinct from failed install.
		return exitCannotVerify
	}
	if status.InstallState == repoobserver.StateFail {
		return 1
	}
	return 0
}

func runVerify(_ context.Context, args []string, stdout, stderr io.Writer) int {
	runDir, code, ok := parseVerifyArgs(args, stderr)
	if !ok {
		return code
	}
	// VerifyRun computes the verdict and derived artifacts from retained run
	// evidence; artifact writing happens even when verification reports a
	// semantic error so failures remain inspectable.
	result, table, audit, err := verifier.VerifyRun(runDir)
	if writeErr := verifier.WriteVerifierArtifacts(runDir, result, table, audit); writeErr != nil {
		fmt.Fprintf(stderr, "failed writing verifier artifacts for %s: %v\n", runDir, writeErr)
		return 1
	}
	// Stdout carries the structured verifier result after artifacts are written
	// so terminal consumers cannot observe a result that was not retained.
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", data)
	if err != nil {
		// The JSON result is still emitted before the diagnostic so automation
		// can capture structured state.
		fmt.Fprintf(stderr, "%v\n", err)
	}
	return verifierResultExitCode(result.Result)
}

func parseVerifyArgs(args []string, stderr io.Writer) (string, int, bool) {
	if len(args) == 0 {
		// Verify requires a concrete retained run directory.
		fmt.Fprintln(stderr, "verify requires <run-dir>")
		return "", exitUsage, false
	}
	runDir := args[0]
	if !existingDirectory(runDir) {
		// Missing run roots are cannot_verify rather than usage once a path was
		// supplied.
		fmt.Fprintf(stderr, "run directory does not exist: %s\n", runDir)
		return "", exitCannotVerify, false
	}
	return runDir, 0, true
}

func existingDirectory(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

var verifierResultExitCodes = map[trace.VerifierVerdict]int{
	trace.VerdictObserved:     0,
	trace.VerdictNotAssessed:  0,
	trace.VerdictFail:         1,
	trace.VerdictCannotVerify: exitCannotVerify,
}

func verifierResultExitCode(result trace.VerifierVerdict) int {
	code, ok := verifierResultExitCodes[result]
	if !ok {
		// Unknown future verifier verdicts should not fail old automation by
		// default; schema validation remains the stronger compatibility check.
		return 0
	}
	return code
}

func runExplain(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		// Explanation is anchored to a retained run directory, not free-form
		// evidence text.
		fmt.Fprintln(stderr, "explain requires <run-dir>")
		return exitUsage
	}
	runDir := args[0]
	// ExplainRun renders a derived human view and does not change verifier
	// artifacts or verdict state.
	explanation, err := verifier.ExplainRun(runDir)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintln(stdout, explanation)
	return 0
}

func runQuery(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "query"}
	// Query mode is selected explicitly; the run directory is the only
	// positional evidence source.
	opts.setString("query", "")
	if err := opts.parse(args); err != nil {
		// Parser errors are usage failures before any retained run is inspected.
		return exitUsage
	}
	queryName := opts.stringValue("query")
	runDirs := opts.rest()
	if len(runDirs) == 0 {
		// Query reads one retained run directory as its evidence source.
		fmt.Fprintln(stderr, "query requires <run-dir>")
		return exitUsage
	}
	// Only the first retained run is accepted by the current query contract;
	// extra positional arguments remain outside the stable command surface.
	payload, code, ok := runNamedQuery(queryName, runDirs[0], stderr)
	if !ok {
		return code
	}
	// Query payloads are emitted as raw JSON bytes from the query package so the
	// CLI cannot alter diagnostic shape.
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runNamedQuery(queryName, runDir string, stderr io.Writer) ([]byte, int, bool) {
	if queryName == query.QueryCaptureDepth {
		// Capture-depth is a read-only diagnostic query over retained evidence.
		return captureDepthQuery(runDir, stderr)
	}
	if queryName != query.QueryMissingEvidence {
		// Unsupported query names are usage errors, not empty findings.
		fmt.Fprintf(stderr, "unsupported query: %s\n", queryName)
		return nil, exitUsage, false
	}
	return missingEvidenceQuery(runDir, stderr)
}

func captureDepthQuery(runDir string, stderr io.Writer) ([]byte, int, bool) {
	payload, err := query.CaptureDepth(runDir)
	if err != nil {
		// Query load/replay failures mean the retained evidence cannot be
		// verified for this diagnostic.
		fmt.Fprintln(stderr, err)
		return nil, exitCannotVerify, false
	}
	return payload, 0, true
}

func missingEvidenceQuery(runDir string, stderr io.Writer) ([]byte, int, bool) {
	payload, err := query.MissingEvidence(runDir)
	if err != nil {
		// Missing-evidence query failures are cannot_verify for the query
		// result, not an empty missing-evidence list.
		fmt.Fprintln(stderr, err)
		return nil, exitCannotVerify, false
	}
	return payload, 0, true
}

func runQueryPack(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 && args[0] == "explain" {
		// Explain renders an existing query-pack result; build creates one.
		return runQueryPackExplain(args[1:], stdout, stderr)
	}
	return runQueryPackBuild(args, stderr)
}

func runQueryPackBuild(args []string, stderr io.Writer) int {
	opts, err := parseQueryPackArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if err := validateQueryPackOptions(opts); err != nil {
		// Pack/profile validation happens before reading run artifacts so bad
		// command shape cannot be mistaken for unverifiable evidence.
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	// Query-pack build writes a portable JSON artifact for later explanation
	// and review.
	code, err := writeQueryPackArtifact(opts)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return code
	}
	return 0
}

func writeQueryPackArtifact(opts *queryPackOptions) (int, error) {
	result, err := query.ForensicsBasicPack(opts.runPath)
	if err != nil {
		// Pack generation depends on replayable run evidence.
		return exitCannotVerify, err
	}
	if err := writeJSONFile(opts.outPath, result); err != nil {
		// A generated pack that cannot be persisted is not review evidence.
		return 1, err
	}
	return 0, nil
}

func runQueryPackExplain(args []string, stdout, stderr io.Writer) int {
	opts, err := parseQueryPackExplainArgs(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	// Explain renders an existing query-pack artifact; it does not rebuild the
	// pack or re-query the original run evidence.
	result, err := readQueryPackResult(opts.resultPath)
	if err != nil {
		// Explain is artifact-only; it cannot reconstruct missing pack results
		// from the original run.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	if err := validateQueryPackExplainResult(result); err != nil {
		// Schema/profile mismatch means this binary cannot render the artifact
		// without risking stale explanation semantics.
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	fmt.Fprint(stdout, query.ExplainForensicsPack(result))
	return 0
}

type queryPackOptions struct {
	pack    string
	runPath string
	outPath string
}

type queryPackExplainOptions struct {
	resultPath string
}

func parseQueryPackArgs(args []string) (*queryPackOptions, error) {
	opts := &flagSet{name: "query-pack"}
	// Pack, run, and output flags are captured before validation so diagnostics
	// can distinguish unsupported pack names from missing evidence paths.
	opts.setString("pack", "")
	opts.setString("run", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		return nil, err
	}
	if len(opts.rest()) != 0 {
		// Pack, run evidence, and output path must be explicit provenance flags.
		// Positional arguments would be omitted from the persisted pack metadata.
		return nil, fmt.Errorf("query-pack accepts only flags")
	}
	return &queryPackOptions{
		// Trimmed values prevent whitespace-only flags from satisfying required
		// evidence-path checks.
		pack:    strings.TrimSpace(opts.stringValue("pack")),
		runPath: strings.TrimSpace(opts.stringValue("run")),
		outPath: strings.TrimSpace(opts.stringValue("out")),
	}, nil
}

func parseQueryPackExplainArgs(args []string) (*queryPackExplainOptions, error) {
	opts := &flagSet{name: "query-pack explain"}
	// Explanation takes one persisted result artifact; it has no run-directory
	// fallback.
	opts.setString("result", "")
	if err := opts.parse(args); err != nil {
		return nil, err
	}
	if len(opts.rest()) != 0 {
		return nil, fmt.Errorf("query-pack explain accepts only flags")
	}
	resultPath := strings.TrimSpace(opts.stringValue("result"))
	if resultPath == "" {
		// Explanation is bound to a persisted result artifact, not stdin or a
		// transient in-memory pack.
		return nil, fmt.Errorf("query-pack explain requires --result")
	}
	return &queryPackExplainOptions{resultPath: resultPath}, nil
}

func validateQueryPackOptions(opts *queryPackOptions) error {
	if opts.pack == "" {
		return fmt.Errorf("error: ambiguous pack selection; --pack is required")
	}
	if opts.pack != query.QueryPackForensicsBasic {
		// The CLI exposes a closed pack vocabulary so unknown pack names fail as
		// usage errors before any evidence is read or written.
		return fmt.Errorf("error: unknown pack %q", opts.pack)
	}
	return requireQueryPackRequiredInputs(opts.runPath, opts.outPath)
}

func requireQueryPackRequiredInputs(runPath, outPath string) error {
	if runPath == "" {
		// The run path is the replayable source evidence for this pack.
		return fmt.Errorf("query-pack requires --run")
	}
	if outPath == "" {
		// The output path is the durable artifact reviewed by later commands.
		return fmt.Errorf("query-pack requires --out")
	}
	return nil
}

func readQueryPackResult(path string) (query.QueryPackResult, error) {
	var result query.QueryPackResult
	if err := readJSONFile(path, &result); err != nil {
		// Artifact read failures are verification failures, not empty results.
		return query.QueryPackResult{}, err
	}
	return result, nil
}

func validateQueryPackExplainResult(result query.QueryPackResult) error {
	if result.SchemaVersion != query.QueryPackSchemaVersion || result.QueryPackID != query.QueryPackForensicsBasic {
		// Explain only understands the current forensics-basic result contract.
		return fmt.Errorf("unsupported query-pack result")
	}
	// The query package owns detailed result validation; CLI validation only
	// gates the schema/profile pair before rendering.
	return nil
}

func runExport(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if exportTelemetryRequested(args) {
		// Telemetry export consumes posture output and emits Prometheus text.
		return runTelemetryExport(args[1:], stdout, stderr)
	}
	if exportCrossRepoPostureExplainRequested(args) {
		// Cross-repo posture explanations render existing exports only.
		return runCrossRepoPostureExplain(args[2:], stdout, stderr)
	}
	if exportCrossRepoPostureRequested(args) {
		// Cross-repo posture export builds the durable posture artifact.
		return runCrossRepoPostureExport(args[1:], stdout, stderr)
	}
	// Export uses a closed command vocabulary; unsupported exports are usage
	// errors, not unverifiable evidence states.
	fmt.Fprintln(stderr, "export requires cross-repo-posture or telemetry")
	return exitUsage
}

func exportTelemetryRequested(args []string) bool {
	return exportCommandIs(args, "telemetry")
}

func exportCrossRepoPostureExplainRequested(args []string) bool {
	return exportCommandIs(args, "cross-repo-posture") && exportSubcommandIs(args, "explain")
}

func exportCrossRepoPostureRequested(args []string) bool {
	return exportCommandIs(args, "cross-repo-posture")
}

func exportCommandIs(args []string, command string) bool {
	return len(args) > 0 && args[0] == command
}

func exportSubcommandIs(args []string, subcommand string) bool {
	return len(args) > 1 && args[1] == subcommand
}

func runTelemetryExport(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseTelemetryExportArgs(args, stderr)
	if !ok {
		return code
	}
	// Telemetry is rendered from an already-built posture artifact so the export
	// layer cannot silently broaden repository selection.
	rendered, err := renderTelemetryExport(opts.stringValue("cross-repo-posture"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	if err := writeTelemetryExportOutput(opts.stringValue("out"), rendered, stdout); err != nil {
		// Export errors happen after rendering; no partial metric file is
		// accepted as evidence.
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func parseTelemetryExportArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "export telemetry"}
	// Telemetry export accepts only the renderer profile, posture artifact, and
	// output target needed to replay metric generation.
	opts.setString("profile", "")
	opts.setString("cross-repo-posture", "")
	opts.setString("out", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Profile, posture artifact, and destination stay named so exports are
	// auditable from the command line alone.
	if rejectRest(opts, stderr, "export telemetry accepts only flags") {
		return nil, exitUsage, false
	}
	return requireTelemetryExportArgs(opts, stderr)
}

func requireTelemetryExportArgs(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	if err := requireTelemetryExportInputs(opts); err != nil {
		// Required input checks keep unsupported profiles and missing artifacts
		// as usage errors before any metric bytes are emitted.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func renderTelemetryExport(posturePath string) (string, error) {
	var result posture.ExportResult
	if err := readJSONFile(posturePath, &result); err != nil {
		// Missing or malformed posture input means telemetry cannot be verified.
		return "", fmt.Errorf("posture_unreadable")
	}
	rendered, err := telemetry.RenderPrometheus(result)
	if err != nil {
		// Rendering failures preserve cannot_verify instead of emitting partial
		// metrics.
		return "", fmt.Errorf("telemetry_cannot_verify")
	}
	return rendered, nil
}

func writeTelemetryExportOutput(outPath, rendered string, stdout io.Writer) error {
	if outPath == "-" {
		// Dash keeps review-friendly stdout output without changing the rendered
		// payload.
		fmt.Fprint(stdout, rendered)
		return nil
	}
	if err := writeTextFileAtomic(outPath, rendered); err != nil {
		// File output is all-or-nothing; partial metric files are not accepted as
		// evidence.
		return fmt.Errorf("out_unwritable")
	}
	return nil
}

func requireTelemetryExportInputs(opts *flagSet) error {
	if strings.TrimSpace(opts.stringValue("profile")) != telemetry.ProfilePrometheusTextV1 {
		// Telemetry export is intentionally profile-locked so future renderers
		// cannot be selected by typo or stale docs.
		return fmt.Errorf("export telemetry requires --profile prometheus-text-v1")
	}
	if strings.TrimSpace(opts.stringValue("cross-repo-posture")) == "" {
		// The posture artifact is the sole metric source; the CLI does not infer
		// repository posture from the working tree.
		return fmt.Errorf("export telemetry requires --cross-repo-posture")
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		// Metrics must either be explicitly written or deliberately streamed with
		// `--out -`.
		return fmt.Errorf("export telemetry requires --out")
	}
	return nil
}

func runCrossRepoPostureExport(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseCrossRepoPostureExportArgs(args, stderr)
	if !ok {
		return code
	}
	// Build reads the declared selection file and produces one posture artifact;
	// stdout is intentionally ignored so the artifact path stays authoritative.
	result, err := posture.Build(opts.stringValue("selection"), time.Now())
	if err != nil {
		// Build failures are reported without leaking selection parse details
		// into a misleading export artifact.
		fmt.Fprintln(stderr, "no_export_artifact")
		return exitCannotVerify
	}
	_ = stdout
	return writeCrossRepoPostureExport(opts, result, stderr)
}

func parseCrossRepoPostureExportArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "export cross-repo-posture"}
	// Profile names the posture contract, selection names the repository set,
	// and out names the durable result artifact.
	opts.setString("profile", "")
	opts.setString("selection", "")
	opts.setString("out", "")
	opts.setBool("validate-only", false)
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	// Selection, profile, and output remain explicit flags for replayable export
	// provenance.
	if rejectRest(opts, stderr, "export cross-repo-posture accepts only flags") {
		return nil, exitUsage, false
	}
	return requireCrossRepoPostureExportArgs(opts, stderr)
}

func requireCrossRepoPostureExportArgs(opts *flagSet, stderr io.Writer) (*flagSet, int, bool) {
	if err := requireCrossRepoPostureInputs(opts); err != nil {
		// The profile and selection file are mandatory even in validate-only
		// mode because they define the posture evidence boundary.
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	return opts, 0, true
}

func writeCrossRepoPostureExport(opts *flagSet, result posture.ExportResult, stderr io.Writer) int {
	if opts.boolValue("validate-only") {
		// Validate-only proves the selection can build without publishing a new
		// posture artifact.
		return 0
	}
	if strings.TrimSpace(opts.stringValue("out")) == "" {
		// Non-preview exports must name the durable posture artifact path.
		fmt.Fprintln(stderr, "export cross-repo-posture requires --out")
		return exitUsage
	}
	if err := writeJSONFile(opts.stringValue("out"), result); err != nil {
		// Failed publication leaves no reviewable posture export.
		fmt.Fprintln(stderr, "out_unwritable")
		return 1
	}
	return 0
}

func requireCrossRepoPostureInputs(opts *flagSet) error {
	if strings.TrimSpace(opts.stringValue("profile")) != posture.ProfileID {
		// Keep the CLI bound to the only supported cross-repo posture contract.
		return fmt.Errorf("export cross-repo-posture requires --profile cross-repo-evidence-posture-v1")
	}
	if strings.TrimSpace(opts.stringValue("selection")) == "" {
		// The selection artifact is the auditable repository set for posture.
		return fmt.Errorf("export cross-repo-posture requires --selection")
	}
	return nil
}

func runCrossRepoPostureExplain(args []string, stdout, stderr io.Writer) int {
	opts, code, ok := parseCrossRepoPostureExplainArgs(args, stderr)
	if !ok {
		return code
	}
	// Explanation renders a saved posture export; it never rebuilds selection
	// state from the workspace.
	result, code, ok := readCrossRepoPostureExplainResult(opts.stringValue("result"), stderr)
	if !ok {
		return code
	}
	rendered, err := posture.Explain(result)
	if err != nil {
		// Unsafe rendered text is a verification failure for the explanation.
		fmt.Fprintln(stderr, "output_safety_violation")
		return exitCannotVerify
	}
	// The explanation is intentionally stdout-only so it cannot be mistaken for
	// a new posture evidence artifact.
	fmt.Fprint(stdout, rendered)
	return 0
}

func parseCrossRepoPostureExplainArgs(args []string, stderr io.Writer) (*flagSet, int, bool) {
	opts := &flagSet{name: "export cross-repo-posture-explain"}
	// The result flag points at the posture artifact that will be explained.
	opts.setString("result", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return nil, exitUsage, false
	}
	if rejectRest(opts, stderr, "export cross-repo-posture-explain accepts only flags") {
		return nil, exitUsage, false
	}
	if strings.TrimSpace(opts.stringValue("result")) == "" {
		// A persisted export result is required before human explanation.
		fmt.Fprintln(stderr, "export cross-repo-posture-explain requires --result")
		return nil, exitUsage, false
	}
	// Successful parsing only identifies the artifact; schema and output-safety
	// checks happen when the artifact is read and rendered.
	return opts, 0, true
}

func readCrossRepoPostureExplainResult(path string, stderr io.Writer) (posture.ExportResult, int, bool) {
	var result posture.ExportResult
	if err := readJSONFile(path, &result); err != nil {
		// Missing or malformed result artifacts cannot be explained honestly.
		fmt.Fprintln(stderr, "result_unreadable")
		return posture.ExportResult{}, exitCannotVerify, false
	}
	if result.SchemaVersion != posture.SchemaVersion || result.ExportProfileID != posture.ProfileID {
		// Explain only accepts the current posture export schema/profile pair.
		fmt.Fprintln(stderr, "unsupported cross-repo posture export")
		return posture.ExportResult{}, exitCannotVerify, false
	}
	// Shape validation is delegated to posture.Explain, which also enforces
	// output-safety before rendering.
	return result, 0, true
}

func runValidateFixtures(_ context.Context, args []string, stdout, stderr io.Writer) int {
	fixtureRoot := fixtureRootArg(args)
	// Fixture discovery is rooted explicitly so validation cannot wander into
	// unrelated run artifacts.
	runDirs, err := demo.DiscoverRunDirs(fixtureRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if validateFixtureRuns(fixtureRoot, runDirs, stdout, stderr) {
		return 1
	}
	return 0
}

func fixtureRootArg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}
	// Default fixtures to the current directory for local demo validation.
	return "."
}

func validateFixtureRuns(fixtureRoot string, runDirs []string, stdout, stderr io.Writer) bool {
	failed := false
	for _, runDir := range runDirs {
		// Continue through all fixtures so one broken run does not hide other
		// drift in the example corpus.
		if validateFixtureRun(fixtureRoot, runDir, stdout, stderr) {
			failed = true
		}
	}
	return failed
}

func validateFixtureRun(fixtureRoot, runDir string, stdout, stderr io.Writer) bool {
	result, table, audit, verifyErr := verifier.VerifyRun(runDir)
	if err := verifier.WriteVerifierArtifacts(runDir, result, table, audit); err != nil {
		// Verifier artifacts are part of the fixture evidence, even when replay
		// reports semantic verification errors.
		fmt.Fprintf(stderr, "failed writing verifier artifacts for %s: %v\n", runDir, err)
		return true
	}
	fmt.Fprintf(stdout, "%s => %s\n", runDir, result.Result)
	if verifyErr != nil {
		// Surface replay diagnostics but still compare the structured verdict
		// against the fixture expectation.
		fmt.Fprintf(stderr, "%s verification error: %v\n", runDir, verifyErr)
	}
	return fixtureExpectationFailed(fixtureRoot, runDir, result, stderr)
}

func fixtureExpectationFailed(fixtureRoot, runDir string, result trace.VerifierResult, stderr io.Writer) bool {
	expectation, err := readFixtureExpectation(fixtureRoot, runDir)
	if err != nil {
		// Bad expectation metadata is fixture drift, not a verifier pass.
		fmt.Fprintf(stderr, "invalid fixture expectation for %s: %v\n", runDir, err)
		return true
	}
	if expectation.ExpectedResult != "" {
		// Explicit fixture expectations define the authoritative verdict.
		return expectedFixtureResultFailed(runDir, result, expectation, stderr)
	}
	// Fixtures without explicit expectations may still fail if replay proves a
	// hard verifier failure or cannot-verify state.
	return unexpectedFixtureResultFailed(result)
}

func expectedFixtureResultFailed(runDir string, result trace.VerifierResult, expectation fixtureExpectation, stderr io.Writer) bool {
	if expectation.ExpectedResult == string(result.Result) {
		return false
	}
	// Mismatches are printed with the run path so fixture corpus drift is
	// actionable.
	fmt.Fprintf(stderr, "%s expected %s, got %s\n", runDir, expectation.ExpectedResult, result.Result)
	return true
}

func unexpectedFixtureResultFailed(result trace.VerifierResult) bool {
	// Without an explicit expected result, only fail/cannot_verify are treated as
	// fixture failures; observed/not_assessed remain inspectable but nonfatal.
	return result.Result == trace.VerdictFail || result.Result == trace.VerdictCannotVerify
}

type doctorReport struct {
	Command            string        `json:"command"`
	Result             string        `json:"result"`
	Environment        []doctorCheck `json:"environment"`
	ControlPoints      []doctorCheck `json:"control_points"`
	SafeRetentionModes []string      `json:"safe_retention_modes"`
}

type doctorCheck struct {
	ID        string   `json:"id"`
	State     string   `json:"state"`
	Reason    string   `json:"reason"`
	Contract  string   `json:"contract_id,omitempty"`
	Missing   []string `json:"missing,omitempty"`
	Reference string   `json:"reference,omitempty"`
}

type doctorOptions struct {
	ContractPath string
	OutputDir    string
	ReportDir    string
	Env          map[string]string
}

type previewBoundary struct {
	Boundary string `json:"boundary"`
	State    string `json:"state"`
	Reason   string `json:"reason"`
}

type previewOfflineImplication struct {
	Requirement string `json:"requirement"`
	State       string `json:"state"`
	Reason      string `json:"reason"`
}

const (
	defaultRunRoot   = ".sdp-trace-runs"
	defaultReportDir = ".sdp-trace-report"
)

func buildDoctorReport(options doctorOptions) (doctorReport, int) {
	defaultContract := trace.DefaultContract
	result := "offline_dev"
	exitCode := 0
	// Doctor reports local readiness; it never upgrades offline evidence to a
	// CI or external witness claim.
	contract, contractCheck := doctorContractCheck(options.ContractPath, defaultContract)
	result, exitCode = updateDoctorExitForCheck(result, exitCode, contractCheck)
	ciCheck := ciWitnessPrerequisiteCheck(options.Env)
	// Writable path checks prove only local filesystem readiness for future
	// artifacts; they do not create recorder evidence.
	outputDirCheck := writablePathCheck("output_directory", options.OutputDir, "run artifact output directory is writable")
	reportDirCheck := writablePathCheck("report_directory", options.ReportDir, "report artifact directory is writable")
	expectedEvidenceCheck := expectedEvidenceReferenceCheck(contract)
	// CI prerequisites are reported, but local doctor does not require them for
	// offline development readiness.
	// The contract check runs before evidence-reference checks so a bad
	// override cannot be hidden by default-contract coverage.
	result, exitCode = updateDoctorExitForLocalChecks(result, exitCode, outputDirCheck, reportDirCheck, expectedEvidenceCheck)
	// Contract loading and writable probes are live checks; CI identity remains
	// separately marked cannot_verify when absent.
	// Doctor assembles named facts only; the result field reflects verifier
	// state without becoming an opaque score.
	report := doctorReport{
		Command:     "doctor",
		Result:      result,
		Environment: doctorEnvironmentChecks(),
		// Environment and control-point sections stay separate so local process
		// facts cannot be mistaken for gate evidence.
		// Control points mix local passes and cannot-verify prerequisites so the
		// JSON report does not collapse them into one health score.
		ControlPoints:      doctorControlPointChecks(defaultContract, ciCheck, outputDirCheck, reportDirCheck, contractCheck, expectedEvidenceCheck),
		SafeRetentionModes: safeRetentionModes(),
	}
	// The report intentionally contains no aggregate score; callers inspect
	// named control points.
	return report, exitCode
}

func updateDoctorExitForLocalChecks(result string, exitCode int, checks ...doctorCheck) (string, int) {
	for _, check := range checks {
		// Only control points that the local process can inspect affect the
		// offline doctor exit code.
		result, exitCode = updateDoctorExitForCheck(result, exitCode, check)
	}
	return result, exitCode
}

func doctorEnvironmentChecks() []doctorCheck {
	// Environment checks describe what the local process can know; they are not
	// external witness evidence.
	return []doctorCheck{
		{
			ID:     "local_process",
			State:  "pass",
			Reason: "current process can inspect local environment",
		},
		{
			ID:     "offline_development",
			State:  "offline_dev",
			Reason: "external CI identity is not required for local preview or wrapper readiness",
		},
	}
}

func doctorControlPointChecks(defaultContract trace.Contract, ciCheck doctorCheck, checks ...doctorCheck) []doctorCheck {
	// Control-point order is stable because downstream fixtures and reports cite
	// named gaps in the order doctor prints them.
	controlPoints := []doctorCheck{
		{
			ID:     "local_wrapper",
			State:  "pass",
			Reason: "wrap and run commands are registered in this binary",
		},
	}
	controlPoints = append(controlPoints, checks...)
	// The built-in contract is reported before CI prerequisites so local
	// readiness remains separate from unavailable external identity.
	controlPoints = append(controlPoints, doctorDefaultContractCheck(defaultContract))
	return append(controlPoints, ciCheck)
}

func doctorDefaultContractCheck(defaultContract trace.Contract) doctorCheck {
	// The default contract is a local development fallback, not proof that a
	// specific repository run satisfied its contract.
	return doctorCheck{
		ID:        "default_contract",
		State:     "pass",
		Reason:    "built-in contract is available for local development",
		Contract:  defaultContract.ContractID,
		Reference: defaultContract.Version,
	}
}

func doctorContractCheck(contractPath string, defaultContract trace.Contract) (trace.Contract, doctorCheck) {
	if contractPath == "" {
		// No override means doctor can only report local default-contract
		// availability.
		return defaultDoctorContractResult(defaultContract)
	}
	// Loading the requested contract is the only source-bound proof for an
	// override contract path.
	contract, err := trace.LoadContract(contractPath)
	if err != nil {
		// Load failures lower only this control point to cannot_verify.
		return unreadableDoctorContractResult(defaultContract, contractPath)
	}
	// A loaded override becomes the contract returned to later doctor checks.
	return loadedDoctorContractResult(contract, contractPath)
}

func loadedDoctorContractResult(contract trace.Contract, contractPath string) (trace.Contract, doctorCheck) {
	// Loaded override contracts are reported with their id and path so doctor
	// output remains replayable from the local command invocation.
	return contract, doctorCheck{
		ID:        "contract",
		State:     "pass",
		Reason:    "contract can be loaded",
		Contract:  contract.ContractID,
		Reference: contractPath,
	}
}

func defaultDoctorContractResult(defaultContract trace.Contract) (trace.Contract, doctorCheck) {
	// The default contract branch is explicitly local evidence; it says nothing
	// about a repository-specific contract file.
	return defaultContract, doctorCheck{
		ID:        "contract",
		State:     "pass",
		Reason:    "default contract is available",
		Contract:  defaultContract.ContractID,
		Reference: "local-default-v1",
	}
}

func unreadableDoctorContractResult(defaultContract trace.Contract, contractPath string) (trace.Contract, doctorCheck) {
	// A requested contract that cannot load keeps doctor in cannot_verify.
	return defaultContract, doctorCheck{
		ID:        "contract",
		State:     string(trace.VerdictCannotVerify),
		Reason:    "contract cannot be loaded",
		Reference: contractPath,
	}
}

func updateDoctorExitForCheck(result string, exitCode int, check doctorCheck) (string, int) {
	if check.State == string(trace.VerdictCannotVerify) {
		// Any cannot-verify control point lowers the overall doctor exit.
		return string(trace.VerdictCannotVerify), exitCannotVerify
	}
	return result, exitCode
}

func writablePathCheck(id, path, okReason string) doctorCheck {
	if strings.TrimSpace(path) == "" {
		return emptyWritablePathCheck(id)
	}
	target, check, ok := writableProbeTarget(id, path)
	if !ok {
		// Path-shape failures are returned before probing so a file is never
		// treated as an artifact directory candidate.
		return check
	}
	return probeWritablePath(id, path, target, okReason)
}

func probeWritablePath(id, path, target, okReason string) doctorCheck {
	// Missing directories are probed through their parent, while existing
	// directories are probed directly.
	// Probe with a temporary file so doctor validates actual write capability,
	// not just path syntax.
	probe, err := os.CreateTemp(target, ".sdp-trace-doctor-")
	if err != nil {
		// A failed probe is recorded as cannot_verify rather than inferred from
		// permissions text or path shape.
		return doctorCheck{
			ID:        id,
			State:     string(trace.VerdictCannotVerify),
			Reason:    "directory is not writable",
			Reference: path,
		}
	}
	probeName := probe.Name()
	// Probe cleanup failures are intentionally ignored; the check is about
	// whether a report/run artifact could be written.
	_ = probe.Close()
	_ = os.Remove(probeName)
	return writablePathPassCheck(id, path, okReason)
}

func emptyWritablePathCheck(id string) doctorCheck {
	// Empty configured paths cannot support durable run/report artifacts.
	return doctorCheck{
		ID:     id,
		State:  string(trace.VerdictCannotVerify),
		Reason: "path is empty",
	}
}

func writablePathPassCheck(id, path, okReason string) doctorCheck {
	// The temporary probe is removed immediately; doctor reports capability,
	// not an artifact to retain.
	return doctorCheck{
		ID:        id,
		State:     "pass",
		Reason:    okReason,
		Reference: path,
	}
}

func writableProbeTarget(id, path string) (string, doctorCheck, bool) {
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		// Existing non-directories cannot become artifact roots.
		return "", doctorCheck{
			ID:        id,
			State:     string(trace.VerdictCannotVerify),
			Reason:    "path exists but is not a directory",
			Reference: path,
		}, false
	}
	if os.IsNotExist(err) {
		// Missing artifact roots are checked by probing their parent directory.
		return writableProbeParent(path), doctorCheck{}, true
	}
	// Other stat errors are left to CreateTemp so the observed write failure is
	// reported through one path.
	return path, doctorCheck{}, true
}

func writableProbeParent(path string) string {
	target := filepath.Dir(path)
	if target == "" {
		// Empty dirname resolves to the current directory for local probes.
		return "."
	}
	return target
}

func expectedEvidenceReferenceCheck(contract trace.Contract) doctorCheck {
	if len(contract.RequiredEvents) == 0 {
		return expectedEvidenceNoRequiredEventsCheck(contract)
	}
	// Required event and evidence references are checked separately so drift is
	// reported with concrete missing keys.
	missing := expectedEvidenceReferenceGaps(contract)
	if len(missing) > 0 {
		return expectedEvidenceUnsupportedReferenceCheck(contract, missing)
	}
	// A pass here means only that the local event vocabulary can represent the
	// contract's required evidence references.
	return doctorCheck{
		ID:       "expected_evidence_references",
		State:    "pass",
		Reason:   "contract required events and evidence references are supported by the current local event model",
		Contract: contract.ContractID,
	}
}

func expectedEvidenceNoRequiredEventsCheck(contract trace.Contract) doctorCheck {
	// A contract with no required events cannot prove evidence coverage.
	return doctorCheck{
		ID:       "expected_evidence_references",
		State:    string(trace.VerdictCannotVerify),
		Reason:   "contract has no required_events",
		Contract: contract.ContractID,
	}
}

func expectedEvidenceUnsupportedReferenceCheck(contract trace.Contract, missing []string) doctorCheck {
	// Unsupported references are reported as explicit gaps, not hidden in a
	// generic contract failure.
	return doctorCheck{
		ID:       "expected_evidence_references",
		State:    string(trace.VerdictCannotVerify),
		Reason:   "contract references unsupported event types",
		Contract: contract.ContractID,
		Missing:  missing,
	}
}

func expectedEvidenceReferenceGaps(contract trace.Contract) []string {
	missing := make([]string, 0)
	for _, eventType := range contract.RequiredEvents {
		if !knownEventType(eventType) {
			// Required events must map to this binary's local event vocabulary.
			missing = append(missing, "required_events:"+eventType)
		}
	}
	// Evidence-specific gaps are appended after required event gaps so the
	// output separates event vocabulary drift from requirement-shape drift.
	for _, evidence := range contract.RequiredEvidence {
		missing = append(missing, expectedEvidenceGaps(evidence)...)
	}
	return missing
}

func expectedEvidenceGaps(evidence trace.EvidenceRequirement) []string {
	missing := make([]string, 0, 2)
	if strings.TrimSpace(evidence.ID) == "" {
		// Missing evidence IDs make later diagnostics ambiguous.
		missing = append(missing, "required_evidence:<missing_id>")
	}
	if strings.TrimSpace(evidence.EventType) == "" {
		// Missing event type is reported separately from unsupported types.
		return append(missing, "required_evidence:"+evidence.ID+":<missing_event_type>")
	}
	if !knownEventType(evidence.EventType) {
		// Evidence requirements must reference event types this binary can emit.
		missing = append(missing, "required_evidence:"+evidence.ID+":"+evidence.EventType)
	}
	return missing
}

func knownEventType(eventType string) bool {
	switch trace.EventType(eventType) {
	case trace.EventRecorderAttached,
		trace.EventRunStarted,
		trace.EventCommandStarted,
		trace.EventCommandFinished,
		trace.EventRunClosed,
		trace.EventPolicyOverrideRequested:
		// Keep doctor scoped to the stable local recorder event model.
		// Unsupported future events remain explicit spec-drift gaps.
		return true
	default:
		return false
	}
}

func ciWitnessPrerequisiteCheck(env map[string]string) doctorCheck {
	missing := missingCIWitnessFields(env)
	if len(missing) > 0 {
		// Local environments usually cannot produce CI witness evidence.
		return doctorCheck{
			ID:      "ci_witness_prerequisites",
			State:   string(trace.VerdictCannotVerify),
			Reason:  "GitHub Actions identity or OIDC prerequisite is unavailable in this environment",
			Missing: missing,
		}
	}
	// Passing prerequisites only means the environment exposes the fields needed
	// for witness construction; it is not a witness verdict.
	return doctorCheck{
		ID:     "ci_witness_prerequisites",
		State:  "pass",
		Reason: "GitHub Actions identity and OIDC prerequisites are present",
	}
}

func missingCIWitnessFields(env map[string]string) []string {
	missing := missingEnvFields(env, requiredCIWitnessEnvFields())
	if env["GITHUB_ACTIONS"] != "" && env["GITHUB_ACTIONS"] != "true" {
		// GitHub exposes the flag as literal true for Actions-backed identity.
		missing = append(missing, "GITHUB_ACTIONS=true")
	}
	return missing
}

func requiredCIWitnessEnvFields() []string {
	// Require both OIDC request fields and workflow identity fields for CI
	// witness construction.
	return []string{
		"ACTIONS_ID_TOKEN_REQUEST_TOKEN",
		"ACTIONS_ID_TOKEN_REQUEST_URL",
		"GITHUB_ACTIONS",
		"GITHUB_ACTOR",
		"GITHUB_JOB",
		"GITHUB_REF",
		"GITHUB_REPOSITORY",
		"GITHUB_RUN_ATTEMPT",
		"GITHUB_RUN_ID",
		"GITHUB_SERVER_URL",
		"GITHUB_SHA",
		"GITHUB_WORKFLOW",
	}
}

func missingEnvFields(env map[string]string, required []string) []string {
	missing := make([]string, 0)
	for _, key := range required {
		if strings.TrimSpace(env[key]) == "" {
			// Trim whitespace so empty exported variables are treated as absent.
			missing = append(missing, key)
		}
	}
	return missing
}

func safeRetentionModes() []string {
	// Doctor publishes retention modes that preserve replay without raw secret
	// exposure by default.
	return []string{
		string(trace.RetentionModeDigestOnly),
		string(trace.RetentionModeSanitizedExcerpt),
		string(trace.RetentionModeEncryptedRawRef),
		string(trace.RetentionModeExternalArtifactRef),
		string(trace.RetentionModeNotAssessed),
	}
}

var previewBoundaryRows = []previewBoundary{
	{
		Boundary: string(trace.ObservationBoundaryProcessWrapper),
		State:    "pass",
		Reason:   "preview covers local process-wrapper capture only",
	},
	{
		Boundary: string(trace.ObservationBoundaryAdapterSocket),
		State:    string(trace.ObservationStateNotIntegrated),
		Reason:   "adapter socket/API capture is not configured in Block 13B",
	},
	{
		Boundary: string(trace.ObservationBoundaryToolWrapper),
		State:    string(trace.ObservationStateUnsupported),
		Reason:   "tool-level wrapping is a future observation boundary",
	},
	{
		Boundary: string(trace.ObservationBoundaryVCSPRObserver),
		State:    string(trace.ObservationStateNotIntegrated),
		Reason:   "VCS/PR observer is not configured in Block 13B",
	},
	{
		Boundary: string(trace.ObservationBoundaryCIObserver),
		State:    string(trace.ObservationStateOfflineDev),
		Reason:   "CI witness cannot be produced by local preview",
	},
	{
		// External witness state remains not_integrated until a concrete witness
		// profile exists.
		Boundary: string(trace.ObservationBoundaryExternalWitness),
		State:    string(trace.ObservationStateNotIntegrated),
		Reason:   "external witness profile is not implemented in Block 13B",
	},
}

func previewBoundaries() []previewBoundary {
	// Preview is explicit about which observation boundaries are local,
	// unsupported, or not integrated.
	return append([]previewBoundary(nil), previewBoundaryRows...)
}

func previewOfflineImplications() []previewOfflineImplication {
	// Offline implications tell the user which evidence must be re-collected in
	// CI or external systems before trust can be upgraded.
	return []previewOfflineImplication{
		{
			Requirement: "ci_witnessed",
			State:       string(trace.ObservationStateOfflineDev),
			Reason:      "rerun in CI with OIDC before using CI witness evidence",
		},
		{
			Requirement: "external_witnessed",
			State:       string(trace.ObservationStateNotIntegrated),
			Reason:      "external witness profile is not implemented in Block 13B",
		},
	}
}

const usageText = `sdp-trace local recorder and verifier commands.

Usage:
  sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>
  sdp-trace version
  sdp-trace run --task <task-ref> [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace dry-run [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace preview [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace doctor [--contract <file>]
  sdp-trace doctor --profile github-actions-git-hooks-v1 [--out <file>]
  sdp-trace install repo-observer --profile github-actions-git-hooks-v1 [--repository-id <safe-id>] [--write] [--force] [--out <file>]
  sdp-trace interaction relay --task-id <safe-id> --event-type <type> --out <file> -- <forward-command...>
  sdp-trace interaction import-transcript --source preclassified-transcript-import --task-id <safe-id> --events-jsonl <file> --out <file>
  sdp-trace interaction summarize --trace <file> [--out <file>]
  sdp-trace observe setup --profile <session-profile.json> --out <run-dir> [--command <harness-command-preview>]
  sdp-trace observe collect --profile <session-profile.json> --run <run-dir>
  sdp-trace observe session --profile <session-profile.json> --out <run-dir> -- <harness-command...>
  sdp-trace harness observe --profile <harness-profile.json> --source <harness-events.jsonl> --out <run-dir>
  sdp-trace harness validate --profile <harness-profile.json> --run <run-dir> --out <validation.json>
  sdp-trace harness summarize --validation <validation.json>
  sdp-trace envelope summarize --envelope <file> [--out <file>]
  sdp-trace verify <run-dir>
  sdp-trace explain <run-dir>
  sdp-trace query --query <missing-evidence|capture-depth> <run-dir>
  sdp-trace query-pack --pack forensics-basic-v1 --run <run-dir> --out <file>
  sdp-trace query-pack explain --result <file>
  sdp-trace export cross-repo-posture --profile cross-repo-evidence-posture-v1 --selection <file> --out <file>
  sdp-trace export cross-repo-posture explain --result <file>
  sdp-trace export telemetry --profile prometheus-text-v1 --cross-repo-posture <file> --out <file|->
  sdp-trace assess --profile adapter-capture --out <file> --run <run-dir>
  sdp-trace assess --profile managed-harness --out <file> --contract <file> --run <run-dir> --adapter-registry <file> --managed-policy <file> --managed-witness <file>
  sdp-trace assess --profile forensic-retention --out <file> --run <run-dir> --redaction-policy <file>
  sdp-trace assess --profile ci-artifact-observation --out <file> --artifact-manifest <file>
  sdp-trace assess --profile authority-envelope --out <file> --authority-package <file>
  sdp-trace assess preview --profile <adapter-capture|managed-harness|forensic-retention|ci-artifact-observation|authority-envelope> [profile inputs]
  sdp-trace assess explain --assessment-result <file>
  sdp-trace report --out <dir> <runs-root-or-run-dir>
  sdp-trace gate --out <file> <runs-root-or-run-dir>
  sdp-trace witness --kind <github-actions|gitlab-ci|buildkite|customer-pki> --out <file> [--report-dir <dir>] [--witness-envelope <file>] [--customer-pki-authority-policy <file>] [--customer-pki-public-cert <file> | --customer-pki-public-key <file>] [--customer-pki-payload-digest <sha256>] [--customer-pki-freshness-evidence <file>] <runs-root-or-run-dir>
  sdp-trace release-proof --manifest <file> --out <file>
  sdp-trace pr-review packet --out <dir> --repo-id <safe-id> --change-ref <pr|mr|change-id> --base <sha> --head <sha> --diff <file> [--ci-state <state>] [--created-by <actor>]
  sdp-trace pr-review run --packet <dir> --profile <file> --out <dir> [--preview] [--work-dir <dir>] [--allow-external-runner <runner>]...
  sdp-trace pr-review synthesize --packet <dir> --runs <dir> --out <file>
  sdp-trace pr-review validate --packet <dir> --profile <file> --runs <dir> --ledger <file> --out <file>
  sdp-trace pr-review summarize --validation <file> --ledger <file> [--out <file>]
  sdp-trace pr-review check --out <dir> --repo-id <safe-id> --change-ref <pr|mr|change-id> --base <sha> --head <sha> --diff <file> --profile <file> [--work-dir <dir>] [--allow-external-runner <runner>]...
  sdp-trace packet build-pr --source <github-actions|github-fixture> --out <dir> [--github-event <file>] [--checks-json <file>] [--artifacts-json <file>] [--route-manifest <file>] [--github-api-url <url>]
  sdp-trace packet build-github --github-input <file> --out <file>
  sdp-trace packet validate --bundle <file>
  sdp-trace packet check-demo --bundle <file>
  sdp-trace packet render --bundle <file> --out <file>
  sdp-trace validate-fixtures [root-dir]
`

func printUsage(w io.Writer) {
	// Global help is the canonical local command contract for this small CLI.
	fmt.Fprint(w, usageText)
}

type fixtureExpectation struct {
	ExpectedResult string `json:"expected_result"`
}

func readFixtureExpectation(root, runDir string) (fixtureExpectation, error) {
	// Expectations are optional corpus metadata; absence leaves default verifier
	// failure handling in place.
	expectations, err := readFixtureExpectations(root)
	if err != nil {
		return fixtureExpectation{}, err
	}
	if len(expectations) == 0 {
		return fixtureExpectation{}, nil
	}
	name := filepath.Base(runDir)
	// Fixture expectations are keyed by run directory basename so the corpus can
	// move as a whole.
	return fixtureExpectation{ExpectedResult: expectations[name]}, nil
}

func readFixtureExpectations(root string) (map[string]string, error) {
	path := filepath.Join(root, "fixture-expectations.json")
	// Fixture expectations are optional metadata outside the verifier result.
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Missing expectations file means defaults, not a broken corpus.
			return nil, nil
		}
		return nil, err
	}
	var expectations map[string]string
	if err := json.Unmarshal(data, &expectations); err != nil {
		// Malformed expectation metadata is reported to the fixture validator.
		return nil, err
	}
	return expectations, nil
}

// flagSet is a tiny local parser for limited CLI needs.
type flagSet struct {
	name  string
	data  map[string]string
	bools map[string]bool
	args  []string
}

func (f *flagSet) setString(key, defaultValue string) {
	if f.data == nil {
		// Lazily allocate so tiny commands only register the flags they own.
		f.data = map[string]string{}
	}
	f.data[key] = defaultValue
}

func (f *flagSet) setBool(key string, defaultValue bool) {
	if f.bools == nil {
		// Boolean flags are tracked separately to reject string-only forms.
		f.bools = map[string]bool{}
	}
	f.bools[key] = defaultValue
}

func (f *flagSet) parse(args []string) error {
	rest := make([]string, 0)
	for i := 0; i < len(args); i++ {
		// The loop index is passed by pointer so string flags can consume their
		// following value without reparsing it as positional input.
		// consumeArg owns index advancement for flags with following values.
		done, err := f.consumeArg(args, &i, &rest)
		if err != nil {
			return err
		}
		if done {
			f.args = rest
			return nil
		}
	}
	f.args = rest
	return nil
}

func (f *flagSet) consumeArg(args []string, idx *int, rest *[]string) (bool, error) {
	arg := args[*idx]
	if arg == "--" {
		// Everything after -- is command payload, not parser flags.
		*rest = append(*rest, args[*idx+1:]...)
		return true, nil
	}
	if !strings.HasPrefix(arg, "--") {
		// Positional arguments are preserved for command-specific validation.
		*rest = append(*rest, arg)
		return false, nil
	}
	flag, flagValue, hasValue := splitFlag(arg)
	return false, f.consumeFlag(flag, flagValue, hasValue, args, idx)
}

func splitFlag(arg string) (string, string, bool) {
	parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
	if len(parts) == 1 {
		// Bare flags may either be boolean flags or string flags with next value.
		return parts[0], "", false
	}
	return parts[0], parts[1], true
}

func (f *flagSet) consumeFlag(flag string, flagValue string, hasValue bool, args []string, idx *int) error {
	isString, isBool := f.isKnownFlag(flag)
	if !isString && !isBool {
		// Unknown flags fail early before command code interprets inputs.
		return fmt.Errorf("unknown flag --%s", flag)
	}
	if hasValue {
		return f.consumeValue(flag, flagValue, isBool)
	}
	return f.consumeNoEqualsValue(flag, args, idx, isBool)
}

func (f *flagSet) isKnownFlag(flag string) (bool, bool) {
	// Return both flag classes so parsing can reject value syntax for booleans
	// without losing the unknown-flag distinction.
	_, isString := f.data[flag]
	_, isBool := f.bools[flag]
	return isString, isBool
}
func (f *flagSet) consumeValue(flag, flagValue string, isBool bool) error {
	if !isBool {
		// --flag=value is the direct string assignment form.
		f.data[flag] = flagValue
		return nil
	}
	return f.consumeBoolValue(flag, flagValue)
}

func (f *flagSet) consumeNoEqualsValue(flag string, args []string, idx *int, isBool bool) error {
	if !isBool {
		// String flags without equals consume the next argument as their value.
		return f.consumeStringFromNext(flag, args, idx)
	}
	nextIdx := *idx + 1
	if !isBoolValueAt(args, nextIdx) {
		// Bare boolean flags imply true unless followed by a boolean literal.
		f.bools[flag] = true
		return nil
	}
	*idx = nextIdx
	return f.consumeBoolValue(flag, args[*idx])
}

func (f *flagSet) consumeStringFromNext(flag string, args []string, idx *int) error {
	nextIdx := *idx + 1
	if nextIdx >= len(args) {
		// String flags must have a concrete following value.
		return fmt.Errorf("flag --%s requires value", flag)
	}
	value := args[nextIdx]
	if strings.HasPrefix(value, "--") {
		// Another flag cannot stand in for a missing string value.
		return fmt.Errorf("flag --%s requires value", flag)
	}
	*idx = nextIdx
	f.data[flag] = value
	return nil
}

func isBoolValueAt(args []string, idx int) bool {
	return idx < len(args) && isBoolLiteral(args[idx])
}

func (f *flagSet) consumeBoolValue(flag, flagValue string) error {
	switch strings.ToLower(flagValue) {
	case "false", "0":
		// Accept compact false spellings for generated command lines.
		f.bools[flag] = false
	case "true", "1", "":
		// Empty value covers --flag= and keeps legacy true semantics.
		f.bools[flag] = true
	default:
		// Invalid boolean values are usage errors, not ignored arguments.
		return fmt.Errorf("invalid boolean value for --%s: %s", flag, flagValue)
	}
	return nil
}

func (f *flagSet) stringValue(key string) string {
	if f.data == nil {
		// Unregistered string maps read as absent flags, matching parse defaults.
		return ""
	}
	return f.data[key]
}

func (f *flagSet) boolValue(key string) bool {
	if f.bools == nil {
		// Unregistered bool maps default to false instead of implying a flag.
		return false
	}
	return f.bools[key]
}

func (f *flagSet) rest() []string {
	return f.args
}

func isHelp(args []string) bool {
	return len(args) == 1 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help")
}

func isBoolLiteral(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "false" || lower == "1" || lower == "0"
}
