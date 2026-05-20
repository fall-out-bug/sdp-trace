package main

import (
	"context"
	"fmt"
	"os/exec"
	"sort"
	"time"
)

// benchmarkDef defines a single benchmark target.
type benchmarkDef struct {
	Name        string
	Description string
	Cmd         string
	Args        []string
}

// benchmarkResult holds the measured statistics for one benchmark.
type benchmarkResult struct {
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Iterations  int       `json:"iterations"`
	MinMs       float64   `json:"min_ms"`
	MaxMs       float64   `json:"max_ms"`
	MedianMs    float64   `json:"median_ms"`
	AllMs       []float64 `json:"all_ms,omitempty"`
	Error       string    `json:"error,omitempty"`
}

// runBenchmark executes a benchmark definition for n iterations.
func runBenchmark(def benchmarkDef, iterations int) benchmarkResult {
	if iterations <= 0 {
		iterations = 20
	}
	if def.Cmd == "" {
		return benchmarkResult{
			Name:       def.Name,
			Iterations: iterations,
			Error:      "no command specified",
		}
	}

	times := make([]time.Duration, 0, iterations)
	var lastErr string
	for i := 0; i < iterations; i++ {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		cmd := exec.CommandContext(ctx, def.Cmd, def.Args...)
		err := cmd.Run()
		cancel()
		if err != nil {
			lastErr = fmt.Sprintf("iteration %d failed: %v", i, err)
			continue
		}
		times = append(times, time.Since(start))
	}

	if len(times) == 0 {
		return benchmarkResult{
			Name:       def.Name,
			Iterations: iterations,
			Error:      lastErr,
		}
	}

	ms := make([]float64, len(times))
	for i, d := range times {
		ms[i] = float64(d) / float64(time.Millisecond)
	}
	min, max, median := stats(ms)
	return benchmarkResult{
		Name:        def.Name,
		Description: def.Description,
		Iterations:  iterations,
		MinMs:       min,
		MaxMs:       max,
		MedianMs:    median,
		AllMs:       ms,
		Error:       lastErr,
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
