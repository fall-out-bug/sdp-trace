package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func skipUnlessIntegration(t *testing.T) {
	if os.Getenv("SDPTRACE_INTEGRATION") != "1" {
		t.Skip("set SDPTRACE_INTEGRATION=1 to run tests that build the product binary")
	}
}

func TestRun_List(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-list"}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	for _, b := range builtIns {
		if !strings.Contains(out.String(), b.Name) {
			t.Errorf("output missing benchmark %q", b.Name)
		}
	}
}

func TestRun_CustomCommand(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-n", "3", "go", "env", "GOOS"}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out.String(), "go env GOOS") {
		t.Fatalf("expected output for custom command, got: %s", out.String())
	}
}

func TestRun_JSON(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-json", "-n", "2", "go", "env", "GOOS"}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	var results []benchmarkResult
	if err := json.Unmarshal(out.Bytes(), &results); err != nil {
		t.Fatalf("expected valid JSON output, got: %s, err: %v", out.String(), err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 benchmark result, got %d", len(results))
	}
	if results[0].Name != "go env GOOS" {
		t.Fatalf("expected name %q, got %q", "go env GOOS", results[0].Name)
	}
	if results[0].Error != "" {
		t.Fatalf("unexpected error: %s", results[0].Error)
	}
}

func TestRun_BuiltinBench(t *testing.T) {
	skipUnlessIntegration(t)
	var out bytes.Buffer
	code := run([]string{"-bench", "sdp-trace-version", "-n", "1"}, &out, &out)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d", code)
	}
	if !strings.Contains(out.String(), "sdp-trace-version") {
		t.Fatalf("expected output for builtin bench, got: %s", out.String())
	}
}

func TestRun_UnknownBench(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"-bench", "nonexistent"}, &out, &out)
	if code != 2 {
		t.Fatalf("expected exit 2, got %d", code)
	}
}

func TestRun_MissingCommand(t *testing.T) {
	var out bytes.Buffer
	code := run([]string{"this-command-does-not-exist-017"}, &out, &out)
	if code != 1 {
		t.Fatalf("expected exit 1, got %d", code)
	}
}

func TestExitCode(t *testing.T) {
	tests := []struct {
		name    string
		results []benchmarkResult
		want    int
	}{
		{"all ok", []benchmarkResult{{Name: "a", MedianMs: 1}}, 0},
		{"one error", []benchmarkResult{{Name: "a", Error: "fail"}}, 1},
		{"mixed", []benchmarkResult{{Name: "a", MedianMs: 1}, {Name: "b", Error: "fail"}}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exitCode(tt.results); got != tt.want {
				t.Errorf("exitCode() = %d, want %d", got, tt.want)
			}
		})
	}
}

func saveBuiltins() (savedBuiltIns, savedBuiltInsOrig []benchmarkDef, savedTempPath string) {
	savedBuiltIns = make([]benchmarkDef, len(builtIns))
	copy(savedBuiltIns, builtIns)
	savedBuiltInsOrig = make([]benchmarkDef, len(builtInsOrig))
	copy(savedBuiltInsOrig, builtInsOrig)
	return savedBuiltIns, savedBuiltInsOrig, tempBinaryPath
}

func restoreBuiltins(savedBuiltIns, savedBuiltInsOrig []benchmarkDef, savedTempPath string) {
	builtIns = savedBuiltIns
	builtInsOrig = savedBuiltInsOrig
	tempBinaryPath = savedTempPath
}

func TestFindBuiltin(t *testing.T) {
	def, ok := findBuiltin("sdp-trace-version")
	if !ok {
		t.Fatal("expected to find sdp-trace-version")
	}
	if def.Name != "sdp-trace-version" {
		t.Errorf("expected name %q, got %q", "sdp-trace-version", def.Name)
	}
	_, ok = findBuiltin("nonexistent")
	if ok {
		t.Error("expected not to find nonexistent")
	}
}

