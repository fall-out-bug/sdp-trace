package main

import (
	"fmt"
	"os"
	"path/filepath"
)

// tempBinaryPath is set when the harness builds sdp-trace into a temp dir.
var tempBinaryPath string

// resolveBuiltIns builds the sdp-trace binary from current source into a temp
// dir and updates builtIns. It never uses a pre-existing repo-root binary so
// that benchmark results always reflect the checked-out source.
func resolveBuiltIns() error {
	tmpDir, err := os.MkdirTemp("", "ossbench-bin-*")
	if err != nil {
		return fmt.Errorf("mkdir temp for build: %w", err)
	}
	bin := filepath.Join(tmpDir, "sdp-trace")
	if err := buildBinary(bin); err != nil {
		_ = os.RemoveAll(tmpDir)
		return fmt.Errorf("sdp-trace build failed: %w", err)
	}
	tempBinaryPath = bin
	for i := range builtIns {
		builtIns[i].Cmd = bin
		builtIns[i].Source = "temp-build"
	}
	return nil
}

// cleanupTempBinary removes the temp-built binary if one was created and
// restores builtIns to their original values so subsequent runs in the same
// process do not reference a deleted path.
func cleanupTempBinary() {
	if tempBinaryPath != "" {
		_ = os.RemoveAll(filepath.Dir(tempBinaryPath))
		tempBinaryPath = ""
	}
	copy(builtIns, builtInsOrig)
}
