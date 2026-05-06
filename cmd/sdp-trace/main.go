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
	if result.LocalGate != demo.GatePass {
		return 1
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
	if isHelp(args) {
		printUsage(stdout)
		return 0
	}
	opts := &flagSet{name: "dry-run"}
	opts.setString("contract", "")
	opts.setBool("use-default-contract", true)
	opts.setString("name", "")
	if err := opts.parse(args); err != nil {
		return exitUsage
	}
	command := opts.rest()
	if len(command) == 0 {
		fmt.Fprintln(stderr, "dry-run requires a command")
		return exitUsage
	}
	contractPath := opts.stringValue("contract")
	useDefault := opts.boolValue("use-default-contract")
	if contractPath == "" && !useDefault {
		fmt.Fprintln(stderr, "dry-run requires --contract unless --use-default-contract is set")
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
		"mode":     "simulation",
		"command":  command,
		"contract": contract,
		"warning":  "no run artifacts were written",
	}
	data, _ := json.MarshalIndent(payload, "", "  ")
	fmt.Fprintf(stdout, "%s\n", data)
	return 0
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

func printUsage(w io.Writer) {
	const usage = `sdp-trace local recorder and verifier commands.

Usage:
  sdp-trace wrap --name <name> [--contract <file>] [--output-dir <dir>] -- <command...>
  sdp-trace run --task <task-ref> [--contract <file> | --use-default-contract] -- <command...>
  sdp-trace dry-run [--contract <file> | --use-default-contract] -- <command...>
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
