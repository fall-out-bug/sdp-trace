package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseAndJoinRowsComputesCRAP(t *testing.T) {
	coverage, err := parseCoverage(strings.NewReader("github.com/fall_out_bug/sdp-trace/internal/demo/demo.go:10:\tBuild\t100.0%\n"))
	if err != nil {
		t.Fatalf("parse coverage: %v", err)
	}
	complexity, err := parseComplexity(strings.NewReader("4 demo Build internal/demo/demo.go:10:1\n"))
	if err != nil {
		t.Fatalf("parse complexity: %v", err)
	}
	rows, err := joinRows(complexity, coverage, false)
	if err != nil {
		t.Fatalf("join rows: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("rows = %d", len(rows))
	}
	if rows[0].crap != 4 {
		t.Fatalf("crap = %.2f", rows[0].crap)
	}
}

func TestJoinRowsSkipsTestsAndRejectsMissingCoverage(t *testing.T) {
	coverage, err := parseCoverage(strings.NewReader(""))
	if err != nil {
		t.Fatalf("parse coverage: %v", err)
	}
	complexity, err := parseComplexity(strings.NewReader("8 demo TestBuild internal/demo/demo_test.go:10:1\n6 demo Build internal/demo/demo.go:20:1\n"))
	if err != nil {
		t.Fatalf("parse complexity: %v", err)
	}
	if _, err := joinRows(complexity, coverage, false); err == nil {
		t.Fatalf("expected missing production coverage to fail")
	}
	rows, err := joinRows(complexity, coverage, true)
	if err != nil {
		t.Fatalf("join rows with missing coverage allowed: %v", err)
	}
	if len(rows) != 1 || rows[0].function != "Build" {
		t.Fatalf("rows = %+v", rows)
	}
}

func TestReadCoverageAndComplexityFromFiles(t *testing.T) {
	dir := t.TempDir()
	coveragePath := filepath.Join(dir, "coverage.txt")
	complexityPath := filepath.Join(dir, "gocyclo.txt")
	if err := os.WriteFile(coveragePath, []byte("github.com/fall_out_bug/sdp-trace/internal/demo/demo.go:10:\tBuild\t75.0%\n"), 0o600); err != nil {
		t.Fatalf("write coverage: %v", err)
	}
	if err := os.WriteFile(complexityPath, []byte("4 demo Build internal/demo/demo.go:10:1\n"), 0o600); err != nil {
		t.Fatalf("write complexity: %v", err)
	}

	coverage, err := readCoverage(coveragePath)
	if err != nil {
		t.Fatalf("read coverage: %v", err)
	}
	complexity, err := readComplexity(complexityPath)
	if err != nil {
		t.Fatalf("read complexity: %v", err)
	}

	if got := coverage["internal/demo/demo.go:10"].coverage; got != 75 {
		t.Fatalf("coverage = %.1f", got)
	}
	if len(complexity) != 1 || complexity[0].function != "Build" {
		t.Fatalf("complexity = %+v", complexity)
	}
}

func TestParseCoverageSkipsNonFunctionLinesAndReportsBadPercent(t *testing.T) {
	coverage, err := parseCoverage(strings.NewReader("\ntotal:\t(statements)\t80.0%\nnot coverage\n"))
	if err != nil {
		t.Fatalf("parse coverage skips: %v", err)
	}
	if len(coverage) != 0 {
		t.Fatalf("coverage = %+v", coverage)
	}

	_, err = parseCoverage(strings.NewReader("internal/demo/demo.go:10:\tBuild\t1.2.3%\n"))
	if err == nil || !strings.Contains(err.Error(), "parse coverage") {
		t.Fatalf("coverage parse err = %v", err)
	}
}

func TestParseComplexityRejectsUnrecognizedLine(t *testing.T) {
	_, err := parseComplexity(strings.NewReader("not gocyclo output\n"))
	if err == nil || !strings.Contains(err.Error(), "unrecognized gocyclo line") {
		t.Fatalf("complexity parse err = %v", err)
	}
}

func TestInputReadErrorsAreReported(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "missing.txt")
	if _, err := readCoverage(missingPath); err == nil {
		t.Fatalf("expected readCoverage to fail")
	}
	if _, err := readComplexity(missingPath); err == nil {
		t.Fatalf("expected readComplexity to fail")
	}
}

