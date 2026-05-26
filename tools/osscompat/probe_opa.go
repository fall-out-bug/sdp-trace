package main

import (
	"context"
	"time"
)

// runOPAPolicy evaluates the checked-in adapter.rego against the test fixture.
func runOPAPolicy() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if s, r := opaPreflight(ctx); s != "" {
		return s, r
	}
	return runOPAPolicyEval(ctx, "test-fixture.json", true, "test-fixture.json")
}

// runOPANegativeFixture evaluates the checked-in adapter.rego against the
// negative test fixture and asserts pass is false.
func runOPANegativeFixture() (verifierState, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if s, r := opaPreflight(ctx); s != "" {
		return s, r
	}
	return runOPAPolicyEval(ctx, "test-fixture-fail.json", false, "test-fixture-fail.json")
}
