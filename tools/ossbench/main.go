package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	cfg, err := parseFlagsAndArgs(args, stderr)
	if err != nil {
		return 2
	}
	if cfg.list {
		return handleList(stdout, stderr, cfg.args)
	}
	return executeBenchmarks(cfg, stdout, stderr)
}

type runConfig struct {
	asJSON     bool
	iterations int
	list       bool
	name       string
	raw        bool
	args       []string
}

func parseFlagsAndArgs(args []string, stderr io.Writer) (runConfig, error) {
	fs := flag.NewFlagSet("ossbench", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: ossbench [flags] [command args...]\n")
		fmt.Fprintf(stderr, "Run built-in or custom benchmarks with min/max/median stats.\n")
		fmt.Fprintf(stderr, "Use -list to see built-in benchmarks.\n\n")
		fs.PrintDefaults()
	}
	var cfg runConfig
	fs.BoolVar(&cfg.asJSON, "json", false, "emit JSON output")
	fs.IntVar(&cfg.iterations, "n", 20, "number of iterations")
	fs.BoolVar(&cfg.list, "list", false, "list built-in benchmarks")
	fs.StringVar(&cfg.name, "bench", "", "run a single built-in benchmark by name")
	fs.BoolVar(&cfg.raw, "raw", false, "include all_ms in JSON output")
	if err := fs.Parse(args); err != nil {
		return runConfig{}, err
	}
	cfg.args = fs.Args()
	return cfg, nil
}

func handleList(stdout, stderr io.Writer, args []string) int {
	if len(args) > 0 {
		fmt.Fprintf(stderr, "unexpected positional args with -list: %s\n", strings.Join(args, " "))
		return 2
	}
	for _, b := range builtIns {
		fmt.Fprintf(stdout, "%s\t%s\n", b.Name, b.Description)
	}
	return 0
}

func executeBenchmarks(cfg runConfig, stdout, stderr io.Writer) int {
	defs, cleanup, err := resolveBenchmarkDefs(cfg)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		return 2
	}
	if cleanup != nil {
		defer cleanup()
	}
	results, err := runAllBenchmarks(defs, cfg.iterations, cfg.raw)
	if err != nil {
		fmt.Fprintf(stderr, "%v\n", err)
		return 2
	}
	return printAndExit(stdout, stderr, results, cfg.asJSON)
}

func resolveBenchmarkDefs(cfg runConfig) ([]benchmarkDef, func(), error) {
	if cfg.name != "" {
		if len(cfg.args) > 0 {
			return nil, nil, fmt.Errorf("unexpected positional args with -bench: %s", strings.Join(cfg.args, " "))
		}
		return resolveSingleBuiltin(cfg.name)
	}
	if len(cfg.args) > 0 {
		defs := []benchmarkDef{{
			Name:        strings.Join(cfg.args, " "),
			Description: "custom command",
			Cmd:         cfg.args[0],
			Args:        cfg.args[1:],
		}}
		return defs, nil, nil
	}
	return resolveAllBuiltins()
}

func findBuiltin(name string) (benchmarkDef, bool) {
	for _, b := range builtIns {
		if b.Name == name {
			return b, true
		}
	}
	return benchmarkDef{}, false
}

func resolveSingleBuiltin(name string) ([]benchmarkDef, func(), error) {
	def, ok := findBuiltin(name)
	if !ok {
		return nil, nil, fmt.Errorf("unknown benchmark: %s", name)
	}
	if err := resolveBuiltIns(); err != nil {
		return nil, nil, fmt.Errorf("resolve built-ins: %w", err)
	}
	def.Cmd = tempBinaryPath
	def.Source = "temp-build"
	return []benchmarkDef{def}, cleanupTempBinary, nil
}

func resolveAllBuiltins() ([]benchmarkDef, func(), error) {
	if err := resolveBuiltIns(); err != nil {
		return nil, nil, fmt.Errorf("resolve built-ins: %w", err)
	}
	return builtIns, cleanupTempBinary, nil
}

func runAllBenchmarks(defs []benchmarkDef, iterations int, raw bool) ([]benchmarkResult, error) {
	results := make([]benchmarkResult, 0, len(defs))
	for i := range defs {
		d := &defs[i]
		if err := setupWrap(d); err != nil {
			return nil, err
		}
		res := runBenchmark(*d, iterations)
		res = finalizeResult(res, d.Cleanup, raw)
		results = append(results, res)
	}
	return results, nil
}

