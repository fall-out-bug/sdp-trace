package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fall_out_bug/sdp-trace/internal/demo"
	"github.com/fall_out_bug/sdp-trace/internal/query"
	"github.com/fall_out_bug/sdp-trace/internal/recorder"
	"github.com/fall_out_bug/sdp-trace/internal/trace"
	"github.com/fall_out_bug/sdp-trace/internal/verifier"
	"github.com/fall_out_bug/sdp-trace/internal/witness"
)

const (
	exitUsage        = 2
	exitCannotVerify = 3
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printUsage(stdout)
		return 0
	}
	cmd := args[0]
	cmdArgs := args[1:]
	ctx := context.Background()

	switch cmd {
	case "wrap":
		return runWrap(ctx, cmdArgs, stdout, stderr)
	case "run":
		return runWrappedCommand(ctx, cmdArgs, stdout, stderr)
	case "dry-run":
		return runDryRun(ctx, cmdArgs, stdout, stderr)
	case "preview":
		return runPreview(ctx, cmdArgs, stdout, stderr)
	case "doctor":
		return runDoctor(ctx, cmdArgs, stdout, stderr)
	case "verify":
		return runVerify(ctx, cmdArgs, stdout, stderr)
	case "explain":
		return runExplain(ctx, cmdArgs, stdout, stderr)
	case "query":
		return runQuery(ctx, cmdArgs, stdout, stderr)
	case "report":
		return runReport(ctx, cmdArgs, stdout, stderr)
	case "gate":
		return runGate(ctx, cmdArgs, stdout, stderr)
	case "override":
		return runOverride(ctx, cmdArgs, stdout, stderr)
	case "witness":
		return runWitness(ctx, cmdArgs, stdout, stderr)
	case "validate-fixtures":
		return runValidateFixtures(ctx, cmdArgs, stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", cmd)
		printUsage(stderr)
		return 1
	}
}

