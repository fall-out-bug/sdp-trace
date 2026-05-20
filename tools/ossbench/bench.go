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

// repoRoot returns the repository root by walking up from the current
// working directory until it finds a .git directory or reaches the filesystem
// root. It falls back to "." if the root cannot be determined.
func repoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(cwd, ".git")); err == nil {
			return cwd
		}
		parent := filepath.Dir(cwd)
		if parent == cwd {
			return "."
		}
		cwd = parent
	}
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

// benchmarkResult holds the measured statistics for one benchmark.
type benchmarkResult struct {
	Name                string    `json:"name"`
	Description         string    `json:"description,omitempty"`
	Command             string    `json:"command"`
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

	times := make([]time.Duration, 0, iterations)
	var lastErr string
	for i := 0; i < iterations; i++ {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		cmd := exec.CommandContext(ctx, def.Cmd, def.Args...)
		if def.Dir != "" {
			cmd.Dir = def.Dir
		}
		err := cmd.Run()
		cancel()
		if err != nil {
			lastErr = fmt.Sprintf("iteration %d failed: %v", i, err)
			continue
		}
		times = append(times, time.Since(start))
	}

	cmdDisplay := filepath.Base(def.Cmd) + " " + strings.Join(def.Args, " ")
	if len(times) == 0 {
		return benchmarkResult{
			Name:                def.Name,
			Description:         def.Description,
			Command:             cmdDisplay,
			BinaryPath:          def.Cmd,
			BinarySource:        def.Source,
			Environment:         getEnv(),
			AttemptedIterations: iterations,
			Error:               lastErr,
		}
	}

	ms := make([]float64, len(times))
	for i, d := range times {
		ms[i] = float64(d) / float64(time.Millisecond)
	}
	min, max, median := stats(ms)
	return benchmarkResult{
		Name:                def.Name,
		Description:         def.Description,
		Command:             cmdDisplay,
		BinaryPath:          def.Cmd,
		BinarySource:        def.Source,
		Environment:         getEnv(),
		AttemptedIterations: iterations,
		SucceededIterations: len(times),
		MinMs:               min,
		MaxMs:               max,
		MedianMs:            median,
		AllMs:               ms,
		Error:               lastErr,
	}
}

// envInfo captures the benchmark environment.
type envInfo struct {
	Platform string `json:"platform"`
	GoOS     string `json:"goos"`
	GoArch   string `json:"goarch"`
}

func getEnv() envInfo {
	return envInfo{
		Platform: runtime.GOOS + "/" + runtime.GOARCH,
		GoOS:     runtime.GOOS,
		GoArch:   runtime.GOARCH,
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
