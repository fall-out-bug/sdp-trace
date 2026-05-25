package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestStats(t *testing.T) {
	tests := []struct {
		name       string
		values     []float64
		wantMin    float64
		wantMax    float64
		wantMedian float64
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
		Name: "go-env",
		Cmd:  "go",
		Args: []string{"env", "GOOS"},
		Dir:  os.TempDir(),
	}
	res := runBenchmark(def, 5)
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if res.AttemptedIterations != 5 {
		t.Errorf("attempted = %d, want 5", res.AttemptedIterations)
	}
	if res.SucceededIterations != 5 {
		t.Errorf("succeeded = %d, want 5", res.SucceededIterations)
	}
	if len(res.AllMs) != 5 {
		t.Errorf("got %d measurements, want 5", len(res.AllMs))
	}
	if res.MinMs < 0 || res.MaxMs < res.MinMs {
		t.Errorf("invalid min/max: min=%v max=%v", res.MinMs, res.MaxMs)
	}
	if len(res.Argv) != 3 || res.Argv[0] != "go" {
		t.Errorf("unexpected argv: %v", res.Argv)
	}
	if res.WorkingDirectory != os.TempDir() {
		t.Errorf("unexpected working_directory: %q", res.WorkingDirectory)
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
	if res.AttemptedIterations != 1 {
		t.Errorf("expected AttemptedIterations=1, got %d", res.AttemptedIterations)
	}
}

func TestRunBenchmark_DefaultIterations(t *testing.T) {
	def := benchmarkDef{
		Name: "go-env",
		Cmd:  "go",
		Args: []string{"env", "GOOS"},
	}
	res := runBenchmark(def, 0)
	if res.Error != "" {
		t.Fatalf("unexpected error: %s", res.Error)
	}
	if res.AttemptedIterations != 20 {
		t.Errorf("expected default 20 attempted, got %d", res.AttemptedIterations)
	}
}

func TestRunBenchmark_PartialFailure(t *testing.T) {
	// Verify that a command that always fails still produces partial results.
	def2 := benchmarkDef{
		Name: "always-fail",
		Cmd:  "go",
		Args: []string{"env", "-badflag"},
	}
	res2 := runBenchmark(def2, 3)
	if res2.Error == "" {
		t.Fatal("expected error for always-failing command")
	}
	if res2.AttemptedIterations != 3 {
		t.Errorf("expected AttemptedIterations=3, got %d", res2.AttemptedIterations)
	}
}

func TestRunBenchmark_MixedSuccessFailure(t *testing.T) {
	// Build a tiny helper that fails on first invocation and succeeds on second.
	src := `package main
import "os"
func main() {
	f := os.Args[1]
	if _, err := os.Stat(f); err == nil {
		os.Remove(f)
		return
	}
	os.WriteFile(f, []byte("x"), 0644)
	os.Exit(1)
}`
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "toggle.go")
	if err := os.WriteFile(srcFile, []byte(src), 0644); err != nil {
		t.Fatalf("write toggle source: %v", err)
	}
	bin := filepath.Join(tmpDir, "toggle")
	if err := exec.Command("go", "build", "-o", bin, srcFile).Run(); err != nil {
		t.Skip("cannot build toggle helper:", err)
	}
	sentinel := filepath.Join(tmpDir, "sentinel")
	def := benchmarkDef{
		Name: "toggle",
		Cmd:  bin,
		Args: []string{sentinel},
	}
	// First iteration fails (no sentinel), second succeeds (sentinel exists),
	// third fails (sentinel removed).
	res := runBenchmark(def, 3)
	if res.AttemptedIterations != 3 {
		t.Errorf("attempted = %d, want 3", res.AttemptedIterations)
	}
	if res.SucceededIterations != 1 {
		t.Errorf("succeeded = %d, want 1", res.SucceededIterations)
	}
	if len(res.AllMs) != 1 {
		t.Errorf("got %d measurements, want 1", len(res.AllMs))
	}
	if res.Error == "" {
		t.Error("expected error because at least one iteration failed")
	}
}

func TestRepoRoot(t *testing.T) {
	tmpDir := t.TempDir()
	gitDir := filepath.Join(tmpDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}
	subDir := filepath.Join(tmpDir, "a", "b")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("mkdir subdir: %v", err)
	}
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	}()
	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	got := repoRoot()
	if got != tmpDir {
		t.Errorf("repoRoot() = %q, want %q", got, tmpDir)
	}
}

func TestRepoRoot_Fallback(t *testing.T) {
	tmpDir := t.TempDir()
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		if err := os.Chdir(originalWd); err != nil {
			t.Fatalf("restore wd: %v", err)
		}
	}()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	got := repoRoot()
	if got != "." {
		t.Errorf("repoRoot() = %q, want %q", got, ".")
	}
}

func TestBuildSDPTrace_Error(t *testing.T) {
	tmpDir := t.TempDir()
	// Create a file where a directory is expected.
	blockingFile := filepath.Join(tmpDir, "block")
	if err := os.WriteFile(blockingFile, []byte("x"), 0644); err != nil {
		t.Fatalf("write blocking file: %v", err)
	}
	outPath := filepath.Join(blockingFile, "sdp-trace")
	err := buildSDPTrace(outPath)
	if err == nil {
		t.Fatal("expected error for invalid output path")
	}
}

func TestBuildSDPTrace(t *testing.T) {
	skipUnlessIntegration(t)
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "sdp-trace")
	if err := buildSDPTrace(outPath); err != nil {
		t.Fatalf("buildSDPTrace: %v", err)
	}
	if _, err := os.Stat(outPath); err != nil {
		t.Fatalf("output binary not found: %v", err)
	}
}
