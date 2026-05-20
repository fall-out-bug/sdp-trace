package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("ossbench", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: ossbench [flags] [command args...]\n")
		fmt.Fprintf(stderr, "Run built-in or custom benchmarks with min/max/median stats.\n")
		fmt.Fprintf(stderr, "Use -list to see built-in benchmarks.\n\n")
		fs.PrintDefaults()
	}
	var (
		asJSON     = fs.Bool("json", false, "emit JSON output")
		iterations = fs.Int("n", 20, "number of iterations")
		list       = fs.Bool("list", false, "list built-in benchmarks")
		name       = fs.String("bench", "", "run a single built-in benchmark by name")
		raw        = fs.Bool("raw", false, "include all_ms in JSON output")
	)
	if err := fs.Parse(args); err != nil {
		return 2
	}

	if *list {
		for _, b := range builtIns {
			fmt.Fprintf(stdout, "%s\t%s\n", b.Name, b.Description)
		}
		return 0
	}

	var defs []benchmarkDef
	if *name != "" {
		found := false
		for _, b := range builtIns {
			if b.Name == *name {
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(stderr, "unknown benchmark: %s\n", *name)
			return 2
		}
		if err := resolveBuiltIns(); err != nil {
			fmt.Fprintf(stderr, "resolve built-ins: %v\n", err)
			return 2
		}
		for _, b := range builtIns {
			if b.Name == *name {
				defs = append(defs, b)
				break
			}
		}
	} else {
		// If no name given and extra args look like a command, run a custom benchmark.
		remaining := fs.Args()
		if len(remaining) > 0 {
			defs = append(defs, benchmarkDef{
				Name:        strings.Join(remaining, " "),
				Description: "custom command",
				Cmd:         remaining[0],
				Args:        remaining[1:],
			})
		} else {
			if err := resolveBuiltIns(); err != nil {
				fmt.Fprintf(stderr, "resolve built-ins: %v\n", err)
				return 2
			}
			defs = builtIns
		}
	}

	results := make([]benchmarkResult, 0, len(defs))
	for i := range defs {
		d := &defs[i]
		if d.Name == "sdp-trace-wrap" {
			tmpDir, err := os.MkdirTemp("", "ossbench-wrap-*")
			if err != nil {
				fmt.Fprintf(stderr, "mkdir temp: %v\n", err)
				return 2
			}
			d.Dir = tmpDir
			d.Cleanup = func() {
				_ = os.RemoveAll(tmpDir)
			}
		}
		res := runBenchmark(*d, *iterations)
		if d.Cleanup != nil {
			d.Cleanup()
		}
		if !*raw {
			res.AllMs = nil
		}
		results = append(results, res)
	}

	if err := printResults(stdout, results, *asJSON); err != nil {
		fmt.Fprintf(stderr, "print results: %v\n", err)
		cleanupTempBinary()
		return 2
	}
	cleanupTempBinary()
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
		var err error
		if r.Error != "" {
			_, err = fmt.Fprintf(w, "%*s  error: %s  median=%6.2f ms  min=%6.2f ms  max=%6.2f ms  attempted=%d succeeded=%d\n",
				-width, r.Name, r.Error, r.MedianMs, r.MinMs, r.MaxMs, r.AttemptedIterations, r.SucceededIterations)
		} else {
			_, err = fmt.Fprintf(w, "%*s  median=%6.2f ms  min=%6.2f ms  max=%6.2f ms  attempted=%d succeeded=%d\n",
				-width, r.Name, r.MedianMs, r.MinMs, r.MaxMs, r.AttemptedIterations, r.SucceededIterations)
		}
		if err != nil {
			return err
		}
	}
	return nil
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
	},
	{
		Name:        "sdp-trace-wrap",
		Description: "sdp-trace wrap /bin/true",
		Cmd:         "sdp-trace",
		Args:        []string{"wrap", "/bin/true"},
	},
}

// tempBinaryPath is set when the harness builds sdp-trace into a temp dir.
var tempBinaryPath string

// resolveBuiltIns resolves or builds the sdp-trace binary and updates builtIns.
// It prefers a repo-local binary; otherwise it builds into a temp dir so the
// checkout is not mutated.
func resolveBuiltIns() error {
	bin := resolveBinary("sdp-trace")
	if bin == "sdp-trace" {
		tmpDir, err := os.MkdirTemp("", "ossbench-bin-*")
		if err != nil {
			return fmt.Errorf("mkdir temp for build: %w", err)
		}
		bin = filepath.Join(tmpDir, "sdp-trace")
		if err := buildSDPTrace(bin); err != nil {
			_ = os.RemoveAll(tmpDir)
			return fmt.Errorf("sdp-trace binary not found and build failed: %w", err)
		}
		tempBinaryPath = bin
	}
	for i := range builtIns {
		builtIns[i].Cmd = bin
	}
	return nil
}

// cleanupTempBinary removes the temp-built binary if one was created.
func cleanupTempBinary() {
	if tempBinaryPath != "" {
		_ = os.RemoveAll(filepath.Dir(tempBinaryPath))
		tempBinaryPath = ""
	}
}
