package main

import (
	"testing"
)

func TestStats(t *testing.T) {
	tests := []struct {
		name         string
		values       []float64
		wantMin      float64
		wantMax      float64
		wantMedian   float64
	}{
		{"odd count", []float64{3, 1, 2}, 1, 3, 2},
		{"even count", []float64{4, 1, 3, 2}, 1, 4, 2.5},
		{"single", []float64{42}, 42, 42, 42},
		{"sorted", []float64{1, 2, 3, 4, 5}, 1, 5, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max, median := stats(tt.values)
			if min != tt.wantMin {
				t.Errorf("min = %v, want %v", min, tt.wantMin)
			}
			if max != tt.wantMax {
				t.Errorf("max = %v, want %v", max, tt.wantMax)
			}
			if median != tt.wantMedian {
				t.Errorf("median = %v, want %v", median, tt.wantMedian)
			}
		})
	}
}

func TestStats_Empty(t *testing.T) {
	min, max, median := stats(nil)
	if min != 0 || max != 0 || median != 0 {
		t.Errorf("expected all zeros for empty slice, got min=%v max=%v median=%v", min, max, median)
	}
}

func TestRunBenchmark_CustomCommand(t *testing.T) {
	def := benchmarkDef{
		Name: "true",
		Cmd:  "true",
		Args: []string{},
	}
	res := runBenchmark(def, 5)
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if res.Iterations != 5 {
		t.Errorf("iterations = %d, want 5", res.Iterations)
	}
	if len(res.AllMs) != 5 {
		t.Errorf("got %d measurements, want 5", len(res.AllMs))
	}
	if res.MinMs < 0 || res.MaxMs < res.MinMs {
		t.Errorf("invalid min/max: min=%v max=%v", res.MinMs, res.MaxMs)
	}
}

func TestRunBenchmark_MissingCommand(t *testing.T) {
	def := benchmarkDef{
		Name: "missing",
		Cmd:  "this-command-does-not-exist-017",
	}
	res := runBenchmark(def, 1)
	if res.Error == "" {
		t.Fatal("expected error for missing command")
	}
}

func TestRunBenchmark_NoCommand(t *testing.T) {
	res := runBenchmark(benchmarkDef{Name: "empty"}, 1)
	if res.Error == "" {
		t.Fatal("expected error when no command specified")
	}
}
