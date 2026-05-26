package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// runOPAEval executes OPA evaluation and returns the boolean result of the query.
func runOPAEval(ctx context.Context, regoPath, fixturePath, query string) (bool, error) {
	cmd := exec.CommandContext(ctx, "opa", "eval",
		"--data", regoPath,
		"--input", fixturePath,
		"--format", "json",
		query,
	)
	stdout, err := cmd.Output()
	if err != nil {
		var stderr []byte
		if exitErr, ok := err.(*exec.ExitError); ok {
			stderr = exitErr.Stderr
		}
		return false, fmt.Errorf("opa eval failed: %w\nstderr: %s", err, strings.TrimSpace(string(stderr)))
	}
	return parseOPAEvalResult(stdout)
}