func runReport(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "report"}
	opts.setString("out", "")
	opts.setString("contract", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	targets := opts.rest()
	if len(targets) != 1 {
		fmt.Fprintln(stderr, "report requires <runs-root-or-run-dir>")
		return exitUsage
	}
	outDir := opts.stringValue("out")
	if outDir == "" {
		fmt.Fprintln(stderr, "report requires --out <dir>")
		return exitUsage
	}
	artifacts, err := demo.WriteReport(targets[0], outDir, opts.stringValue("contract"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(artifacts.Summary, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	if artifacts.Summary.CannotVerifyCount > 0 {
		return exitCannotVerify
	}
	if artifacts.Summary.FailedCount > 0 {
		return 1
	}
	return 0
}

func runGate(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) > 0 {
		switch args[0] {
		case "explain":
			return runGateExplain(args[1:], stdout, stderr)
		case "preview":
			return runGatePreview(args[1:], stdout, stderr)
		}
	}
	opts := &flagSet{name: "gate"}
	opts.setString("out", "")
	opts.setString("contract", "")
	opts.setString("witness", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	targets := opts.rest()
	if len(targets) != 1 {
		fmt.Fprintln(stderr, "gate requires <runs-root-or-run-dir>")
		return exitUsage
	}
	outPath := opts.stringValue("out")
	if outPath == "" {
		fmt.Fprintln(stderr, "gate requires --out <file>")
		return exitUsage
	}
	result, err := demo.WriteGate(targets[0], outPath, opts.stringValue("contract"), opts.stringValue("witness"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	return gateExitCode(result)
}

func runGateExplain(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "gate explain"}
	opts.setString("gate-result", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "gate explain accepts only flags")
		return exitUsage
	}
	path := opts.stringValue("gate-result")
	if path == "" {
		fmt.Fprintln(stderr, "gate explain requires --gate-result <file>")
		return exitUsage
	}
	var result demo.GateResult
	if err := readJSONFile(path, &result); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "Gate mode: %s\n", result.GateMode)
	fmt.Fprintf(stdout, "Trust cap: %s\n", result.TrustCap)
	fmt.Fprintf(stdout, "Local gate: %s\n", result.LocalGate)
	fmt.Fprintf(stdout, "CI witness gate: %s\n", result.CIWitnessGate)
	fmt.Fprintf(stdout, "Audit-grade gate: %s\n", result.AuditGradeGate)
	for _, requiredRun := range result.RequiredRuns {
		fmt.Fprintf(stdout, "Required run %s: %s\n", requiredRun.ID, requiredRun.State)
	}
	for _, binding := range result.WitnessBindings {
		fmt.Fprintf(stdout, "Witness binding %s: %s\n", binding.ID, binding.State)
	}
	for _, missing := range result.MissingAuditEvidence {
		fmt.Fprintf(stdout, "Missing audit evidence: %s\n", missing)
	}
	for _, override := range result.OverrideRequests {
		fmt.Fprintf(stdout, "Override %s: %s\n", override.OverrideID, override.State)
	}
	for _, reason := range result.Reasons {
		fmt.Fprintf(stdout, "Reason: %s\n", reason)
	}
	for _, action := range result.NextActions {
		fmt.Fprintf(stdout, "Next action: %s\n", action)
	}
	return 0
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

func runGatePreview(args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "gate preview"}
	opts.setString("contract", "")
	opts.setString("witness", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	targets := opts.rest()
	if len(targets) != 1 {
		fmt.Fprintln(stderr, "gate preview requires <runs-root-or-run-dir>")
		return exitUsage
	}
	contract, err := trace.LoadContract(opts.stringValue("contract"))
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	report := gatePreviewReport{
		Command:          "gate preview",
		GateMode:         previewGateMode(contract),
		TrustCap:         string(trace.TrustScopeLocalObserved),
		RequiredRuns:     requiredRunIDs(contract),
		RequiredEvidence: requiredEvidenceIDsForCLI(contract),
		Claim:            "preview is read-only and does not claim the gate will pass",
	}
	witnessPath := opts.stringValue("witness")
	if witnessPath != "" {
		report.WitnessInspectable, report.WitnessMismatches = demo.PreviewWitnessBinding(witnessPath, targets[0])
	}
	payload, _ := json.MarshalIndent(report, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	_ = targets[0]
	return 0
}

func runOverride(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 || args[0] != "request" {
		fmt.Fprintln(stderr, "override requires request")
		return exitUsage
	}
	opts := &flagSet{name: "override request"}
	opts.setString("out", "")
	opts.setString("id", "")
	opts.setString("by", "")
	opts.setString("reason", "")
	opts.setString("source-ref", "")
	opts.setString("scope", "")
	opts.setString("external-reference", "")
	if err := opts.parse(args[1:]); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "override request accepts only flags")
		return exitUsage
	}
	runDir := opts.stringValue("out")
	required := map[string]string{
		"--out":        runDir,
		"--id":         opts.stringValue("id"),
		"--by":         opts.stringValue("by"),
		"--reason":     opts.stringValue("reason"),
		"--source-ref": opts.stringValue("source-ref"),
		"--scope":      opts.stringValue("scope"),
	}
	for flag, value := range required {
		if strings.TrimSpace(value) == "" {
			fmt.Fprintf(stderr, "override request requires %s\n", flag)
			return exitUsage
		}
	}
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
		payload["external_reference"] = external
	}
	event, err := trace.AppendRunEvent(runDir, trace.EventPolicyOverrideRequested, payload, "sdp-trace-cli")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "override_event: %s\n", event.EventID)
	return 0
}

func readJSONFile(path string, dst any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dst)
}

func previewGateMode(contract trace.Contract) string {
	mode := demo.GateModeObservation
	for _, required := range contract.RequiredRuns {
		switch required.Profile {
		case demo.GateModeProtectedFuture:
			return demo.GateModeProtectedFuture
		case demo.GateModeAdvisoryCI:
			mode = demo.GateModeAdvisoryCI
		}
	}
	return mode
}

func requiredRunIDs(contract trace.Contract) []string {
	ids := make([]string, 0, len(contract.RequiredRuns))
	for _, required := range contract.RequiredRuns {
		if required.ID != "" {
			ids = append(ids, required.ID)
		}
	}
	return ids
}