func setupWrap(def *benchmarkDef) error {
	if def.Name != "sdp-trace-wrap" {
		return nil
	}
	tmpDir, err := os.MkdirTemp("", "ossbench-wrap-*")
	if err != nil {
		return fmt.Errorf("mkdir temp: %w", err)
	}
	def.Dir = tmpDir
	def.Cleanup = func() {
		_ = os.RemoveAll(tmpDir)
	}
	return nil
}

func finalizeResult(res benchmarkResult, cleanup func(), raw bool) benchmarkResult {
	if cleanup != nil {
		cleanup()
	}
	if !raw {
		res.AllMs = nil
	}
	return res
}

func printAndExit(stdout, stderr io.Writer, results []benchmarkResult, asJSON bool) int {
	if err := printResults(stdout, results, asJSON); err != nil {
		fmt.Fprintf(stderr, "print results: %v\n", err)
		return 2
	}
	return exitCode(results)
}

func printResults(w io.Writer, results []benchmarkResult, asJSON bool) error {
	if asJSON {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(results)
	}
	width := maxNameWidth(results)
	for _, r := range results {
		if _, err := fmt.Fprint(w, formatResultLine(r, width)); err != nil {
			return err
		}
	}
	return nil
}

func formatResultLine(r benchmarkResult, width int) string {
	if r.Error != "" {
		return fmt.Sprintf("%*s  error: %s  median=%6.2f ms  min=%6.2f ms  max=%6.2f ms  attempted=%d succeeded=%d\n",
			-width, r.Name, r.Error, r.MedianMs, r.MinMs, r.MaxMs, r.AttemptedIterations, r.SucceededIterations)
	}
	return fmt.Sprintf("%*s  median=%6.2f ms  min=%6.2f ms  max=%6.2f ms  attempted=%d succeeded=%d\n",
		-width, r.Name, r.MedianMs, r.MinMs, r.MaxMs, r.AttemptedIterations, r.SucceededIterations)
}

func maxNameWidth(results []benchmarkResult) int {
	maxW := 24
	for _, r := range results {
		if len(r.Name) > maxW {
			maxW = len(r.Name)
		}
	}
	return maxW
}

func exitCode(results []benchmarkResult) int {
	for _, r := range results {
		if r.Error != "" {
			return 1
		}
	}
	return 0
}

// builtIns are the standard OSS tool benchmarks.
var builtIns = []benchmarkDef{
	{
		Name:        "sdp-trace-version",
		Description: "sdp-trace version command",
		Cmd:         "sdp-trace",
		Args:        []string{"version"},
		Source:      "PATH",
	},
	{
		Name:        "sdp-trace-wrap",
		Description: "sdp-trace wrap no-op",
		Cmd:         "sdp-trace",
		Args: func() []string {
			if runtime.GOOS == "windows" {
				return []string{"wrap", "cmd", "/c", "exit", "0"}
			}
			return []string{"wrap", "true"}
		}(),
		Source: "PATH",
	},
}

// builtInsOrig holds the original Cmd and Source values so resolveBuiltIns
// can mutate the global slice and cleanupTempBinary can restore it.
var builtInsOrig = make([]benchmarkDef, len(builtIns))

func init() {
	copy(builtInsOrig, builtIns)
}

// tempBinaryPath is set when the harness builds sdp-trace into a temp dir.
var tempBinaryPath string

// resolveBuiltIns builds the sdp-trace binary from current source into a temp
// dir and updates builtIns. It never uses a pre-existing repo-root binary so
// that benchmark results always reflect the checked-out source.
func resolveBuiltIns() error {
	tmpDir, err := os.MkdirTemp("", "ossbench-bin-*")
	if err != nil {
		return fmt.Errorf("mkdir temp for build: %w", err)
	}
	bin := filepath.Join(tmpDir, "sdp-trace")
	if err := buildBinary(bin); err != nil {
		_ = os.RemoveAll(tmpDir)
		return fmt.Errorf("sdp-trace build failed: %w", err)
	}
	tempBinaryPath = bin
	for i := range builtIns {
		builtIns[i].Cmd = bin
		builtIns[i].Source = "temp-build"
	}
	return nil
}

// cleanupTempBinary removes the temp-built binary if one was created and
// restores builtIns to their original values so subsequent runs in the same
// process do not reference a deleted path.
func cleanupTempBinary() {
	if tempBinaryPath != "" {
		_ = os.RemoveAll(filepath.Dir(tempBinaryPath))
		tempBinaryPath = ""
	}
	copy(builtIns, builtInsOrig)
}
