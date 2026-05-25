package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// benchmarkDef defines a single benchmark target.
type benchmarkDef struct {
	Name        string
	Description string
	Cmd         string
	Args        []string
	Dir         string // working directory; empty means current directory
	Cleanup     func() // optional cleanup after benchmark completes
	Source      string // "repo-root", "temp-build", or "PATH"
}

func repoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	for cwd != "" {
		if _, err := os.Stat(filepath.Join(cwd, ".git")); err == nil {
			return cwd
		}
		cwd = parentDir(cwd)
	}
	return "."
}

func parentDir(dir string) string {
	parent := filepath.Dir(dir)
	if parent == dir {
		return ""
	}
	return parent
}

// buildSDPTrace compiles the sdp-trace binary from source into outPath.
func buildSDPTrace(outPath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "build", "-o", outPath, "./cmd/sdp-trace")
	cmd.Dir = repoRoot()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go build failed: %v\n%s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

// buildBinary is used by resolveBuiltIns to compile the sdp-trace binary.
// Tests may replace it to avoid real builds.
var buildBinary = buildSDPTrace

// benchmarkResult holds the measured statistics for one benchmark.
type benchmarkResult struct {
	Name                string    `json:"name"`
	Description         string    `json:"description,omitempty"`
	Command             string    `json:"command"`
	Argv                []string  `json:"argv,omitempty"`
	WorkingDirectory    string    `json:"working_directory,omitempty"`
	BinaryPath          string    `json:"binary_path,omitempty"`
	BinarySource        string    `json:"binary_source,omitempty"`
	Environment         envInfo   `json:"environment,omitempty"`
	AttemptedIterations int       `json:"attempted_iterations"`
	SucceededIterations int       `json:"succeeded_iterations"`
	MinMs               float64   `json:"min_ms"`
	MaxMs               float64   `json:"max_ms"`
	MedianMs            float64   `json:"median_ms"`
	AllMs               []float64 `json:"all_ms,omitempty"`
	Error               string    `json:"error,omitempty"`
}

// runBenchmark executes a benchmark definition for n iterations.
func runBenchmark(def benchmarkDef, iterations int) benchmarkResult {
	if iterations <= 0 {
		iterations = 20
	}
	if def.Cmd == "" {
		return benchmarkResult{
			Name:                def.Name,
			Description:         def.Description,
			Command:             "",
			BinaryPath:          "",
			BinarySource:        def.Source,
			Environment:         getEnv(),
			AttemptedIterations: iterations,
			Error:               "no command specified",
		}
	}

	times, attempted, lastErr := runIterations(def, iterations)
	argv := append([]string{filepath.Base(def.Cmd)}, def.Args...)
	cmdDisplay := strings.Join(argv, " ")
	return buildResult(def, cmdDisplay, argv, attempted, times, lastErr)
}

func runIterations(def benchmarkDef, iterations int) ([]time.Duration, int, string) {
	times := make([]time.Duration, 0, iterations)
	var lastErr string
	attempted := 0
	for i := 0; i < iterations; i++ {
		attempted++
		dur, err, timedOut := runSingleCommand(def.Cmd, def.Args, def.Dir)
		if err != nil {
			lastErr = fmt.Sprintf("iteration %d failed: %v", i, err)
			if timedOut {
				break
			}
			continue
		}
		times = append(times, dur)
	}
	return times, attempted, lastErr
}

func runSingleCommand(cmd string, args []string, dir string) (time.Duration, error, bool) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	c := exec.CommandContext(ctx, cmd, args...)
	if dir != "" {
		c.Dir = dir
	}
	err := c.Run()
	if err != nil {
		return 0, err, ctx.Err() == context.DeadlineExceeded
	}
	return time.Since(start), nil, false
}

func buildResult(def benchmarkDef, cmdDisplay string, argv []string, attempted int, times []time.Duration, lastErr string) benchmarkResult {
	res := benchmarkResult{
		Name:                def.Name,
		Description:         def.Description,
		Command:             cmdDisplay,
		Argv:                argv,
		WorkingDirectory:    def.Dir,
		BinaryPath:          def.Cmd,
		BinarySource:        def.Source,
		Environment:         getEnv(),
		AttemptedIterations: attempted,
		Error:               lastErr,
	}
	if len(times) == 0 {
		return res
	}
	ms := make([]float64, len(times))
	for i, d := range times {
		ms[i] = float64(d) / float64(time.Millisecond)
	}
	min, max, median := stats(ms)
	res.SucceededIterations = len(times)
	res.MinMs = min
	res.MaxMs = max
	res.MedianMs = median
	res.AllMs = ms
	return res
}

// envInfo captures the benchmark environment.
type envInfo struct {
	Platform  string `json:"platform"`
	GoOS      string `json:"goos"`
	GoArch    string `json:"goarch"`
	GoVersion string `json:"go_version"`
}

func getEnv() envInfo {
	return envInfo{
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
		GoOS:      runtime.GOOS,
		GoArch:    runtime.GOARCH,
		GoVersion: runtime.Version(),
	}
}

// stats returns min, max, and median of a non-empty slice.
func stats(values []float64) (min, max, median float64) {
	if len(values) == 0 {
		return 0, 0, 0
	}
	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)
	min = sorted[0]
	max = sorted[len(sorted)-1]
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		median = sorted[mid]
	} else {
		median = (sorted[mid-1] + sorted[mid]) / 2
	}
	return min, max, median
}
