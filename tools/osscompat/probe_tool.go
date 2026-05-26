package main

import (
	"context"
	"os/exec"
)

// hasTool reports whether tool is in $PATH.
func hasTool(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}

// runExternalTool executes an external command and returns combined output.
func runExternalTool(ctx context.Context, name string, args ...string) ([]byte, error) {
	return exec.CommandContext(ctx, name, args...).CombinedOutput()
}