func requiredEvidenceIDsForCLI(contract trace.Contract) []string {
	ids := make([]string, 0, len(contract.RequiredEvidence))
	for _, requirement := range contract.RequiredEvidence {
		if requirement.ID != "" {
			ids = append(ids, requirement.ID)
		}
	}
	return ids
}

func gateExitCode(result demo.GateResult) int {
	for _, requiredRun := range result.RequiredRuns {
		if requiredRun.State == demo.GateFail || requiredRun.State == demo.GateMissingTelemetry {
			return 1
		}
	}
	for _, state := range []string{result.LocalGate, result.CIWitnessGate, result.AuditGradeGate} {
		if state == demo.GateFail || state == demo.GateMissingTelemetry {
			return 1
		}
	}
	for _, requiredRun := range result.RequiredRuns {
		if requiredRun.State == demo.GateCannotVerify {
			return exitCannotVerify
		}
	}
	for _, state := range []string{result.LocalGate, result.CIWitnessGate, result.AuditGradeGate} {
		if state == demo.GateCannotVerify {
			return exitCannotVerify
		}
	}
	return 0
}

func runWitness(_ context.Context, args []string, stdout, stderr io.Writer) int {
	opts := &flagSet{name: "witness"}
	opts.setString("kind", "")
	opts.setString("out", "")
	opts.setString("report-dir", "")
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	targets := opts.rest()
	if len(targets) != 1 {
		fmt.Fprintln(stderr, "witness requires <runs-root-or-run-dir>")
		return exitUsage
	}
	if opts.stringValue("kind") != witness.KindGitHubActions {
		fmt.Fprintln(stderr, "witness requires --kind github-actions")
		return exitUsage
	}
	outPath := opts.stringValue("out")
	if outPath == "" {
		fmt.Fprintln(stderr, "witness requires --out <file>")
		return exitUsage
	}
	record, err := witness.WriteGitHubActions(outPath, targets[0], opts.stringValue("report-dir"), witness.EnvironmentFromOS())
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	payload, _ := json.MarshalIndent(record, "", "  ")
	fmt.Fprintf(stdout, "%s\n", payload)
	if record.Status == witness.StatusCannotVerify {
		return exitCannotVerify
	}
	return 0
}

func runWrap(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts := &flagSet{name: "wrap"}
	opts.setString("name", "")
	opts.setString("contract", "")
	opts.setString("output-dir", "")
	if err := opts.parse(args); err != nil {
		return exitUsage
	}
	command := opts.rest()
	if len(command) == 0 {
		fmt.Fprintln(stderr, "wrap requires a command")
		return exitUsage
	}
	res, err := recorder.Run(ctx, recorder.RecorderOptions{
		WrapperName:        opts.stringValue("name"),
		ContractPath:       opts.stringValue("contract"),
		UseDefaultContract: true,
		OutputDir:          opts.stringValue("output-dir"),
		Command:            command,
	})
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	fmt.Fprintf(stdout, "run_dir: %s\n", res.RunDir)
	return res.ExitCode
}