func TestLoadRowsReportsComplexityReadError(t *testing.T) {
	dir := t.TempDir()
	coveragePath := filepath.Join(dir, "coverage.txt")
	if err := os.WriteFile(coveragePath, []byte("internal/demo/demo.go:10:\tBuild\t100.0%\n"), 0o600); err != nil {
		t.Fatalf("write coverage: %v", err)
	}

	_, err := loadRows(options{coverPath: coveragePath, cycloPath: filepath.Join(dir, "missing-gocyclo.txt")})
	if err == nil || !strings.Contains(err.Error(), "read complexity") {
		t.Fatalf("loadRows err = %v", err)
	}
}

func TestRunReportsRowsAndThresholdFailures(t *testing.T) {
	dir := t.TempDir()
	coveragePath := filepath.Join(dir, "coverage.txt")
	complexityPath := filepath.Join(dir, "gocyclo.txt")
	if err := os.WriteFile(coveragePath, []byte("internal/demo/demo.go:10:\tBuild\t100.0%\n"), 0o600); err != nil {
		t.Fatalf("write coverage: %v", err)
	}
	if err := os.WriteFile(complexityPath, []byte("4 demo Build internal/demo/demo.go:10:1\n"), 0o600); err != nil {
		t.Fatalf("write complexity: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-cover-func", coveragePath, "-gocyclo", complexityPath, "-threshold", "5", "-strict-less"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("run exit = %d, stderr = %q", code, stderr.String())
	}
	if want := "internal/demo/demo.go:10 Build complexity=4 coverage=100.0 crap=4.00\n"; stdout.String() != want {
		t.Fatalf("stdout = %q, want %q", stdout.String(), want)
	}

	stdout.Reset()
	stderr.Reset()
	code = run([]string{"-cover-func", coveragePath, "-gocyclo", complexityPath, "-threshold", "4", "-strict-less"}, &stdout, &stderr)
	if code != 1 {
		t.Fatalf("threshold exit = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stderr.String(), "CRAP threshold 4.00 exceeded by 1 function(s)") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunRejectsMissingInputPath(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(nil, &stdout, &stderr)
	if code != 2 {
		t.Fatalf("run exit = %d", code)
	}
	if !strings.Contains(stderr.String(), "usage: crapcheck -cover-func cover.txt -gocyclo cyclo.txt [-threshold 5]") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestNormalizeFileVariants(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{" /tmp/work/sdp-trace/internal/demo/demo.go ", "internal/demo/demo.go"},
		{"cmd/sdp-trace/harness_cli.go", "cmd/sdp-trace/harness_cli.go"},
		{"sdp-trace/internal/demo/demo.go", "internal/demo/demo.go"},
		{"github.com/fall_out_bug/sdp-trace/internal/demo/demo.go", "internal/demo/demo.go"},
		{"./internal/demo/demo.go", "internal/demo/demo.go"},
	}
	for _, tc := range tests {
		if got := normalizeFile(tc.input); got != tc.want {
			t.Fatalf("normalizeFile(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestExceedsThresholdSupportsStrictLess(t *testing.T) {
	if exceedsThreshold(5, 5, false) {
		t.Fatalf("non-strict threshold should allow equality")
	}
	if !exceedsThreshold(5, 5, true) {
		t.Fatalf("strict threshold should reject equality")
	}
}

func TestCRAPScoreUsesCubicUncoveredRiskFormula(t *testing.T) {
	if got := crapScore(10, 50); got != 22.5 {
		t.Fatalf("crapScore(10, 50) = %.2f, want 22.50", got)
	}
}
