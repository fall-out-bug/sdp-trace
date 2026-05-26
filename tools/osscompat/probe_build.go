package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// buildSDPTraceInTemp builds the sdp-trace binary in a temporary directory.
func buildSDPTraceInTemp(ctx context.Context, root string) (bin string, tmpDir string, err error) {
	tmpDir, err = os.MkdirTemp("", "osscompat-wrap-*")
	if err != nil {
		return "", "", fmt.Errorf("mkdir temp: %w", err)
	}
	bin = filepath.Join(tmpDir, "sdp-trace")
	buildCmd := exec.CommandContext(ctx, "go", "build", "-o", bin, "./cmd/sdp-trace")
	buildCmd.Dir = root
	buildOut, err := buildCmd.CombinedOutput()
	if err != nil {
		return "", tmpDir, fmt.Errorf("go build failed: %w\n%s", err, strings.TrimSpace(string(buildOut)))
	}
	return bin, tmpDir, nil
}