func runWrappedCommand(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts := &flagSet{name: "run"}
	opts.setString("task", "")
	opts.setString("contract", "")
	opts.setBool("use-default-contract", false)
	opts.setString("name", "")
	opts.setString("output-dir", "")
	if err := opts.parse(args); err != nil {
		return exitUsage
	}
	command := opts.rest()
	if len(command) == 0 {
		fmt.Fprintln(stderr, "run requires a command")
		return exitUsage
	}
	useDefault := opts.boolValue("use-default-contract")
	task := opts.stringValue("task")
	if task == "" {
		fmt.Fprintln(stderr, "run requires --task")
		return exitUsage
	}
	if opts.stringValue("contract") == "" && !useDefault {
		fmt.Fprintln(stderr, "run requires --contract unless --use-default-contract is set")
		return exitUsage
	}
	res, err := recorder.Run(ctx, recorder.RecorderOptions{
		Task:               task,
		WrapperName:        opts.stringValue("name"),
		ContractPath:       opts.stringValue("contract"),
		UseDefaultContract: useDefault,
		OutputDir:          opts.stringValue("output-dir"),
		Command:            command,
	})
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}
	fmt.Fprintf(stdout, "run_dir: %s\n", res.RunDir)
	return res.ExitCode
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
	opts := &flagSet{name: commandName}
	opts.setString("contract", "")
	opts.setBool("use-default-contract", true)
	opts.setString("name", "")
	if err := opts.parse(args); err != nil {
		return exitUsage
	}
	command := opts.rest()
	if len(command) == 0 {
		fmt.Fprintf(stderr, "%s requires a command\n", commandName)
		return exitUsage
	}
	contractPath := opts.stringValue("contract")
	useDefault := opts.boolValue("use-default-contract")
	if contractPath == "" && !useDefault {
		fmt.Fprintf(stderr, "%s requires --contract unless --use-default-contract is set\n", commandName)
		return exitUsage
	}
	contract := trace.DefaultContract
	if contractPath != "" {
		loaded, err := trace.LoadContract(contractPath)
		if err != nil {
			fmt.Fprintf(stderr, "failed to load contract: %v\n", err)
			return exitCannotVerify
		}
		contract = loaded
	}
	payload := map[string]any{
		"mode":                 mode,
		"command_descriptor":   trace.NewCommandDescriptor(command),
		"contract":             contract,
		"boundaries":           previewBoundaries(),
		"offline_implications": previewOfflineImplications(),
		"writes_artifacts":     false,
		"safe_retention_modes": safeRetentionModes(),
		"warning":              "no run artifacts were written",
	}
	data, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintf(stdout, "%s\n", data)
	return 0
}

func runDoctor(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts := &flagSet{name: "doctor"}
	opts.setString("contract", "")
	opts.setString("output-dir", defaultRunRoot)
	opts.setString("report-dir", defaultReportDir)
	if err := opts.parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	if len(opts.rest()) != 0 {
		fmt.Fprintln(stderr, "doctor does not accept positional arguments")
		return exitUsage
	}
	report, exitCode := buildDoctorReport(doctorOptions{
		ContractPath: opts.stringValue("contract"),
		OutputDir:    opts.stringValue("output-dir"),
		ReportDir:    opts.stringValue("report-dir"),
		Env:          witness.EnvironmentFromOS(),
	})
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	fmt.Fprintf(stdout, "%s\n", data)
	return exitCode
}

func runVerify(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "verify requires <run-dir>")
		return exitUsage
	}
	runDir := args[0]
	if info, err := os.Stat(runDir); err != nil || !info.IsDir() {
		fmt.Fprintf(stderr, "run directory does not exist: %s\n", runDir)
		return exitCannotVerify
	}
	result, table, audit, err := verifier.VerifyRun(runDir)
	if writeErr := verifier.WriteVerifierArtifacts(runDir, result, table, audit); writeErr != nil {
		fmt.Fprintf(stderr, "failed writing verifier artifacts for %s: %v\n", runDir, writeErr)
		return 1
	}
	data, _ := json.MarshalIndent(result, "", "  ")
	fmt.Fprintf(stdout, "%s\n", data)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
	}
	switch result.Result {
	case trace.VerdictObserved, trace.VerdictNotAssessed:
		return 0
	case trace.VerdictFail:
		return 1
	case trace.VerdictCannotVerify:
		return exitCannotVerify
	default:
		return 0
	}
}

func runExplain(_ context.Context, args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "explain requires <run-dir>")
		return exitUsage
	}
	runDir := args[0]
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
	opts.setString("query", "")
	if err := opts.parse(args); err != nil {
		return exitUsage
	}
	queryName := opts.stringValue("query")
	runDirs := opts.rest()
	if len(runDirs) == 0 {
		fmt.Fprintln(stderr, "query requires <run-dir>")
		return exitUsage
	}
	if queryName != query.QueryMissingEvidence {
		fmt.Fprintf(stderr, "unsupported query: %s\n", queryName)
		return exitUsage
	}
	payload, err := query.MissingEvidence(runDirs[0])
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitCannotVerify
	}
	fmt.Fprintf(stdout, "%s\n", payload)
	return 0
}