func TestFormatResultLine(t *testing.T) {
	tests := []struct {
		name   string
		result benchmarkResult
		width  int
		want   string
	}{
		{
			name: "success",
			result: benchmarkResult{
				Name:                "bench",
				MedianMs:            1.5,
				MinMs:               1.0,
				MaxMs:               2.0,
				AttemptedIterations: 10,
				SucceededIterations: 10,
			},
			width: 5,
			want:  "bench  median=  1.50 ms  min=  1.00 ms  max=  2.00 ms  attempted=10 succeeded=10\n",
		},
		{
			name: "error",
			result: benchmarkResult{
				Name:                "bench",
				Error:               "fail",
				MedianMs:            0,
				MinMs:               0,
				MaxMs:               0,
				AttemptedIterations: 1,
				SucceededIterations: 0,
			},
			width: 5,
			want:  "bench  error: fail  median=  0.00 ms  min=  0.00 ms  max=  0.00 ms  attempted=1 succeeded=0\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatResultLine(tt.result, tt.width)
			if got != tt.want {
				t.Errorf("formatResultLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResolveBenchmarkDefs_Custom(t *testing.T) {
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	cfg := runConfig{args: []string{"go", "env", "GOOS"}}
	defs, cleanup, err := resolveBenchmarkDefs(cfg)
	if err != nil {
		t.Fatalf("resolveBenchmarkDefs: %v", err)
	}
	if cleanup != nil {
		t.Fatal("expected no cleanup for custom command")
	}
	if len(defs) != 1 {
		t.Fatalf("expected 1 def, got %d", len(defs))
	}
	if defs[0].Name != "go env GOOS" {
		t.Errorf("expected name %q, got %q", "go env GOOS", defs[0].Name)
	}
}

func TestResolveBenchmarkDefs_Unknown(t *testing.T) {
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	cfg := runConfig{name: "nonexistent"}
	_, _, err := resolveBenchmarkDefs(cfg)
	if err == nil {
		t.Fatal("expected error for unknown benchmark")
	}
}

func TestResolveBenchmarkDefs_UnexpectedArgs(t *testing.T) {
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	cfg := runConfig{name: "sdp-trace-version", args: []string{"extra"}}
	_, _, err := resolveBenchmarkDefs(cfg)
	if err == nil {
		t.Fatal("expected error for unexpected args")
	}
}

func TestResolveBenchmarkDefs_AllBuiltins(t *testing.T) {
	skipUnlessIntegration(t)
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	cfg := runConfig{}
	defs, cleanup, err := resolveBenchmarkDefs(cfg)
	if err != nil {
		t.Fatalf("resolveBenchmarkDefs: %v", err)
	}
	if cleanup == nil {
		t.Fatal("expected cleanup for builtins")
	}
	defer cleanup()
	if len(defs) != len(builtIns) {
		t.Errorf("expected %d defs, got %d", len(builtIns), len(defs))
	}
}

func TestResolveBuiltIns(t *testing.T) {
	skipUnlessIntegration(t)
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	if err := resolveBuiltIns(); err != nil {
		t.Fatalf("resolveBuiltIns: %v", err)
	}
	if tempBinaryPath == "" {
		t.Fatal("expected tempBinaryPath to be set")
	}
	if _, err := os.Stat(tempBinaryPath); err != nil {
		t.Fatalf("temp binary not found: %v", err)
	}
	for _, b := range builtIns {
		if b.Source != "temp-build" {
			t.Errorf("expected Source='temp-build', got %q", b.Source)
		}
		if b.Cmd != tempBinaryPath {
			t.Errorf("expected Cmd=%q, got %q", tempBinaryPath, b.Cmd)
		}
	}

	cleanupTempBinary()
	if tempBinaryPath != "" {
		t.Error("expected tempBinaryPath to be cleared")
	}
	for i, b := range builtIns {
		if b.Cmd != savedBI[i].Cmd {
			t.Errorf("builtIns[%d].Cmd not restored", i)
		}
		if b.Source != savedBI[i].Source {
			t.Errorf("builtIns[%d].Source not restored", i)
		}
	}
}

func TestCleanupTempBinary(t *testing.T) {
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "sdp-trace")
	if err := os.WriteFile(binPath, []byte("fake"), 0644); err != nil {
		t.Fatalf("write fake binary: %v", err)
	}

	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	tempBinaryPath = binPath
	builtIns[0].Cmd = binPath
	builtIns[0].Source = "temp-build"

	cleanupTempBinary()
	if tempBinaryPath != "" {
		t.Error("expected tempBinaryPath to be cleared")
	}
	if _, err := os.Stat(binPath); !os.IsNotExist(err) {
		t.Error("expected binary to be removed")
	}
	if builtIns[0].Cmd != savedBI[0].Cmd {
		t.Errorf("builtIns[0].Cmd not restored, got %q", builtIns[0].Cmd)
	}
	if builtIns[0].Source != savedBI[0].Source {
		t.Errorf("builtIns[0].Source not restored, got %q", builtIns[0].Source)
	}
}

func TestSetupWrap(t *testing.T) {
	def := benchmarkDef{Name: "sdp-trace-wrap"}
	if err := setupWrap(&def); err != nil {
		t.Fatalf("setupWrap: %v", err)
	}
	if def.Dir == "" {
		t.Fatal("expected Dir to be set")
	}
	if def.Cleanup == nil {
		t.Fatal("expected Cleanup to be set")
	}
	def.Cleanup()
	if _, err := os.Stat(def.Dir); !os.IsNotExist(err) {
		t.Error("expected temp dir to be removed")
	}
}

func TestSetupWrap_Other(t *testing.T) {
	def := benchmarkDef{Name: "other"}
	if err := setupWrap(&def); err != nil {
		t.Fatalf("setupWrap: %v", err)
	}
	if def.Dir != "" {
		t.Error("expected Dir to be empty for non-wrap benchmark")
	}
	if def.Cleanup != nil {
		t.Error("expected Cleanup to be nil for non-wrap benchmark")
	}
}

func TestResolveBuiltIns_BuildError(t *testing.T) {
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	oldBuildBinary := buildBinary
	buildBinary = func(outPath string) error {
		return fmt.Errorf("mock build failure")
	}
	defer func() { buildBinary = oldBuildBinary }()

	err := resolveBuiltIns()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "mock build failure") {
		t.Fatalf("unexpected error: %v", err)
	}
	if tempBinaryPath != "" {
		t.Error("expected tempBinaryPath to be empty")
	}
}

func TestResolveBuiltIns_MockSuccess(t *testing.T) {
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	oldBuildBinary := buildBinary
	buildBinary = func(outPath string) error {
		return os.WriteFile(outPath, []byte("fake"), 0644)
	}
	defer func() { buildBinary = oldBuildBinary }()

	if err := resolveBuiltIns(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tempBinaryPath == "" {
		t.Fatal("expected tempBinaryPath to be set")
	}
	for _, b := range builtIns {
		if b.Source != "temp-build" {
			t.Errorf("expected Source='temp-build', got %q", b.Source)
		}
		if b.Cmd != tempBinaryPath {
			t.Errorf("expected Cmd=%q, got %q", tempBinaryPath, b.Cmd)
		}
	}
	cleanupTempBinary()
}

func TestResolveAllBuiltins_BuildError(t *testing.T) {
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	oldBuildBinary := buildBinary
	buildBinary = func(outPath string) error {
		return fmt.Errorf("mock build failure")
	}
	defer func() { buildBinary = oldBuildBinary }()

	defs, cleanup, err := resolveAllBuiltins()
	if err == nil {
		t.Fatal("expected error")
	}
	if cleanup != nil {
		t.Error("expected nil cleanup on error")
	}
	if len(defs) != 0 {
		t.Errorf("expected no defs, got %d", len(defs))
	}
}

func TestResolveAllBuiltins_MockSuccess(t *testing.T) {
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	oldBuildBinary := buildBinary
	buildBinary = func(outPath string) error {
		return os.WriteFile(outPath, []byte("fake"), 0644)
	}
	defer func() { buildBinary = oldBuildBinary }()

	defs, cleanup, err := resolveAllBuiltins()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cleanup == nil {
		t.Fatal("expected cleanup")
	}
	defer cleanup()
	if len(defs) != len(builtIns) {
		t.Errorf("expected %d defs, got %d", len(builtIns), len(defs))
	}
}

func TestResolveSingleBuiltin_BuildError(t *testing.T) {
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	oldBuildBinary := buildBinary
	buildBinary = func(outPath string) error {
		return fmt.Errorf("mock build failure")
	}
	defer func() { buildBinary = oldBuildBinary }()

	defs, cleanup, err := resolveSingleBuiltin("sdp-trace-version")
	if err == nil {
		t.Fatal("expected error")
	}
	if cleanup != nil {
		t.Error("expected nil cleanup on error")
	}
	if len(defs) != 0 {
		t.Errorf("expected no defs, got %d", len(defs))
	}
}

func TestResolveSingleBuiltin_MockSuccess(t *testing.T) {
	savedBI, savedBIO, savedTP := saveBuiltins()
	defer restoreBuiltins(savedBI, savedBIO, savedTP)

	oldBuildBinary := buildBinary
	buildBinary = func(outPath string) error {
		return os.WriteFile(outPath, []byte("fake"), 0644)
	}
	defer func() { buildBinary = oldBuildBinary }()

	defs, cleanup, err := resolveSingleBuiltin("sdp-trace-version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cleanup == nil {
		t.Fatal("expected cleanup")
	}
	defer cleanup()
	if len(defs) != 1 {
		t.Fatalf("expected 1 def, got %d", len(defs))
	}
	if defs[0].Name != "sdp-trace-version" {
		t.Errorf("expected name sdp-trace-version, got %q", defs[0].Name)
	}
	if defs[0].Cmd != tempBinaryPath {
		t.Errorf("expected Cmd=%q, got %q", tempBinaryPath, defs[0].Cmd)
	}
	if defs[0].Source != "temp-build" {
		t.Errorf("expected Source=temp-build, got %q", defs[0].Source)
	}
}
