package main

import (
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