func runValidateFixtures(_ context.Context, args []string, stdout, stderr io.Writer) int {
	fixtureRoot := "."
	if len(args) > 0 {
		fixtureRoot = args[0]
	}
	runDirs, err := demo.DiscoverRunDirs(fixtureRoot)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	failed := false
	for _, runDir := range runDirs {
		result, table, audit, verifyErr := verifier.VerifyRun(runDir)
		if err := verifier.WriteVerifierArtifacts(runDir, result, table, audit); err != nil {
			fmt.Fprintf(stderr, "failed writing verifier artifacts for %s: %v\n", runDir, err)
			failed = true
			continue
		}
		fmt.Fprintf(stdout, "%s => %s\n", runDir, result.Result)
		if verifyErr != nil {
			fmt.Fprintf(stderr, "%s verification error: %v\n", runDir, verifyErr)
		}
		expectation, err := readFixtureExpectation(fixtureRoot, runDir)
		if err != nil {
			fmt.Fprintf(stderr, "invalid fixture expectation for %s: %v\n", runDir, err)
			failed = true
			continue
		}
		if expectation.ExpectedResult != "" && expectation.ExpectedResult != string(result.Result) {
			fmt.Fprintf(stderr, "%s expected %s, got %s\n", runDir, expectation.ExpectedResult, result.Result)
			failed = true
			continue
		}
		if expectation.ExpectedResult == "" && result.Result == trace.VerdictFail {
			failed = true
		}
		if expectation.ExpectedResult == "" && result.Result == trace.VerdictCannotVerify {
			failed = true
		}
	}
	if failed {
		return 1
	}
	return 0
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
	contract := defaultContract
	contractCheck := doctorCheck{
		ID:        "contract",
		State:     "pass",
		Reason:    "default contract is available",
		Contract:  defaultContract.ContractID,
		Reference: "local-default-v1",
	}
	if options.ContractPath != "" {
		loaded, err := trace.LoadContract(options.ContractPath)
		if err != nil {
			result = string(trace.VerdictCannotVerify)
			exitCode = exitCannotVerify
			contractCheck = doctorCheck{
				ID:        "contract",
				State:     string(trace.VerdictCannotVerify),
				Reason:    "contract cannot be loaded",
				Reference: options.ContractPath,
			}
		} else {
			contract = loaded
			contractCheck = doctorCheck{
				ID:        "contract",
				State:     "pass",
				Reason:    "contract can be loaded",
				Contract:  contract.ContractID,
				Reference: options.ContractPath,
			}
		}
	}
	ciCheck := ciWitnessPrerequisiteCheck(options.Env)
	outputDirCheck := writablePathCheck("output_directory", options.OutputDir, "run artifact output directory is writable")
	reportDirCheck := writablePathCheck("report_directory", options.ReportDir, "report artifact directory is writable")
	expectedEvidenceCheck := expectedEvidenceReferenceCheck(contract)
	for _, check := range []doctorCheck{outputDirCheck, reportDirCheck, expectedEvidenceCheck} {
		if check.State == string(trace.VerdictCannotVerify) {
			result = string(trace.VerdictCannotVerify)
			exitCode = exitCannotVerify
		}
	}
	report := doctorReport{
		Command: "doctor",
		Result:  result,
		Environment: []doctorCheck{
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
		},
		ControlPoints: []doctorCheck{
			{
				ID:     "local_wrapper",
				State:  "pass",
				Reason: "wrap and run commands are registered in this binary",
			},
			outputDirCheck,
			reportDirCheck,
			contractCheck,
			expectedEvidenceCheck,
			{
				ID:        "default_contract",
				State:     "pass",
				Reason:    "built-in contract is available for local development",
				Contract:  defaultContract.ContractID,
				Reference: defaultContract.Version,
			},
			ciCheck,
		},
		SafeRetentionModes: safeRetentionModes(),
	}
	return report, exitCode
}

