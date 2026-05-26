package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

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

// buildBinary is used by resolveBuiltIns to compile the sdp-trace binary.
// Tests may replace it to avoid real builds.
var buildBinary = buildSDPTrace
