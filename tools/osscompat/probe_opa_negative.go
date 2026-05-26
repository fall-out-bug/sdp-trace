package main

import (
	"context"
	"time"
)

// runOPANegativeTraceID evaluates adapter.rego against the trace_id-only
// negative fixture and asserts pass is false.
func runOPANegativeTraceID() (verifierState, string) {
	return runOPANegativeFixturePath("test-fixture-fail-traceid.json", "trace_id-only")
}

// runOPANegativeProvenance evaluates adapter.rego against the provenance-only
// negative fixture and asserts pass is false.
func runOPANegativeProvenance() (verifierState, string) {
	return runOPANegativeFixturePath("test-fixture-fail-provenance.json", "provenance-only")
}

// runOPANegativeFixturePath is a shared helper for per-rule negative fixtures.
func runOPANegativeFixturePath(fixtureName, label string) (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if s, r := opaPreflight(ctx); s != "" {
		return s, r
	}
	return runOPAPolicyEval(ctx, fixtureName, false, label)
}