func writablePathCheck(id, path, okReason string) doctorCheck {
	if strings.TrimSpace(path) == "" {
		return doctorCheck{
			ID:     id,
			State:  string(trace.VerdictCannotVerify),
			Reason: "path is empty",
		}
	}
	target := path
	info, err := os.Stat(target)
	if err == nil && !info.IsDir() {
		return doctorCheck{
			ID:        id,
			State:     string(trace.VerdictCannotVerify),
			Reason:    "path exists but is not a directory",
			Reference: path,
		}
	}
	if os.IsNotExist(err) {
		target = filepath.Dir(path)
		if target == "" {
			target = "."
		}
	}
	probe, err := os.CreateTemp(target, ".sdp-trace-doctor-")
	if err != nil {
		return doctorCheck{
			ID:        id,
			State:     string(trace.VerdictCannotVerify),
			Reason:    "directory is not writable",
			Reference: path,
		}
	}
	probeName := probe.Name()
	_ = probe.Close()
	_ = os.Remove(probeName)
	return doctorCheck{
		ID:        id,
		State:     "pass",
		Reason:    okReason,
		Reference: path,
	}
}

func expectedEvidenceReferenceCheck(contract trace.Contract) doctorCheck {
	if len(contract.RequiredEvents) == 0 {
		return doctorCheck{
			ID:       "expected_evidence_references",
			State:    string(trace.VerdictCannotVerify),
			Reason:   "contract has no required_events",
			Contract: contract.ContractID,
		}
	}
	missing := make([]string, 0)
	for _, eventType := range contract.RequiredEvents {
		if !knownEventType(eventType) {
			missing = append(missing, "required_events:"+eventType)
		}
	}
	for _, evidence := range contract.RequiredEvidence {
		if strings.TrimSpace(evidence.ID) == "" {
			missing = append(missing, "required_evidence:<missing_id>")
		}
		if strings.TrimSpace(evidence.EventType) == "" {
			missing = append(missing, "required_evidence:"+evidence.ID+":<missing_event_type>")
			continue
		}
		if !knownEventType(evidence.EventType) {
			missing = append(missing, "required_evidence:"+evidence.ID+":"+evidence.EventType)
		}
	}
	if len(missing) > 0 {
		return doctorCheck{
			ID:       "expected_evidence_references",
			State:    string(trace.VerdictCannotVerify),
			Reason:   "contract references unsupported event types",
			Contract: contract.ContractID,
			Missing:  missing,
		}
	}
	return doctorCheck{
		ID:       "expected_evidence_references",
		State:    "pass",
		Reason:   "contract required events and evidence references are supported by the current local event model",
		Contract: contract.ContractID,
	}
}

func knownEventType(eventType string) bool {
	switch trace.EventType(eventType) {
	case trace.EventRecorderAttached,
		trace.EventRunStarted,
		trace.EventCommandStarted,
		trace.EventCommandFinished,
		trace.EventRunClosed,
		trace.EventPolicyOverrideRequested:
		return true
	default:
		return false
	}
}

func ciWitnessPrerequisiteCheck(env map[string]string) doctorCheck {
	missing := missingCIWitnessFields(env)
	if len(missing) > 0 {
		return doctorCheck{
			ID:      "ci_witness_prerequisites",
			State:   string(trace.VerdictCannotVerify),
			Reason:  "GitHub Actions identity or OIDC prerequisite is unavailable in this environment",
			Missing: missing,
		}
	}
	return doctorCheck{
		ID:     "ci_witness_prerequisites",
		State:  "pass",
		Reason: "GitHub Actions identity and OIDC prerequisites are present",
	}
}

func missingCIWitnessFields(env map[string]string) []string {
	required := []string{
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
	missing := make([]string, 0)
	for _, key := range required {
		if strings.TrimSpace(env[key]) == "" {
			missing = append(missing, key)
		}
	}
	if env["GITHUB_ACTIONS"] != "" && env["GITHUB_ACTIONS"] != "true" {
		missing = append(missing, "GITHUB_ACTIONS=true")
	}
	return missing
}

func safeRetentionModes() []string {
	return []string{
		string(trace.RetentionModeDigestOnly),
		string(trace.RetentionModeSanitizedExcerpt),
		string(trace.RetentionModeEncryptedRawRef),
		string(trace.RetentionModeExternalArtifactRef),
		string(trace.RetentionModeNotAssessed),
	}
}

func previewBoundaries() []previewBoundary {
	return []previewBoundary{
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
			Boundary: string(trace.ObservationBoundaryExternalWitness),
			State:    string(trace.ObservationStateNotIntegrated),
			Reason:   "external witness profile is not implemented in Block 13B",
		},
	}
}

func previewOfflineImplications() []previewOfflineImplication {
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

func printUsage(w io.Writer) {
	const usage = `sdp-trace local recorder and verifier commands.

Usage:
  sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>
  sdp-trace run --task <task-ref> [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace dry-run [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace preview [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace doctor [--contract <file>]
  sdp-trace verify <run-dir>
  sdp-trace explain <run-dir>
  sdp-trace query --query missing-evidence <run-dir>
  sdp-trace report --out <dir> <runs-root-or-run-dir>
  sdp-trace gate --out <file> <runs-root-or-run-dir>
  sdp-trace witness --kind github-actions --out <file> [--report-dir <dir>] <runs-root-or-run-dir>
  sdp-trace validate-fixtures [root-dir]
`
	fmt.Fprint(w, usage)
}

type fixtureExpectation struct {
	ExpectedResult string `json:"expected_result"`
}

func readFixtureExpectation(root, runDir string) (fixtureExpectation, error) {
	expectations, err := readFixtureExpectations(root)
	if err != nil {
		return fixtureExpectation{}, err
	}
	if len(expectations) == 0 {
		return fixtureExpectation{}, nil
	}
	name := filepath.Base(runDir)
	return fixtureExpectation{ExpectedResult: expectations[name]}, nil
}

func readFixtureExpectations(root string) (map[string]string, error) {
	path := filepath.Join(root, "fixture-expectations.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var expectations map[string]string
	if err := json.Unmarshal(data, &expectations); err != nil {
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
		f.data = map[string]string{}
	}
	f.data[key] = defaultValue
}

func (f *flagSet) setBool(key string, defaultValue bool) {
	if f.bools == nil {
		f.bools = map[string]bool{}
	}
	f.bools[key] = defaultValue
}

func (f *flagSet) parse(args []string) error {
	rest := make([]string, 0)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			rest = append(rest, args[i+1:]...)
			break
		}
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(strings.TrimPrefix(arg, "--"), "=", 2)
			key := parts[0]
			_, knownString := f.data[key]
			_, knownBool := f.bools[key]
			if !knownString && !knownBool {
				return fmt.Errorf("unknown flag --%s", key)
			}
			switch len(parts) {
			case 1:
				switch {
				case knownBool:
					if i+1 < len(args) && isBoolLiteral(args[i+1]) {
						f.bools[key] = parseBoolLiteral(args[i+1])
						i++
					} else {
						f.bools[key] = true
					}
				default:
					if i+1 >= len(args) {
						return fmt.Errorf("flag --%s requires value", key)
					}
					val := args[i+1]
					if strings.HasPrefix(val, "--") {
						return fmt.Errorf("flag --%s requires value", key)
					}
					i++
					f.data[key] = val
				}
			case 2:
				if _, ok := f.bools[key]; ok {
					lower := strings.ToLower(parts[1])
					if lower == "false" || lower == "0" {
						f.bools[key] = false
					} else if lower == "true" || lower == "1" || lower == "" {
						f.bools[key] = true
					} else {
						return fmt.Errorf("invalid boolean value for --%s: %s", key, parts[1])
					}
					continue
				}
				f.data[key] = parts[1]
			default:
			}
			continue
		}
		rest = append(rest, arg)
	}
	f.args = rest
	return nil
}

func (f *flagSet) stringValue(key string) string {
	if f.data == nil {
		return ""
	}
	return f.data[key]
}

func (f *flagSet) boolValue(key string) bool {
	if f.bools == nil {
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

func parseBoolLiteral(value string) bool {
	lower := strings.ToLower(value)
	return lower == "true" || lower == "1"
}
